package repository

import (
	"encoding/json"
)

func marshalJSONStrings(items []string) string {
	if items == nil {
		items = []string{}
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(raw)
}

func unmarshalJSONStrings(raw string) []string {
	if raw == "" {
		return []string{}
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []string{}
	}
	if items == nil {
		return []string{}
	}
	return items
}

func marshalJSONMap(m map[string]string) string {
	if m == nil {
		m = map[string]string{}
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(raw)
}

func unmarshalJSONMap(raw string) map[string]string {
	if raw == "" {
		return map[string]string{}
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return map[string]string{}
	}
	if m == nil {
		return map[string]string{}
	}
	return m
}
