import { existsSync, readFileSync } from 'node:fs';
import { hostname } from 'node:os';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { spawn } from 'node:child_process';

const scriptDir = dirname(fileURLToPath(import.meta.url));
const frontendDir = resolve(scriptDir, '..');
const defaultEnvFile = resolve(frontendDir, '..', 'backend', '.env');
const envFile = process.env.FRONTEND_ENV_FILE || process.env.ENV_FILE || defaultEnvFile;
const ngBin = resolve(frontendDir, 'node_modules', '.bin', process.platform === 'win32' ? 'ng.cmd' : 'ng');

loadEnvFile(envFile);

const userArgs = process.argv.slice(2);
const args = ['serve', ...userArgs];

appendAllowedHosts(args);
appendDefaultFlag(args, '--hmr', envBool('FRONTEND_HMR', false));
appendDefaultFlag(args, '--live-reload', envBool('FRONTEND_LIVE_RELOAD', false));
appendDefaultFlag(args, '--prebundle', envBool('FRONTEND_PREBUNDLE', false));

const child = spawn(ngBin, args, {
  cwd: frontendDir,
  env: process.env,
  stdio: 'inherit',
  shell: process.platform === 'win32',
});

child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }

  process.exit(code ?? 0);
});

function loadEnvFile(path) {
  if (!existsSync(path)) {
    return;
  }

  const text = readFileSync(path, 'utf8');
  for (const line of text.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('#')) {
      continue;
    }

    const match = trimmed.match(/^([A-Za-z_][A-Za-z0-9_]*)=(.*)$/);
    if (!match) {
      continue;
    }

    const [, key, rawValue] = match;
    if (process.env[key] !== undefined) {
      continue;
    }

    process.env[key] = stripQuotes(rawValue.trim());
  }
}

function stripQuotes(value) {
  if (value.length >= 2) {
    const first = value[0];
    const last = value[value.length - 1];
    if ((first === '"' && last === '"') || (first === "'" && last === "'")) {
      return value.slice(1, -1);
    }
  }

  return value;
}

function appendAllowedHosts(args) {
  if (hasFlag(args, '--allowed-hosts') || hasFlag(args, '--disable-host-check')) {
    return;
  }

  const hosts = unique([
    'localhost',
    '127.0.0.1',
    hostname(),
    ...splitList(process.env.FRONTEND_ALLOWED_HOSTS),
    ...splitList(process.env.ALLOWED_HOSTS),
    ...hostsFromOrigins(process.env.CORS_ORIGINS),
    hostFromUrl(process.env.FRONTEND_PUBLIC_URL),
    hostFromUrl(process.env.PUBLIC_URL),
  ]);

  for (const host of hosts) {
    args.push('--allowed-hosts', host);
  }
}

function appendDefaultFlag(args, name, enabled) {
  if (hasFlag(args, name)) {
    return;
  }

  args.push(`${name}=${enabled ? 'true' : 'false'}`);
}

function hasFlag(args, name) {
  return args.some((arg) => arg === name || arg.startsWith(`${name}=`));
}

function envBool(name, fallback) {
  const value = process.env[name];
  if (value === undefined || value === '') {
    return fallback;
  }

  return ['1', 'true', 'yes', 'on'].includes(value.toLowerCase());
}

function splitList(value) {
  if (!value) {
    return [];
  }

  return value
    .split(',')
    .map((item) => normalizeHost(item))
    .filter(Boolean);
}

function hostsFromOrigins(value) {
  return splitList(value)
    .map((origin) => hostFromUrl(origin))
    .filter(Boolean);
}

function hostFromUrl(value) {
  if (!value) {
    return '';
  }

  try {
    return normalizeHost(new URL(value).hostname);
  } catch {
    return normalizeHost(value);
  }
}

function normalizeHost(value) {
  return value.trim().replace(/^https?:\/\//, '').replace(/\/.*$/, '').replace(/:\d+$/, '');
}

function unique(values) {
  return [...new Set(values.map((value) => normalizeHost(value)).filter(Boolean))];
}
