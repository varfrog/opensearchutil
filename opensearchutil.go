package main

import (
	"bytes"
	_ "embed"
	"github.com/pkg/errors"
	"reflect"
	"text/template"
)

type mappingProperty struct {
	FieldName string
	FieldType string
}

type indexTplData struct {
	MappingProperties []mappingProperty
}

//go:embed index_tpl.txt
var indexTpl string

func GenerateIndexJson(obj interface{}) (string, error) {
	mappingProperties, err := getMappingProperties(obj)
	if err != nil {
		return "", errors.Wrapf(err, "getMappingProperties")
	}

	tpl, err := template.New("IndexTpl").Parse(indexTpl)
	if err != nil {
		return "", errors.Wrapf(err, "parse template")
	}

	var tplResult bytes.Buffer
	err = tpl.Execute(&tplResult, indexTplData{
		MappingProperties: mappingProperties,
	})
	if err != nil {
		return "", errors.Wrapf(err, "tpl.Execute")
	}

	return tplResult.String(), nil
}

func getMappingProperties(obj interface{}) ([]mappingProperty, error) {
	var mappingProperties []mappingProperty
	t := reflect.ValueOf(obj).Type()
	for i := 0; i < t.NumField(); i++ {
		kind := t.Field(i).Type.Kind()
		resolvedKind := kind
		if kind == reflect.Ptr {
			resolvedKind = t.Field(i).Type.Elem().Kind()
		}
		if isPrimitiveNonPtr(resolvedKind) {
			mappingProperties = append(mappingProperties, mappingProperty{
				FieldName: toSnakeCase(t.Field(i).Name),
				FieldType: getOpenSearchType(resolvedKind),
			})
		}
	}
	return mappingProperties, nil
}

func isPrimitiveNonPtr(kind reflect.Kind) bool {
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
		return true
	default:
		return false
	}
}

func getOpenSearchType(kind reflect.Kind) string {
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
		return "text" // can also be keyword, token_count
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
