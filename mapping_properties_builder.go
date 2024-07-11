package opensearchutil

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

type MappingPropertiesBuilder struct { //.
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
		fieldName := tField.Name
		resolvedField := b.resolveField(tField, v.Field(i))
		resolvedField = b.unslice(resolvedField)

		if err := validateField(resolvedField); err != nil {
			return nil, errors.Wrapf(err, "validateField")
		}

		fieldType, err := b.resolveFieldType(resolvedField)
		if err != nil {
			return nil, errors.Wrapf(err, "resolveFieldType")
		}

		fieldFormat, err := b.resolveFieldFormat(resolvedField)
		if err != nil {
			return nil, errors.Wrapf(err, "resolveFieldFormat")
		}

		transformedFieldName, err := b.optionContainer.fieldNameTransformer.TransformFieldName(fieldName)
		if err != nil {
			return nil, errors.Wrapf(err, "TransformFieldName")
		}

		if fieldType != "" {
			mappingProperty := MappingProperty{
				FieldName:   transformedFieldName,
				FieldType:   fieldType,
				FieldFormat: fieldFormat,
			}
			if err := b.addProperties(resolvedField, &mappingProperty); err != nil {
				return nil, errors.Wrapf(err, "addProperties")
			}
			mappingProperties = append(mappingProperties, mappingProperty)
			continue
		} else if resolvedField.kind == reflect.Struct {
			if nthLevel+1 > b.optionContainer.maxDepth {
				continue
			}
			children, err := b.doBuildMappingProperties(resolvedField.value.Interface(), nthLevel+1)
			if err != nil {
				return nil, errors.Wrapf(err, "nested b.doBuildMappingProperties")
			}
			mappingProperties = append(mappingProperties, MappingProperty{
				FieldName:   transformedFieldName,
				Children:    children,
				FieldFormat: fieldFormat,
			})
			continue
		} else if !b.optionContainer.omitUnsupportedTypes {
			return nil, fmt.Errorf(
				"field not supported: %s, please use opensearchutil.OmitUnsupportedTypes to skip"+
					" fields of unsupported types",
				resolvedField.field.Name)
		}
	}
	return mappingProperties, nil
}

func (b *MappingPropertiesBuilder) addProperties(resolvedField *fieldWrapper, mappingProperty *MappingProperty) error {
	indexPrefixes := getTagOptionValue(resolvedField.field, tagKey, tagOptionIndexPrefixes)
	if indexPrefixes != "" {
		opts := parseCustomPropertyValue(indexPrefixes)
		mappingProperty.IndexPrefixes = MakePtr(make(map[string]string, len(opts)))
		for k, v := range opts {
			(*mappingProperty.IndexPrefixes)[k] = v
		}
	}

	analyzer := getTagOptionValue(resolvedField.field, tagKey, tagOptionAnalyzer)
	if analyzer != "" {
		mappingProperty.Analyzer = MakePtr(analyzer)
	}

	return nil
}

func validateField(field *fieldWrapper) error {
	if field.kind == reflect.Struct {
		switch field.value.Interface().(type) {
		case time.Time:
			return ErrGotBuiltInTimeField
		}
	}
	return nil
}

func (b *MappingPropertiesBuilder) resolveFieldType(field *fieldWrapper) (string, error) {
	fieldTypeOverride := getTagOptionValue(field.field, tagKey, tagOptionType)
	if fieldTypeOverride != "" {
		return fieldTypeOverride, nil
	}
	if field.isPrimitive {
		return b.getDefaultOSTypeFromPrimitiveKind(field.kind), nil
	}
	if field.kind == reflect.Struct {
		if x, ok := field.value.Interface().(OpenSearchDateType); ok {
			if x.GetOpenSearchFieldType() != "" {
				return "date", nil
			}
		}
	}
	return "", nil
}

func (b *MappingPropertiesBuilder) resolveFieldFormat(field *fieldWrapper) (*string, error) {
	fieldFormatOverride := getTagOptionValue(field.field, tagKey, tagOptionFormat)
	if fieldFormatOverride != "" {
		return &fieldFormatOverride, nil
	}
	if field.kind == reflect.Struct {
		if x, ok := field.value.Interface().(OpenSearchDateType); ok {
			if x.GetOpenSearchFieldType() != "" {
				return MakePtr(x.GetOpenSearchFieldType()), nil
			}
		}
	}
	return nil, nil
}

// resolveField returns a wrapper object for the given field. If the field is a pointer, it returns a wrapper
// for the dereferenced field, since we treat both pointer and value fields the same.
func (b *MappingPropertiesBuilder) resolveField(structField reflect.StructField, value reflect.Value) *fieldWrapper {
	var kind reflect.Kind
	var val reflect.Value
	if structField.Type.Kind() == reflect.Ptr {
		kind = structField.Type.Elem().Kind()
		val = reflect.New(structField.Type.Elem()).Elem()
	} else {
		kind = structField.Type.Kind()
		val = value
	}

	return &fieldWrapper{
		field:       structField,
		kind:        kind,
		value:       val,
		isPrimitive: b.isPrimitive(kind),
	}
}

// unslice "un-slices" the field by resolving the underlying element type. I.e. if it's a slice of struct Foo
// the returned object contains reflection objects for just Foo (not its slice).
func (b *MappingPropertiesBuilder) unslice(wrapper *fieldWrapper) *fieldWrapper {
	if wrapper.kind != reflect.Slice {
		return wrapper
	}

	elemType := wrapper.field.Type.Elem()
	var (
		newKind reflect.Kind
		newVal  reflect.Value
	)
	if elemType.Kind() == reflect.Ptr {
		newKind = elemType.Elem().Kind()
		newVal = reflect.New(elemType.Elem()).Elem()
	} else {
		newKind = elemType.Kind()
		newVal = reflect.New(elemType).Elem()
	}

	return &fieldWrapper{
		field:       wrapper.field,
		kind:        newKind,
		value:       newVal,
		isPrimitive: b.isPrimitive(newKind),
	}
}

func (b *MappingPropertiesBuilder) isPrimitive(kind reflect.Kind) bool {
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
