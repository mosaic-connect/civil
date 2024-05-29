package civil

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

// DateTime represents a date-time without a timezone.
// Calculations on DateTime are performed using the standard
// library's time.Time type. For these calculations the
// timezone is UTC.
//
// DateTime is useful in situations where a date and time
// are specified, without reference to a timezone. Although not
// common, it can be useful. For example, a dose of medication
// may be scheduled for a particular date and time, regardless
// of the timezone that the patient is residing in at the time.
//
// Because DateTime does not specify a unique instant in
// time, it has never been necessary to specify to sub-second
// accuracy. For this reason DateTime only specifies the
// time to second accuracy. In actual fact, DateTime would
// probably be fine if it only specified to minute accuracy.
type DateTime struct {
	t time.Time
}

// After reports whether the civil date-time d is after e
func (dt DateTime) After(e DateTime) bool {
	return dt.t.After(e.t)
}

// Before reports whether the civil date-time d is before e
func (dt DateTime) Before(e DateTime) bool {
	return dt.t.Before(e.t)
}

// Equal reports whether dt and e represent the same civil date-time.
func (dt DateTime) Equal(e DateTime) bool {
	return dt.t.Equal(e.t)
}

// IsZero reports whether dt represents the zero civil date-time,
// Midnight, January 1, year 1.
func (dt DateTime) IsZero() bool {
	return dt.t.IsZero()
}

// Date returns the year, month and day on which dt occurs.
func (dt DateTime) Date() (year int, month time.Month, day int) {
	return dt.t.Date()
}

// Clock returns the hour, minute and second on which dt occurs.
func (dt DateTime) Clock() (hour int, minute int, second int) {
	hour = dt.Hour()
	minute = dt.Minute()
	second = dt.Second()
	return
}

// DateTime returns the year, month, day, hour minute, second and nanosecond on which dt occurs.
func (dt DateTime) DateTime() (year int, month time.Month, day int, hour int, minute int, second int) {
	year, month, day = dt.t.Date()
	hour, minute, second = dt.Clock()
	return
}

// Unix returns d as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC to midnight of the date-time UTC.
func (dt DateTime) Unix() int64 {
	return dt.t.Unix()
}

// Year returns the year in which dt occurs.
func (dt DateTime) Year() int {
	return dt.t.Year()
}

// Month returns the month of the year specified by dt.
func (dt DateTime) Month() time.Month {
	return dt.t.Month()
}

// Day returns the day of the month specified by dt.
func (dt DateTime) Day() int {
	return dt.t.Day()
}

// Hour returns the hour specified by dt.
func (dt DateTime) Hour() int {
	return dt.t.Hour()
}

// Minute returns the minute specified by dt.
func (dt DateTime) Minute() int {
	return dt.t.Minute()
}

// Second returns the second specified by dt.
func (dt DateTime) Second() int {
	return dt.t.Second()
}

// Weekday returns the day of the week specified by d.
func (dt DateTime) Weekday() time.Weekday {
	return dt.t.Weekday()
}

// ISOWeek returns the ISO 8601 year and week number in which d occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
func (dt DateTime) ISOWeek() (year, week int) {
	year, week = dt.t.ISOWeek()
	return
}

// YearDay returns the day of the year specified by D, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
func (dt DateTime) YearDay() int {
	return dt.t.YearDay()
}

// Add returns the civil date-time d + duration.
func (dt DateTime) Add(duration time.Duration) DateTime {
	t := dt.t.Add(toSeconds(duration))
	return DateTime{t: t}
}

// Sub returns the duration dt-e, which will be an integral number of seconds.
// If the result exceeds the maximum (or minimum) value that can be stored
// in a Duration, the maximum (or minimum) duration will be returned.
// To compute dt-duration, use dt.Add(-duration).
func (dt DateTime) Sub(e DateTime) time.Duration {
	return dt.t.Sub(e.t)
}

// AddDate returns the civil date-time corresponding to adding the given number of years,
// months, and days to t. For example, AddDate(-1, 2, 3) applied to January 1, 2011
// returns March 4, 2010.
//
// AddDate normalizes its result in the same way that Date does, so, for example,
// adding one month to October 31 yields December 1, the normalized form for November 31.
func (dt DateTime) AddDate(years int, months int, days int) DateTime {
	t := dt.t.AddDate(years, months, days)
	return DateTime{t: t}
}

