package opensearchutil

import (
	_ "embed"
)

const (
	tagKey          = "opensearch"
	tagOptionType   = "type"
	tagOptionFormat = "format"

	DefaultTimeFormat = "basic_date_time"
	DefaultMaxDepth   = 2
)

// MappingProperty corresponds to mappings.properties of a mapping JSON. See
// https://opensearch.org/docs/1.3/opensearch/mappings/#explicit-mapping.
// MappingProperty defines either a primitive data type, in which case FieldType != "", or an object, in which case
// len(Children) > 0.
type MappingProperty struct {
	FieldName   string
	FieldType   string
	FieldFormat *string
	Children    []MappingProperty
}

type JsonFormatter interface {
	FormatJson(str []byte) ([]byte, error)
}

type FieldNameTransformer interface {
	TransformFieldName(name string) (string, error)
}
