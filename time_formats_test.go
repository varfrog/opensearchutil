package opensearchutil

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestTimeBasicDateTime_MarshalText(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := TimeBasicDateTime(time.Date(2019, 3, 23, 21, 34, 46, 567000000, time.UTC))
	res, err := customTime.MarshalText()
	g.Expect(err).To(BeNil())

	g.Expect(string(res)).To(Equal("20190323T213446.567+00:00"))
}

func TestTimeBasicDateTime_UnmarshalText(t *testing.T) {
	g := NewGomegaWithT(t)

	var obj TimeBasicDateTime

	g.Expect(obj.UnmarshalText([]byte("20190323T213446.567+00:00"))).To(Succeed())
	g.Expect(time.Time(obj).Equal(time.Date(2019, 3, 23, 21, 34, 46, 567000000, time.UTC))).To(BeTrue())
}

func TestTimeBasicDateTimeNoMillis_MarshalText(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := TimeBasicDateTimeNoMillis(time.Date(2019, 3, 23, 21, 34, 46, 0, time.UTC))
	res, err := customTime.MarshalText()
	g.Expect(err).To(BeNil())

	g.Expect(string(res)).To(Equal("20190323T213446+00:00"))
}

func TestTimeBasicDateTimeNoMillis_UnmarshalText(t *testing.T) {
	g := NewGomegaWithT(t)

	var obj TimeBasicDateTimeNoMillis

	g.Expect(obj.UnmarshalText([]byte("20190323T213446+00:00"))).To(Succeed())
	g.Expect(time.Time(obj).Equal(time.Date(2019, 3, 23, 21, 34, 46, 0, time.UTC))).To(BeTrue())
}

func TestTimeBasicDate_MarshalText(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := TimeBasicDate(time.Date(2019, 3, 23, 21, 34, 46, 0, time.UTC))
	res, err := customTime.MarshalText()
	g.Expect(err).To(BeNil())

	g.Expect(string(res)).To(Equal("20190323"))
}

func TestTimeBasicDate_UnmarshalText(t *testing.T) {
	g := NewGomegaWithT(t)

	var obj TimeBasicDate

	g.Expect(obj.UnmarshalText([]byte("20190323"))).To(Succeed())
	g.Expect(time.Time(obj).Equal(time.Date(2019, 3, 23, 0, 0, 0, 0, time.UTC))).To(BeTrue())
}
