package opensearchutil

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"time"
)

type MappingPropertiesBuilder struct {
	optionContainer mappingPropertiesBuilderOptionContainer
}

type fieldWrapper struct {
	field       reflect.StructField
	kind        reflect.Kind
	value       reflect.Value
	isPrimitive bool
}

func NewMappingPropertiesBuilder(options ...MappingPropertiesBuilderOption) *MappingPropertiesBuilder {
	optContainer := mappingPropertiesBuilderOptionContainer{
		maxDepth: DefaultMaxDepth,
	}
	for _, o := range options {
		o.apply(&optContainer)
	}
	if optContainer.fieldNameTransformer == nil {
		optContainer.fieldNameTransformer = NewSnakeCaser()
	}
	if optContainer.jsonFormatter == nil {
		optContainer.jsonFormatter = NewMarshalIndentJsonFormatter()
	}

	return &MappingPropertiesBuilder{optionContainer: optContainer}
}

func (b *MappingPropertiesBuilder) BuildMappingProperties(obj interface{}) ([]MappingProperty, error) {
	mps, err := b.doBuildMappingProperties(obj, 1)
	if err != nil {
		return nil, errors.Wrapf(err, "b.doBuildMappingProperties")
	}
	return mps, nil
}

func (b *MappingPropertiesBuilder) doBuildMappingProperties(
	obj interface{},
	nthLevel uint8,
) ([]MappingProperty, error) {
	var mappingProperties []MappingProperty
	v := reflect.ValueOf(obj)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		resolvedField := b.resolveField(tField, v.Field(i))
		fieldType, err := b.resolveFieldType(resolvedField)
		if err != nil {
			return nil, errors.Wrapf(err, "resolveFieldType")
		}

		transformedFieldName, err := b.optionContainer.fieldNameTransformer.TransformFieldName(tField.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "TransformFieldName")
		}

		if fieldType != "" {
			mappingProperties = append(mappingProperties, MappingProperty{
				FieldName: transformedFieldName,
				FieldType: fieldType,
			})
			continue
		}

		if resolvedField.kind == reflect.Struct && nthLevel+1 <= b.optionContainer.maxDepth {
			children, err := b.doBuildMappingProperties(resolvedField.value.Interface(), nthLevel+1)
			if err != nil {
				return nil, errors.Wrapf(err, "nested b.doBuildMappingProperties")
			}
			mappingProperties = append(mappingProperties, MappingProperty{
				FieldName: transformedFieldName,
				Children:  children,
			})
			continue
		}
	}
	return mappingProperties, nil
}

func (b *MappingPropertiesBuilder) resolveFieldType(field fieldWrapper) (string, error) {
	fieldTypeOverride := b.getFieldTypeOverride(field.field)
	if fieldTypeOverride != "" {
		return fieldTypeOverride, nil
	}
	if field.isPrimitive {
		return b.getDefaultOSTypeFromPrimitiveKind(field.kind), nil
	}
	if field.kind == reflect.Struct {
		switch field.value.Interface().(type) {
		case time.Time:
			return defaultTimeType, nil
		}
	}
	return "", nil
}

// getFieldTypeOverride returns a type of the given field if it is overriden by a tag,
// returns "" if it is not overriden.
func (b *MappingPropertiesBuilder) getFieldTypeOverride(structField reflect.StructField) string {
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

// resolveField returns a wrapper object for the given field. If the field is a pointer, it returns a wrapper
// for the dereferences field, since we treat both pointer and value fields the same.
func (b *MappingPropertiesBuilder) resolveField(structField reflect.StructField, value reflect.Value) fieldWrapper {
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
		field:       structField,
		kind:        kind,
		value:       val,
		isPrimitive: primitive,
	}
}

func (b *MappingPropertiesBuilder) getDefaultOSTypeFromPrimitiveKind(kind reflect.Kind) string {
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
