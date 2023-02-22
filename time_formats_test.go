package opensearchutil

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestTimeBasicDateTime_MarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := TimeBasicDateTime(time.Date(2019, 3, 23, 21, 34, 46, 567000000, time.UTC))
	res, err := customTime.MarshalText()
	g.Expect(err).To(BeNil())

	g.Expect(string(res)).To(Equal("20190323T213446.567+00:00"))
}

func TestTimeBasicDateTimeNoMillis_MarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := TimeBasicDateTimeNoMillis(time.Date(2019, 3, 23, 21, 34, 46, 0, time.UTC))
	res, err := customTime.MarshalText()
	g.Expect(err).To(BeNil())

	g.Expect(string(res)).To(Equal("20190323T213446+00:00"))
}

func TestTimeBasicDate_MarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := TimeBasicDate(time.Date(2019, 3, 23, 21, 34, 46, 0, time.UTC))
	res, err := customTime.MarshalText()
	g.Expect(err).To(BeNil())

	g.Expect(string(res)).To(Equal("20190323"))
}
