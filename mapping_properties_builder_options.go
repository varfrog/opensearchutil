package opensearchutil

type MappingPropertiesBuilderOption interface {
	apply(*mappingPropertiesBuilderOptionContainer)
}

type mappingPropertiesBuilderOptionContainer struct {
	maxDepth             uint8
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
