package opensearchutil

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/gomega"
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

func TestTimeFormat_marshallingIntoJson(t *testing.T) {
	g := NewGomegaWithT(t)

	type foo struct {
		A TimeBasicDateTime
		B TimeBasicDateTimeNoMillis
		C TimeBasicDate
	}

	someTime := time.Date(2019, 3, 23, 21, 34, 46, 567000000, time.UTC)
	jsonBytes, err := json.Marshal(foo{
		A: TimeBasicDateTime(someTime),
		B: TimeBasicDateTimeNoMillis(someTime),
		C: TimeBasicDate(someTime),
	})

	var fooUnmarshalled foo
	err = json.Unmarshal(jsonBytes, &fooUnmarshalled)
	g.Expect(err).To(BeNil())
	g.Expect(time.Time(fooUnmarshalled.A).Equal(time.Date(2019, 3, 23, 21, 34, 46, 567000000, time.UTC))).To(BeTrue())
	g.Expect(time.Time(fooUnmarshalled.B).Equal(time.Date(2019, 3, 23, 21, 34, 46, 0, time.UTC))).To(BeTrue())
	g.Expect(time.Time(fooUnmarshalled.C).Equal(time.Date(2019, 3, 23, 0, 0, 0, 0, time.UTC))).To(BeTrue())
}

func TestNumericTime_MarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	customTime := NumericTime{Time: time.Date(2023, 1, 7, 15, 0, 0, 0, time.UTC)}
	jsonBytes, err := json.Marshal(customTime)
	g.Expect(err).To(BeNil())
	g.Expect(string(jsonBytes)).To(Equal("1673103600")) // Unix timestamp for 2023-01-07 15:00:00 UTC
}

func TestNumericTime_UnmarshalJSON(t *testing.T) {
	g := NewGomegaWithT(t)

	var customTime NumericTime
	err := json.Unmarshal([]byte("1673103600"), &customTime)
	g.Expect(err).To(BeNil())
	g.Expect(customTime.Time.Equal(time.Date(2023, 1, 7, 15, 0, 0, 0, time.UTC))).To(BeTrue())
}

func TestNumericTime_JSONMarshallingAndUnmarshalling(t *testing.T) {
	g := NewGomegaWithT(t)

	type TestStruct struct {
		Timestamp NumericTime `json:"timestamp"`
	}

	// Original time for testing
	originalTime := time.Date(2023, 1, 7, 15, 0, 0, 0, time.UTC)
	testObj := TestStruct{Timestamp: NumericTime{Time: originalTime}}

	// Marshal the struct
	jsonBytes, err := json.Marshal(testObj)
	g.Expect(err).To(BeNil())
	g.Expect(string(jsonBytes)).To(Equal(`{"timestamp":1673103600}`))

	// Unmarshal back into a struct
	var unmarshalledObj TestStruct
	err = json.Unmarshal(jsonBytes, &unmarshalledObj)
	g.Expect(err).To(BeNil())
	g.Expect(unmarshalledObj.Timestamp.Time.Equal(originalTime)).To(BeTrue())
}
