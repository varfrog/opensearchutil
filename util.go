package opensearchutil

import (
	"reflect"
	"strings"
)

func MakePtr[V any](v V) *V {
	return &v
}

func getTagOptionValue(structField reflect.StructField, tagKey string, optionKey string) string {
	if tag := structField.Tag.Get(tagKey); tag != "" {
		for _, kvs := range strings.Split(tag, ",") {
			kv := strings.Split(kvs, ":")
			if len(kv) == 2 {
				if kv[0] == optionKey {
					return kv[1]
				}
			}
		}
	}
	return ""
}
