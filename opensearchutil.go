package opensearchutil

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"text/template"
)

const (
	tagKey     = "opensearch"
	tagKeyType = "type"
)

//go:embed index.gotmpl
var indexTmpl string

type fieldWrapper struct {
	kind        reflect.Kind
	isPrimitive bool
	value       reflect.Value
}

// MappingProperty corresponds to mappings.properties of a mapping JSON. See
// https://opensearch.org/docs/1.3/opensearch/mappings/#explicit-mapping.
// MappingProperty defines either a primitive data type, in which case FieldType != "", or an object, in which case
// len(Children) > 0.
type MappingProperty struct {
	FieldName string

	FieldType string
	Children  []MappingProperty
}

func BuildMappingProperties(obj interface{}) ([]MappingProperty, error) {
	var mappingProperties []MappingProperty
	v := reflect.ValueOf(obj)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		resolvedField := resolveField(tField, v.Field(i))
		fieldTypeOverride := getFieldTypeOverride(tField)

		if resolvedField.isPrimitive {
			var fieldType string
			if fieldTypeOverride != "" {
				fieldType = fieldTypeOverride
			} else {
				fieldType = getDefaultOSTypeFromPrimitiveKind(resolvedField.kind)
			}
			mappingProperties = append(mappingProperties, MappingProperty{
				FieldName: toSnakeCase(tField.Name),
				FieldType: fieldType,
			})
		} else if resolvedField.kind == reflect.Struct {
			if fieldTypeOverride != "" {
				mappingProperties = append(mappingProperties, MappingProperty{
					FieldName: toSnakeCase(tField.Name),
					FieldType: fieldTypeOverride,
				})
			} else {
				mp, err := BuildMappingProperties(resolvedField.value.Interface())
				if err != nil {
					return nil, errors.Wrapf(err, "nested BuildMappingProperties")
				}
				mappingProperties = append(mappingProperties, MappingProperty{
					FieldName: toSnakeCase(tField.Name),
					Children:  mp,
				})
			}
		}
	}
	return mappingProperties, nil
}

// getFieldTypeOverride returns a type of the given field if it is overriden by a tag,
// returns "" if it is not overriden.
func getFieldTypeOverride(structField reflect.StructField) string {
	if tag := structField.Tag.Get(tagKey); tag != "" {
		for _, kvs := range strings.Split(tag, ",") {
			kv := strings.Split(kvs, ":")
			if len(kv) == 2 {
				if kv[0] == tagKeyType {
					return kv[1]
				}
			}
		}
	}
	return ""
}

// resolveField returns the kind of the type. If it's a pointer, it returns the referenced type's kind.
func resolveField(structField reflect.StructField, value reflect.Value) fieldWrapper {
	var kind reflect.Kind
	var val reflect.Value
	if structField.Type.Kind() == reflect.Ptr {
		kind = structField.Type.Elem().Kind()
		val = reflect.New(structField.Type.Elem()).Elem()
	} else {
		kind = structField.Type.Kind()
		val = value
	}

	var primitive bool
	switch kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		primitive = true
	}

	return fieldWrapper{
		kind:        kind,
		isPrimitive: primitive,
		value:       val,
	}
}

func GenerateIndexJson(mappingProperties []MappingProperty) ([]byte, error) {
	type indexTmplData struct {
		MappingProperties []MappingProperty
	}

	var funcMap template.FuncMap = map[string]interface{}{
		"notLast": func(index int, len int) bool {
			return index+1 < len
		},
	}

	tmpl, err := template.New("IndexTmpl").Funcs(funcMap).Parse(indexTmpl)
	if err != nil {
		return nil, errors.Wrapf(err, "parse template")
	}

	var tmplResult bytes.Buffer
	err = tmpl.Execute(&tmplResult, indexTmplData{
		MappingProperties: mappingProperties,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "tpl.Execute")
	}

	formattedJson, err := formatJson(tmplResult.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "formatJson")
	}

	return formattedJson, nil
}

func formatJson(str []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(str, &obj); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal")
	}

	jsonBytes, err := json.MarshalIndent(&obj, "", "   ")
	if err != nil {
		return nil, errors.Wrapf(err, "json.Marshal")
	}
	return jsonBytes, nil
}

func getDefaultOSTypeFromPrimitiveKind(kind reflect.Kind) string {
	switch kind {
	case reflect.Bool:
		return "boolean"
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return "integer"
	case reflect.Float32,
		reflect.Float64:
		return "float"
	case reflect.String:
		return "text"
	default:
		return ""
	}
}

func toSnakeCase(name string) string {
	nameBytes := []byte(name)
	nameLen := len(nameBytes)

	var result []byte
	for i := 0; i < nameLen; i++ {
		c := nameBytes[i]
		if c >= 'A' && c <= 'Z' {
			if len(result) > 0 {
				if !(i > 0 && nameBytes[i-1] >= 'A' && nameBytes[i-1] <= 'Z') { // Previous character uppercase?
					result = append(result, '_')
				}
			}
			result = append(result, c+32)
		} else if c >= '0' && c <= '9' {
			if len(result) > 0 {
				if !(i > 0 && nameBytes[i-1] >= '0' && nameBytes[i-1] <= '9') { // Previous character number?
					result = append(result, '_')
				}
			}
			result = append(result, c)
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}
