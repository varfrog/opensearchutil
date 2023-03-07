package opensearchutil

import (
	"github.com/pkg/errors"
	"time"
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
)

// OpenSearchDateType tells MappingPropertiesBuilder that a type is a "date" OpenSearch type.
type OpenSearchDateType interface {
	GetOpenSearchFieldType() string
}

func (t *TimeBasicDateTime) MarshalText() ([]byte, error) {
	return []byte(time.Time(*t).Format(FormatTimeBasicDateTime)), nil
}

func (t *TimeBasicDateTime) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse(FormatTimeBasicDateTime, string(text))
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	*t = TimeBasicDateTime(parsedTime)
	return nil
}

func (t TimeBasicDateTime) GetOpenSearchFieldType() string {
	return "basic_date_time"
}

// TimeBasicDateTimeNoMillis

func (t *TimeBasicDateTimeNoMillis) MarshalText() ([]byte, error) {
	return []byte(time.Time(*t).Format(FormatTimeBasicDateTimeNoMillis)), nil
}

func (t *TimeBasicDateTimeNoMillis) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse(FormatTimeBasicDateTimeNoMillis, string(text))
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	*t = TimeBasicDateTimeNoMillis(parsedTime)
	return nil
}

func (t TimeBasicDateTimeNoMillis) GetOpenSearchFieldType() string {
	return "basic_date_time_no_millis"
}

// TimeBasicDate

func (t *TimeBasicDate) MarshalText() ([]byte, error) {
	return []byte(time.Time(*t).Format(FormatTimeBasicDate)), nil
}

func (t *TimeBasicDate) UnmarshalText(text []byte) error {
	parsedTime, err := time.Parse(FormatTimeBasicDate, string(text))
	if err != nil {
		return errors.Wrap(err, "time.Parse")
	}
	*t = TimeBasicDate(parsedTime)
	return nil
}

func (t TimeBasicDate) GetOpenSearchFieldType() string {
	return "basic_date"
}
