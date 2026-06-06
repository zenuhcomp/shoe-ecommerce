#!/usr/bin/env bash
# Start Luna Collections backend (Go) and frontend (Angular) using backend/.env
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${ENV_FILE:-$ROOT/backend/.env}"
BACKEND_DIR="$ROOT/backend"
FRONTEND_DIR="$ROOT/frontend"
RUN_DIR="$ROOT/.run"
BACKEND_BIN="${BACKEND_BIN:-$BACKEND_DIR/backend}"
PROXY_FILE="$FRONTEND_DIR/proxy.conf.runtime.json"
PID_BACKEND="$RUN_DIR/backend.pid"
PID_FRONTEND="$RUN_DIR/frontend.pid"
LOG_BACKEND="$RUN_DIR/backend.log"
LOG_FRONTEND="$RUN_DIR/frontend.log"

mkdir -p "$RUN_DIR" "$BACKEND_DIR/data"

if [[ ! -f "$ENV_FILE" ]]; then
  echo "No $ENV_FILE found — creating from backend/.env.example"
  cp "$BACKEND_DIR/.env.example" "$ENV_FILE"
fi

# shellcheck disable=SC1090
set -a
source "$ENV_FILE"
set +a

FRONTEND_PORT="${FRONTEND_PORT:-${frontend_port:-4204}}"
BACKEND_PORT="${BACKEND_PORT:-${backend_port:-8084}}"
DATABASE_MODE="${DATABASE_MODE:-${database_mode:-test}}"
FRONTEND_DIST_DIR="${FRONTEND_DIST_DIR:-$FRONTEND_DIR/dist/frontend/browser}"
SERVE_STATIC="${SERVE_STATIC:-${serve_static:-}}"
if [[ -z "$SERVE_STATIC" ]]; then
  SERVE_STATIC=0
  if [[ -f "$FRONTEND_DIST_DIR/index.html" && ! -x "$FRONTEND_DIR/node_modules/.bin/ng" ]]; then
    SERVE_STATIC=1
  fi
fi

export PORT="$BACKEND_PORT"
export BACKEND_PORT
export FRONTEND_PORT
export DATABASE_MODE
export ENV_FILE
export FRONTEND_DIST_DIR

if [[ -z "${CORS_ORIGINS:-}" ]]; then
  export CORS_ORIGINS="http://localhost:${FRONTEND_PORT},http://127.0.0.1:${FRONTEND_PORT}"
fi

stop_process() {
  local pid_file="$1"
  local label="$2"
  if [[ -f "$pid_file" ]]; then
    local pid
    pid="$(cat "$pid_file")"
    if kill -0 "$pid" 2>/dev/null; then
      echo "Stopping $label (pid $pid)..."
      kill "$pid" 2>/dev/null || true
      wait "$pid" 2>/dev/null || true
    fi
    rm -f "$pid_file"
  fi
}

cleanup() {
  echo ""
  echo "Shutting down..."
  stop_process "$PID_FRONTEND" "frontend"
  stop_process "$PID_BACKEND" "backend"
  exit 0
}

trap cleanup INT TERM

stop_process "$PID_FRONTEND" "frontend"
stop_process "$PID_BACKEND" "backend"

cat > "$PROXY_FILE" <<EOF
{
  "/api": {
    "target": "http://localhost:${BACKEND_PORT}",
    "secure": false
  }
}
EOF

echo "=============================================="
echo " Luna Collections — dev stack"
echo "=============================================="
echo " Env file       : $ENV_FILE"
echo " DATABASE_MODE  : $DATABASE_MODE"
echo " Backend        : http://localhost:${BACKEND_PORT}"
if [[ "$SERVE_STATIC" == "1" ]]; then
  echo " Frontend       : http://localhost:${BACKEND_PORT} (static)"
  echo " Frontend dist  : $FRONTEND_DIST_DIR"
else
  echo " Frontend       : http://localhost:${FRONTEND_PORT}"
fi
if [[ "$SERVE_STATIC" != "1" ]]; then
  echo " API proxy      : $PROXY_FILE"
fi
echo "=============================================="
echo ""

