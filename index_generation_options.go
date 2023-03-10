package opensearchutil

type IndexGenerationOption interface {
	apply(*indexGenerationOptionContainer)
}

type indexGenerationOptionContainer struct {
	strictMapping bool
}

// Strict mapping

type strictMappingOption bool

func (c strictMappingOption) apply(opts *indexGenerationOptionContainer) {
	opts.strictMapping = bool(c)
}

// WithStrictMapping adds "dynamic: "strict" to "mappings"
func WithStrictMapping(strictMapping bool) IndexGenerationOption {
	return strictMappingOption(strictMapping)
}
