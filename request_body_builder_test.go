package opensearchutil

import (
	"github.com/onsi/gomega"
	"testing"
)

func TestRequestBodyBuilder_BuildIndexBody(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	builder := NewRequestBodyBuilder()

	type address struct {
		PostalCode uint32 `json:"postal_code"`
	}
	type person struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Age     uint8   `json:"age"`
		Address address `json:"address"`
	}

	resBody, err := builder.BuildIndexBody([]ObjectWrapper{
		{
			ID:    "680",
			Index: "people",
			Object: person{
				ID:      "680",
				Name:    "Ann",
				Age:     40,
				Address: address{PostalCode: 10000},
			},
		},
		{
			ID:    "730",
			Index: "people",
			Object: person{
				ID:      "730",
				Name:    "Bob",
				Age:     38,
				Address: address{PostalCode: 35000},
			},
		},
	})
	g.Expect(err).To(gomega.BeNil())

	expBody := `{"index":{"_index":"people","_id":"680"}}
{"id":"680","name":"Ann","age":40,"address":{"postal_code":10000}}
{"index":{"_index":"people","_id":"730"}}
{"id":"730","name":"Bob","age":38,"address":{"postal_code":35000}}
`
	g.Expect(resBody).To(gomega.Equal(expBody))
}
