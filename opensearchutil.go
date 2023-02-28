package opensearchutil

import (
	_ "embed"
)

const (
	DefaultMaxDepth = 2

	tagKey          = "opensearch"
	tagOptionType   = "type"
	tagOptionFormat = "format"
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

// IndexSettings allows to specify settings of an index, at its creation. This struct includes both static (those
// that are specified at index increation) settings, and dynamic settings (those that can be altered after index
// creation).
// Refer to https://opensearch.org/docs/latest/api-reference/index-apis/create-index/ for docs on each setting.
type IndexSettings struct {
	NumberOfShards                  *uint16 `json:"number_of_shards,omitempty"`
	NumberOfRoutingShards           *uint16 `json:"number_of_routing_shards,omitempty"`
	ShardCheckOnStartup             *bool   `json:"shard.check_on_startup,omitempty"`
	Codec                           *string `json:"codec,omitempty"`
	RoutingPartitionSize            *uint16 `json:"routing_partition_size,omitempty"`
	SoftDeletesRetentionLeasePeriod *string `json:"soft_deletes.retention_lease.period,omitempty"`
	LoadFixedBitsetFiltersEagerly   *bool   `json:"load_fixed_bitset_filters_eagerly,omitempty"`
	Hidden                          *bool   `json:"hidden,omitempty"`
	NumberOfReplicas                *uint16 `json:"number_of_replicas,omitempty"`
	AutoExpandReplicas              *string `json:"auto_expand_replicas,omitempty"`
	SearchIdleAfter                 *string `json:"search.idle.after,omitempty"`
	RefreshInterval                 *string `json:"refresh_interval,omitempty"`
	MaxResultWindow                 *uint64 `json:"max_result_window,omitempty"`
	MaxInnerResultWindow            *uint64 `json:"max_inner_result_window,omitempty"`
	MaxRescoreWindow                *uint64 `json:"max_rescore_window,omitempty"`
	MaxDocvalueFieldsSearch         *uint64 `json:"max_docvalue_fields_search,omitempty"`
	MaxScriptFields                 *uint16 `json:"max_script_fields,omitempty"`
	MaxNgramDiff                    *uint16 `json:"max_ngram_diff,omitempty"`
	MaxShingleDiff                  *uint16 `json:"max_shingle_diff,omitempty"`
	MaxRefreshListeners             *uint16 `json:"max_refresh_listeners,omitempty"`
	AnalyzeMaxTokenCount            *uint16 `json:"analyze.max_token_count,omitempty"`
	HighlightMaxAnalyzedOffset      *uint64 `json:"highlight.max_analyzed_offset,omitempty"`
	MaxTermsCount                   *uint64 `json:"max_terms_count,omitempty"`
	MaxRegexLength                  *uint16 `json:"max_regex_length,omitempty"`
	QueryDefaultField               *string `json:"query.default_field,omitempty"`
	RoutingAllocationEnable         *string `json:"routing.allocation_enable,omitempty"`
	RoutingRebalanceEnable          *string `json:"routing.rebalance_enable,omitempty"`
	GcDeletes                       *string `json:"gc_deletes,omitempty"`
	DefaultPipeline                 *string `json:"default_pipeline,omitempty"`
	FinalPipeline                   *string `json:"final_pipeline,omitempty"`
}

type JsonFormatter interface {
	FormatJson(str []byte) ([]byte, error)
}

type FieldNameTransformer interface {
	TransformFieldName(name string) (string, error)
}
