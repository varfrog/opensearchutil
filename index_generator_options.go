package opensearchutil

type IndexGeneratorOption interface {
	apply(*indexGeneratorOptionContainer)
}

type indexGeneratorOptionContainer struct {
	jsonFormatter JsonFormatter
}

type jsonFormatterOption struct {
	jsonFormatter JsonFormatter
}

func (c jsonFormatterOption) apply(opts *indexGeneratorOptionContainer) {
	opts.jsonFormatter = c.jsonFormatter
}

//goland:noinspection GoUnusedExportedFunction
func WithJsonFormatter(jsonFormatter JsonFormatter) IndexGeneratorOption {
	return jsonFormatterOption{jsonFormatter: jsonFormatter}
}
