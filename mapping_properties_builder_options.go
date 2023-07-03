package opensearchutil

type MappingPropertiesBuilderOption interface {
	apply(*mappingPropertiesBuilderOptionContainer)
}

type mappingPropertiesBuilderOptionContainer struct {
	maxDepth             uint8
	omitUnsupportedTypes bool
	fieldNameTransformer FieldNameTransformer
	jsonFormatter        JsonFormatter
}

// MaxDepth option
type maxDepthOption uint8

func (c maxDepthOption) apply(opts *mappingPropertiesBuilderOptionContainer) {
	opts.maxDepth = uint8(c)
}

func WithMaxDepth(maxDepth uint8) MappingPropertiesBuilderOption {
	return maxDepthOption(maxDepth)
}

// FieldNameTransformer option
type fieldNameTransformerOption struct {
	fieldNameTransformer FieldNameTransformer
}

func (c fieldNameTransformerOption) apply(opts *mappingPropertiesBuilderOptionContainer) {
	opts.fieldNameTransformer = c.fieldNameTransformer
}

//goland:noinspection GoUnusedExportedFunction
func WithFieldNameTransformer(fieldNameTransformer FieldNameTransformer) MappingPropertiesBuilderOption {
	return fieldNameTransformerOption{fieldNameTransformer: fieldNameTransformer}
}

// MaxDepth option
type skipUnsupportedTypesOption bool

func (c skipUnsupportedTypesOption) apply(opts *mappingPropertiesBuilderOptionContainer) {
	opts.omitUnsupportedTypes = bool(c)
}

// OmitUnsupportedTypes makes the builder not complain about types it cannot support. Those would then need to be
// generated manually.
func OmitUnsupportedTypes() MappingPropertiesBuilderOption {
	return skipUnsupportedTypesOption(true)
}
