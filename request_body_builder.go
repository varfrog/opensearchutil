package opensearchutil

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type ObjectWrapper struct {
	// ID represents the _id property of a document
	ID string

	Index  string
	Object interface{}
}

type RequestBodyBuilder struct{}

func NewRequestBodyBuilder() *RequestBodyBuilder {
	return &RequestBodyBuilder{}
}

// BuildIndexBody builds the body of a request to the POST _bulk API.
func (r *RequestBodyBuilder) BuildIndexBody(objects []ObjectWrapper) (string, error) {
	bodyParts := make([]string, 0, len(objects)*2) // 2 JSONs per doc
	for _, obj := range objects {
		bodyParts = append(
			bodyParts,
			fmt.Sprintf(`{"index":{"_index":"%s","_id":"%s"}}`, obj.Index, obj.ID))

		jsonBytes, err := json.Marshal(obj.Object)
		if err != nil {
			return "", errors.Wrapf(err, "json.Marshal")
		}
		bodyParts = append(bodyParts, string(jsonBytes))
	}

	return strings.Join(bodyParts, "\n") + "\n", nil
}
