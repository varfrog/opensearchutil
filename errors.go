package opensearchutil

import (
	"errors"
)

var ErrGotBuiltInTimeField = errors.New(`time.Time fields cannot be used, use Time* types or custom types that implement encoding.TextMarshaler and opensearchutil.OpenSearchTime and marshall into OpenSearch date formats`)
