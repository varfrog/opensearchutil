package opensearchutil

import (
	"reflect"
	"strings"
)

func MakePtr[V any](v V) *V {
	return &v
}

// getTagOptionValue gets a tag option value. For example, given a tag "type:keyword", getTagOptionValue("type")
// returns "keyword".
func getTagOptionValue(structField reflect.StructField, tagKey string, optionKey string) string {
	const tagOptionSep = ","
	const keyValSep = ":"
	if tag := structField.Tag.Get(tagKey); tag != "" {
		for _, kvs := range strings.Split(tag, tagOptionSep) {
			kv := strings.Split(kvs, keyValSep)
			if len(kv) > 1 {
				if strings.Trim(kv[0], " ") == optionKey {
					return strings.Join(kv[1:], keyValSep)
				}
			}
		}
	}
	return ""
}

// parseCustomPropertyValue parses a string like "min_chars=2;foo=bar" into a map like
// map[string]string{"min_chars":"2", "foo":"bar"}.
func parseCustomPropertyValue(str string) map[string]string {
	const keyValSep = "="
	pairs := strings.Split(str, ";")
	m := make(map[string]string, len(pairs))
	for _, kvs := range pairs {
		kv := strings.Split(kvs, keyValSep)
		if len(kv) > 1 {
			val := strings.Join(kv[1:], keyValSep)
			m[kv[0]] = val
		}
	}
	return m
}

// parseListPropertyValue parses a comma-separated list like "a,b,c" into []string{"a","b","c"}.
func parseListPropertyValue(str string) []string {
	if strings.TrimSpace(str) == "" {
		return nil
	}
	parts := strings.Split(str, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