// toDate converts the time.Time value into a DateTime.,
func toLocalDateTime(t time.Time) DateTime {
	y, m, d := t.Date()
	hour, minute, second := t.Clock()
	return DateTimeFor(y, m, d, hour, minute, second)
}

// Now returns the current civil date-time.
func Now() DateTime {
	return toLocalDateTime(time.Now())
}

// DateTimeFor returns the DateTime corresponding to year, month, day, hour, minute and second.
//
// The month and day values may be outside their usual ranges
// and will be normalized during the conversion.
// For example, October 32 converts to November 1.
func DateTimeFor(year int, month time.Month, day int, hour int, minute int, second int) DateTime {
	return DateTime{
		t: time.Date(year, month, day, hour, minute, second, 0, time.UTC),
	}
}

// DateTimeOf returns the DateTime corresponding to t in t's location.
func DateTimeOf(t time.Time) DateTime {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return DateTimeFor(year, month, day, hour, minute, second)
}

// Format returns a textual representation of the time value formatted
// according to layout, which takes the same form as the standard library
// time package. Note that with a Date the reference time is
//  Mon Jan 2 2006 15:04:05.
func (dt DateTime) Format(layout string) string {
	return dt.t.Format(layout)
}

// String returns a string representation of d. The date
// format returned is compatible with ISO 8601: yyyy-mm-ddTHH:MM:SS.
func (dt DateTime) String() string {
	return localDateTimeString(dt)
}

// localDateTimeString returns the string representation of the date.
func localDateTimeString(dt DateTime) string {
	year, month, day, hour, minute, second := dt.DateTime()
	sign := ""
	if year < 0 {
		year = -year
		sign = "-"
	}
	return fmt.Sprintf("%s%04d-%02d-%02dT%02d:%02d:%02d", sign, year, int(month), day, hour, minute, second)
}

// localDateQuotedString returns the string representation of the date in quotation marks.
func localDateQuotedString(dt DateTime) string {
	return fmt.Sprintf(`"%s"`, localDateTimeString(dt))
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (dt DateTime) MarshalBinary() ([]byte, error) {
	return dt.t.MarshalBinary()
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (dt *DateTime) UnmarshalBinary(data []byte) error {
	var t time.Time
	if err := t.UnmarshalBinary(data); err != nil {
		return err
	}
	*dt = DateTimeOf(t)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// The date is a quoted string in an ISO 8601 format (yyyy-mm-ddTHH:MM:SS).
func (dt DateTime) MarshalJSON() ([]byte, error) {
	return []byte(localDateQuotedString(dt)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The date is expected to be a quoted string in an ISO 8601
// format (calendar or ordinal).
func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	*dt, err = ParseDateTime(s)
	return
}

// MarshalText implements the encoding.TextMarshaller interface.
// The date format is yyyy-mm-dd.
func (dt DateTime) MarshalText() ([]byte, error) {
	return []byte(localDateTimeString(dt)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaller interface.
// The date is expected to an ISO 8601 format (calendar or ordinal).
func (dt *DateTime) UnmarshalText(data []byte) (err error) {
	s := string(data)
	*dt, err = ParseDateTime(s)
	return
}

// Scan implements the sql.Scanner interface.
func (dt *DateTime) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		{
			d1, err := ParseDateTime(v)
			if err != nil {
				return err
			}
			*dt = d1
		}
	case []byte:
		{
			d1, err := ParseDateTime(string(v))
			if err != nil {
				return err
			}
			*dt = d1
		}
	case time.Time:
		{
			d1 := DateTimeOf(v)
			*dt = d1
		}
	case nil:
		*dt = DateTime{}
	default:
		return errors.New("cannot convert to civil.DateTime")
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (dt DateTime) Value() (driver.Value, error) {
	year, month, day := dt.Date()
	hour, minute, second := dt.Clock()
	return time.Date(year, month, day, hour, minute, second, 0, time.UTC), nil
}
