package opensearchutil

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

const (
	FormatTimeBasicDateTime         = "20060102T150405.999-07:00"
	FormatTimeBasicDateTimeNoMillis = "20060102T150405-07:00"
	FormatTimeBasicDate             = "20060102"
)

type (
	// TimeBasicDateTime marshalls into OpenSearch basic_date_time type
	TimeBasicDateTime time.Time

	// TimeBasicDateTimeNoMillis marshalls into OpenSearch basic_date_time_no_millis type
	TimeBasicDateTimeNoMillis time.Time

	// TimeBasicDate marshalls into OpenSearch basic_date type
	TimeBasicDate time.Time

	// NumericTime marshals to and from Unix timestamps (integer "long" values) to make sorting on dates possible.
	//
	// OpenSearch supports "date" fields for indexing and querying date-time values. However, using "date" fields for sorting
	// introduces significant performance and memory challenges because OpenSearch does not natively optimize sorting on "date"
	// fields. Sorting on "date" fields often requires loading field data into memory, which can lead to increased memory usage,
	// slower queries, and potential errors.
	//
	// To address this, NumericTime represents date-time values as Unix timestamps (in seconds since the epoch) and stores them
	// as "long" fields in OpenSearch. This approach has several advantages:
	//   - **Efficient Sorting:** Numeric fields (e.g., "long") are natively optimized for sorting in OpenSearch, avoiding
	//     the overhead of "date" field sorting.
	//   - **Compact Storage:** Unix timestamps are smaller and simpler to store compared to string-based date formats.
	//   - **Compatibility:** Unix timestamps are a widely used and supported standard for representing date-time values.
	//
	// This type is designed to work alongside "date" fields when necessary. For example, you can use the "date" field for
	// full-text queries and range filters, while relying on NumericTime for efficient sorting operations.
	NumericTime struct{ time.Time }
)

// OpenSearchDateType tells MappingPropertiesBuilder that a type is a "date" OpenSearch type.
type OpenSearchDateType interface {
	GetOpenSearchFieldType() string
}

//goland:noinspection GoMixedReceiverTypes
func (t TimeBasicDateTime) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format(FormatTimeBasicDateTime)), nil
}

//goland:noinspection GoMixedReceiverTypes
func (t *TimeBasicDateTime) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse(FormatTimeBasicDateTime, string(text))
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	*t = TimeBasicDateTime(parsedTime)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (t TimeBasicDateTime) GetOpenSearchFieldType() string {
	return "basic_date_time"
}

// TimeBasicDateTimeNoMillis

//goland:noinspection GoMixedReceiverTypes
func (t TimeBasicDateTimeNoMillis) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format(FormatTimeBasicDateTimeNoMillis)), nil
}

//goland:noinspection GoMixedReceiverTypes
func (t *TimeBasicDateTimeNoMillis) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse(FormatTimeBasicDateTimeNoMillis, string(text))
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	*t = TimeBasicDateTimeNoMillis(parsedTime)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (t TimeBasicDateTimeNoMillis) GetOpenSearchFieldType() string {
	return "basic_date_time_no_millis"
}

// TimeBasicDate

//goland:noinspection GoMixedReceiverTypes
func (t TimeBasicDate) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format(FormatTimeBasicDate)), nil
}

//goland:noinspection GoMixedReceiverTypes
func (t *TimeBasicDate) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse(FormatTimeBasicDate, string(text))
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	*t = TimeBasicDate(parsedTime)
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (t TimeBasicDate) GetOpenSearchFieldType() string {
	return "basic_date"
}

// NumericTime

// MarshalJSON converts NumericTime to a Unix timestamp (seconds) for JSON encoding.
func (nt NumericTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(nt.Time.Unix())
}

// UnmarshalJSON parses a Unix timestamp (seconds) from JSON and converts it to time.Time.
func (nt *NumericTime) UnmarshalJSON(data []byte) error {
	var unixTimestamp int64
	if err := json.Unmarshal(data, &unixTimestamp); err != nil {
		return err
	}
	nt.Time = time.Unix(unixTimestamp, 0).UTC() // Ensure UTC
	return nil
}

// Unix converts NumericTime to its Unix timestamp representation (seconds).
func (nt NumericTime) Unix() int64 {
	return nt.Time.Unix()
}
