package opensearchutil

const DefaultMaxDepth = 2

type Option interface {
	apply(*optionsContainer)
}

type optionsContainer struct {
	maxDepth uint8
}

type maxDepthOption uint8

func (c maxDepthOption) apply(opts *optionsContainer) {
	opts.maxDepth = uint8(c)
}

func WithMaxDepth(maxDepth uint8) Option {
	return maxDepthOption(maxDepth)
}
