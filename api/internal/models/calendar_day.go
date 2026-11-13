package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// CalendarDay is a calendar day stored as midnight UTC and written on the wire as a bare
// yyyy-MM-dd, like the claim dates. It reads a bare day or an RFC 3339 instant (keeping
// the day as written), so a client that sends either lands on the day it picked.
type CalendarDay struct {
	time.Time
}

// NewCalendarDay returns the calendar day of t as written (its own offset is not converted).
func NewCalendarDay(t time.Time) CalendarDay {
	year, month, day := t.Date()
	return CalendarDay{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// ParseCalendarDay reads a bare yyyy-MM-dd or an RFC 3339 string.
func ParseCalendarDay(value string) (CalendarDay, bool) {
	trimmed := strings.TrimSpace(value)
	parsed, err := time.Parse(time.DateOnly, trimmed)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, trimmed)
		if err != nil {
			return CalendarDay{}, false
		}
	}
	return NewCalendarDay(parsed), true
}

func (day CalendarDay) MarshalJSON() ([]byte, error) {
	return []byte(`"` + day.UTC().Format(time.DateOnly) + `"`), nil
}

func (day *CalendarDay) UnmarshalJSON(data []byte) error {
	text := strings.Trim(string(data), `"`)
	if text == "null" || text == "" {
		*day = CalendarDay{}
		return nil
	}

	parsed, ok := ParseCalendarDay(text)
	if !ok {
		return fmt.Errorf("invalid calendar day %q", text)
	}
	*day = parsed
	return nil
}

// Value stores the day as a timestamp, the same shape every other date column has.
func (day CalendarDay) Value() (driver.Value, error) {
	return day.Time, nil
}

func (day *CalendarDay) Scan(value any) error {
	switch typed := value.(type) {
	case time.Time:
		*day = NewCalendarDay(typed)
	case string:
		return day.scanText(typed)
	case []byte:
		return day.scanText(string(typed))
	case nil:
		*day = CalendarDay{}
	default:
		return fmt.Errorf("cannot scan %T into a CalendarDay", value)
	}
	return nil
}

func (day *CalendarDay) scanText(text string) error {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05", time.DateOnly} {
		if parsed, err := time.Parse(layout, text); err == nil {
			*day = NewCalendarDay(parsed)
			return nil
		}
	}
	return fmt.Errorf("cannot scan %q into a CalendarDay", text)
}

// GormDataType keeps the column a datetime, matching the other date columns.
func (CalendarDay) GormDataType() string {
	return "datetime"
}