echo "Starting backend..."
if [[ -x "$BACKEND_BIN" ]]; then
  (
    cd "$BACKEND_DIR"
    exec "$BACKEND_BIN" >>"$LOG_BACKEND" 2>&1
  ) &
else
  (
    cd "$BACKEND_DIR"
    exec go run . >>"$LOG_BACKEND" 2>&1
  ) &
fi
echo $! >"$PID_BACKEND"

echo "Waiting for API health..."
for _ in $(seq 1 40); do
  if curl -sf "http://localhost:${BACKEND_PORT}/api/health" >/dev/null 2>&1; then
    break
  fi
  if ! kill -0 "$(cat "$PID_BACKEND")" 2>/dev/null; then
    break
  fi
  sleep 0.5
done

if ! curl -sf "http://localhost:${BACKEND_PORT}/api/health" >/dev/null 2>&1; then
  echo "Backend failed to start. Last log lines:"
  tail -n 30 "$LOG_BACKEND" || true
  stop_process "$PID_BACKEND" "backend"
  exit 1
fi

health="$(curl -sf "http://localhost:${BACKEND_PORT}/api/health")"
echo "Backend ready: $health"
echo ""

if [[ "$SERVE_STATIC" == "1" ]]; then
  if [[ ! -f "$FRONTEND_DIST_DIR/index.html" ]]; then
    echo "Static frontend build is missing: $FRONTEND_DIST_DIR/index.html"
    echo "Run npm run build in $FRONTEND_DIR before starting with SERVE_STATIC=1."
    stop_process "$PID_BACKEND" "backend"
    exit 1
  fi

  echo "Static app running. Press Ctrl+C to stop."
  echo "  App:    http://localhost:${BACKEND_PORT}"
  echo "  API:    http://localhost:${BACKEND_PORT}/api/health"
  echo "  Logs:   $LOG_BACKEND"
  echo ""

  while true; do
    if ! kill -0 "$(cat "$PID_BACKEND")" 2>/dev/null; then
      echo "Backend exited. Last log lines:"
      tail -n 50 "$LOG_BACKEND" || true
      exit 1
    fi
    sleep 1
  done
fi

echo "Starting frontend (ng serve)..."
if [[ ! -x "$FRONTEND_DIR/node_modules/.bin/ng" ]]; then
  echo "Frontend dependencies are missing. Run npm install in $FRONTEND_DIR."
  stop_process "$PID_BACKEND" "backend"
  exit 1
fi

(
  cd "$FRONTEND_DIR"
  exec npm start -- --port "$FRONTEND_PORT" --proxy-config proxy.conf.runtime.json >>"$LOG_FRONTEND" 2>&1
) &
echo $! >"$PID_FRONTEND"

echo "Waiting for frontend..."
for _ in $(seq 1 80); do
  if curl -sf "http://localhost:${FRONTEND_PORT}" >/dev/null 2>&1; then
    break
  fi
  if ! kill -0 "$(cat "$PID_FRONTEND")" 2>/dev/null; then
    break
  fi
  sleep 0.5
done

if ! curl -sf "http://localhost:${FRONTEND_PORT}" >/dev/null 2>&1; then
  echo "Frontend failed to start. Last log lines:"
  tail -n 50 "$LOG_FRONTEND" || true
  stop_process "$PID_FRONTEND" "frontend"
  stop_process "$PID_BACKEND" "backend"
  exit 1
fi

echo ""
echo "Stack running. Press Ctrl+C to stop."
echo "  App:    http://localhost:${FRONTEND_PORT}"
echo "  API:    http://localhost:${BACKEND_PORT}/api/health"
echo "  Logs:   $LOG_BACKEND"
echo "          $LOG_FRONTEND"
echo ""

while true; do
  if ! kill -0 "$(cat "$PID_BACKEND")" 2>/dev/null; then
    echo "Backend exited. Last log lines:"
    tail -n 50 "$LOG_BACKEND" || true
    stop_process "$PID_FRONTEND" "frontend"
    exit 1
  fi
  if ! kill -0 "$(cat "$PID_FRONTEND")" 2>/dev/null; then
    echo "Frontend exited. Last log lines:"
    tail -n 50 "$LOG_FRONTEND" || true
    stop_process "$PID_BACKEND" "backend"
    exit 1
  fi
  sleep 1
done
