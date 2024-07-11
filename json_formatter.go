package opensearchutil

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type MarshalIndentJsonFormatter struct{}

func NewMarshalIndentJsonFormatter() *MarshalIndentJsonFormatter {
	return &MarshalIndentJsonFormatter{}
}

func (f *MarshalIndentJsonFormatter) FormatJson(str []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(str, &obj); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal")
	}

	jsonBytes, err := json.MarshalIndent(&obj, "", "   ")
	if err != nil {
		return nil, errors.Wrapf(err, "json.Marshal")
	}
	return jsonBytes, nil
}
