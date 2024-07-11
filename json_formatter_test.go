package opensearchutil

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestMarshalIndentJsonFormatter_FormatJson(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	inputJson := `{"name":"tom",   "age":

80,"foo":{}}`

	expectedJson := `{
   "age": 80,
   "foo": {},
   "name": "tom"
}`

	res, err := NewMarshalIndentJsonFormatter().FormatJson([]byte(inputJson))
	g.Expect(err).To(gomega.BeNil())
	g.Expect(string(res)).To(gomega.Equal(expectedJson))
}
