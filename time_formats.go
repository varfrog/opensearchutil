package opensearchutil

import "time"

type (
	// TimeBasicDateTime marshalls into OpenSearch basic_date_time type
	TimeBasicDateTime time.Time

	// TimeBasicDateTimeNoMillis marshalls into OpenSearch basic_date_time_no_millis type
	TimeBasicDateTimeNoMillis time.Time

	// TimeBasicDate marshalls into OpenSearch basic_date type
	TimeBasicDate time.Time
)

// OpenSearchDateType tells MappingPropertiesBuilder that a type is a "date" OpenSearch type.
type OpenSearchDateType interface {
	GetOpenSearchFieldType() string
}

func (t TimeBasicDateTime) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format("20060102T150405.999-07:00")), nil
}

func (t TimeBasicDateTime) GetOpenSearchFieldType() string {
	return "basic_date_time"
}

func (t TimeBasicDateTimeNoMillis) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format("20060102T150405-07:00")), nil
}

func (t TimeBasicDateTimeNoMillis) GetOpenSearchFieldType() string {
	return "basic_date_time_no_millis"
}

func (t TimeBasicDate) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format("20060102")), nil
}

func (t TimeBasicDate) GetOpenSearchFieldType() string {
	return "basic_date"
}
