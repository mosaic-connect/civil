package civil

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	assert := assert.New(t)
	datetime := Now()
	now := time.Now()

	y1, m1, d1 := now.Date()
	y2, m2, d2 := datetime.Date()

	assert.Equal(d1, d2)
	assert.Equal(m1, m2)
	assert.Equal(y1, y2)

	y3, m3, d3, h3, mi3, s3 := datetime.DateTime()
	assert.Equal(d1, d3)
	assert.Equal(m1, m3)
	assert.Equal(y1, y3)
	assert.Equal(now.Hour(), h3)
	assert.Equal(now.Minute(), mi3)
	assert.Equal(now.Second(), s3)
}

func BenchmarkDateTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		year := 1934
		month := time.Month(3)
		day := 18
		hour := 22
		minute := 2
		second := 31
		_ = DateTimeFor(year, month, day, hour, minute, second)
	}
}

func TestDateTimeYears(t *testing.T) {
	for year := -9999; year <= 9999; year++ {
		month := 5
		day := 20
		hour := 22
		minute := 2
		second := 31

		datetime := DateTimeFor(year, time.Month(month), day, hour, minute, second)
		CheckLocalDateTime(t, datetime, year, month, day, hour, minute, second)
	}
}

func TestDateTimeMonths(t *testing.T) {
	for month := 1; month <= 12; month++ {
		year := 1969
		day := 12
		hour := 22
		minute := 2
		second := 31

		datetime := DateTimeFor(year, time.Month(month), day, hour, minute, second)
		CheckLocalDateTime(t, datetime, year, month, day, hour, minute, second)
	}
}

func TestDateTimeDays(t *testing.T) {
	for day := 1; day <= 31; day++ {
		year := 1969
		month := 1
		hour := 22
		minute := 2
		second := 31

		datetime := DateTimeFor(year, time.Month(month), day, hour, minute, second)
		CheckLocalDateTime(t, datetime, year, month, day, hour, minute, second)
	}
}

func CheckLocalDateTime(t *testing.T, datetime DateTime, year, month, day, hour, minute, second int) {
	assert := assert.New(t)
	assert.Equal(year, datetime.Year())
	assert.Equal(month, int(datetime.Month()))
	assert.Equal(day, datetime.Day())
	assert.Equal(hour, datetime.Hour())
	assert.Equal(minute, datetime.Minute())
	assert.Equal(second, datetime.Second())

	// Calculate expected text representation
	var text = datetime.t.Format("2006-01-02T15:04:05")

	assert.Equal(text, datetime.String())

	datetime2, err := ParseDateTime(text)
	assert.NoError(err)
	assert.True(datetime.Equal(datetime2))

	// for non-negative years, can check parsing with time package
	if year >= 0 {
		if tm, err := time.Parse("2006-01-02T15:04:05", text); err != nil {
			t.Errorf("time.Parse: unexpected error parsing %s: %v", text, err)
		} else {
			y := tm.Year()
			m := int(tm.Month())
			d := tm.Day()
			if y != year {
				t.Errorf("time.Parse: Year: expected %d, actual %d", year, y)
			}
			if m != month {
				t.Errorf("time.Parse: Month: expected %d, actual %d", month, m)
			}
			if d != day {
				t.Errorf("time.Parse: Day: expected %d, actual %d", day, d)
			}
		}
	}

	// check marshalling and unmarshalling JSON
	data, err := datetime.MarshalJSON()
	if err != nil {
		t.Errorf("MarshalJSON: %s: unexpected error: %v", text, err)
	} else {
		assert.Equal(`"`+text+`"`, string(data))
		var datetime2 DateTime
		err = datetime2.UnmarshalJSON(data)
		if err != nil {
			t.Errorf("UnmarshalJSON: %s: unexpected error: %v", text, err)
		} else {
			if !datetime.Equal(datetime2) {
				t.Errorf("UnmarshalJSON: expected %s, actual %s", datetime.String(), datetime2.String())
			}
		}
	}

	// check marshalling and unmarshalling text
	data, err = datetime.MarshalText()
	if err != nil {
		t.Errorf("MarshalText: %s: unexpected error: %v", text, err)
	} else {
		assert.Equal(text, string(data))
		var datetime2 DateTime
		err = datetime2.UnmarshalText(data)
		if err != nil {
			t.Errorf("UnmarshalText: %s: unexpected error: %v", text, err)
		} else {
			if !datetime.Equal(datetime2) {
				t.Errorf("UnmarshalText: expected %s, actual %s", datetime.String(), datetime2.String())
			}
		}
	}

	// marshal and unmarshal binary
	data, err = datetime.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary: %s: unexpected error: %v", text, err)
	} else {
		// binary should be the same as the equivalent time binary
		tdata, _ := datetime.t.MarshalBinary()
		assert.Equal(tdata, data)
		var datetime2 DateTime
		err = datetime2.UnmarshalBinary(data)
		assert.NoError(err, datetime.String())
	}
}

func TestParseDateDateTime(t *testing.T) {
	testCases := []struct {
		Text   string
		Valid  bool
		Day    int
		Month  time.Month
		Year   int
		Hour   int
		Minute int
		Second int
	}{
		{
			Text:  "2095-09-30",
			Valid: true,
			Day:   30,
			Month: time.September,
			Year:  2095,
		},
		{
			Text:  "2195-060",
			Valid: true,
			Day:   1,
			Month: time.March,
			Year:  2195,
		},
		{
			Text:  "2095.09.30",
			Valid: true,
			Day:   30,
			Month: time.September,
			Year:  2095,
		},
		{
			Text:  "2095/09/30",
			Valid: true,
			Day:   30,
			Month: time.September,
			Year:  2095,
		},
		{
			Text:  "20951030",
			Valid: true,
			Day:   30,
			Month: time.October,
			Year:  2095,
		},
		{
			Text:  "2195-060",
			Valid: true,
			Day:   1,
			Month: time.March,
			Year:  2195,
		},
		{
			Text:  "2195074",
			Valid: true,
			Day:   15,
			Month: time.March,
			Year:  2195,
		},
		{
			Text:   "2095-09-30T1:2:3",
			Valid:  true,
			Day:    30,
			Month:  time.September,
			Year:   2095,
			Hour:   1,
			Minute: 2,
			Second: 3,
		},
		{
			Text:   "2195-060T030211",
			Valid:  true,
			Day:    1,
			Month:  time.March,
			Year:   2195,
			Hour:   3,
			Minute: 2,
			Second: 11,
		},
		{
			Text:   "2095.09.30T12:39",
			Valid:  true,
			Day:    30,
			Month:  time.September,
			Year:   2095,
			Hour:   12,
			Minute: 39,
		},
		{
			Text:   "2095/09/30T1147",
			Valid:  true,
			Day:    30,
			Month:  time.September,
			Year:   2095,
			Hour:   11,
			Minute: 47,
		},
		{
			Text:   "20951030T10:11:12.123456789",
			Valid:  true,
			Day:    30,
			Month:  time.October,
			Year:   2095,
			Hour:   10,
			Minute: 11,
			Second: 12,
		},
		{
			Text:   "2195-060T121110.1234",
			Valid:  true,
			Day:    1,
			Month:  time.March,
			Year:   2195,
			Hour:   12,
			Minute: 11,
			Second: 10,
		},
		{
			Text:   "2195074T001122.",
			Valid:  true,
			Day:    15,
			Month:  time.March,
			Year:   2195,
			Hour:   0,
			Minute: 11,
			Second: 22,
		},
	}
	assert := assert.New(t)

	for _, tc := range testCases {
		for _, text := range []string{tc.Text, " \t\t" + tc.Text + "\t\t\t "} {
			ld, err := ParseDateTime(text)
			if tc.Valid {
				assert.NoError(err, text)
				year, month, day, hour, minute, second := ld.DateTime()
				assert.Equal(tc.Day, day, text)
				assert.Equal(tc.Month, month, text)
				assert.Equal(tc.Year, year, text)
				assert.Equal(tc.Hour, hour, text)
				assert.Equal(tc.Minute, minute, text)
				assert.Equal(tc.Second, second, text)
			} else {
				assert.Error(err, text)
			}
		}
	}
}

func TestDateTimeMarshalXML(t *testing.T) {
	assert := assert.New(t)
	type testStruct struct {
		XMLName        xml.Name `xml:"TestCase"`
		Element        DateTime
		AnotherElement DateTime `xml:"another"`
		Attribute1     DateTime `xml:",attr"`
		Attribute2     DateTime `xml:"attribute-2,attr"`
	}

	testCases := []struct {
		st  testStruct
		xml string
	}{
		{
			st: testStruct{
				Element:        mustParseDateTime("2021-01-02T01:02:03"),
				AnotherElement: mustParseDateTime("2021-01-03T04:05:06"),
				Attribute1:     mustParseDateTime("2021-01-04T07:08:09"),
				Attribute2:     mustParseDateTime("2021-01-05T10:11:12"),
			},
			xml: `<TestCase Attribute1="2021-01-04T07:08:09" attribute-2="2021-01-05T10:11:12">` +
				`<Element>2021-01-02T01:02:03</Element><another>2021-01-03T04:05:06</another></TestCase>`,
		},
	}

	for _, tc := range testCases {
		b, err := xml.Marshal(&tc.st)
		assert.NoError(err)
		assert.Equal(tc.xml, string(b))
		var st testStruct
		err = xml.Unmarshal([]byte(tc.xml), &st)
		assert.NoError(err)
		assert.Equal("", st.XMLName.Space)
		assert.Equal("TestCase", st.XMLName.Local)
		st.XMLName.Local = ""
		assert.Equal(tc.st, st)
	}
}

func TestDateTimeAfter(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		DateTime1, DateTime2 DateTime
	}{
		{DateTimeFor(1999, 9, 30, 1, 2, 3), DateTimeFor(1999, 9, 30, 1, 2, 4)},
		{DateTimeFor(0, 9, 30, 3, 2, 1), DateTimeFor(0, 9, 30, 3, 2, 2)},
	}

	for _, tc := range testCases {
		assert.True(tc.DateTime1.Before(tc.DateTime2))
		assert.True(tc.DateTime2.After(tc.DateTime1))
		assert.False(tc.DateTime2.Before(tc.DateTime1))
		assert.False(tc.DateTime1.After(tc.DateTime2))
	}
}

func TestDateTimeWeekday(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		DateTime DateTime
		Weekday  time.Weekday
	}{
		{DateTimeFor(1999, 9, 30, 1, 2, 3), time.Thursday},
		{DateTimeFor(1997, 1, 30, 4, 5, 6), time.Thursday},
		{DateTimeFor(1994, 11, 14, 7, 8, 9), time.Monday},
		{DateTimeFor(1992, 12, 16, 10, 11, 12), time.Wednesday},
		{DateTimeFor(2033, 1, 4, 13, 14, 15), time.Tuesday},
		{DateTimeFor(2033, 4, 8, 16, 17, 18), time.Friday},
		{DateTimeFor(2033, 4, 9, 19, 20, 21), time.Saturday},
		{DateTimeFor(2042, 7, 6, 22, 23, 24), time.Sunday},
	}

	for _, tc := range testCases {
		assert.Equal(tc.Weekday, tc.DateTime.Weekday(), tc.DateTime.String())
	}
}

func TestDateTimeScan(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Value    interface{}
		Error    bool
		Expected DateTime
	}{
		{
			Value:    "2056-11-13T12:34:56",
			Expected: DateTimeFor(2056, 11, 13, 12, 34, 56),
		},
		{
			Value:    "2056-11-13 12:34:56.000",
			Expected: DateTimeFor(2056, 11, 13, 12, 34, 56),
		},
		{
			Value:    time.Date(2056, 10, 31, 16, 34, 12, 0, time.UTC),
			Expected: DateTimeFor(2056, 10, 31, 16, 34, 12),
		},
		{
			Value:    time.Date(2056, 9, 30, 1, 2, 3, 400000, time.FixedZone("Australia/Brisbane", 10*3600)),
			Expected: DateTimeFor(2056, 9, 30, 1, 2, 3),
		},
		{Value: []byte("2157-12-31"), Expected: DateTimeFor(2157, 12, 31, 0, 0, 0)},
		{Value: []byte("zzz"), Error: true},
		{Value: nil, Expected: DateTime{}},
		{Value: int64(11), Error: true},
		{Value: true, Error: true},
		{Value: float64(11.1), Error: true},
	}

	for _, tc := range testCases {
		var d DateTime
		err := d.Scan(tc.Value)
		if tc.Error {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.True(d.Equal(tc.Expected))
		}
	}
}

func TestDateTimeValue(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Date     Date
		Expected time.Time
	}{
		{DateFor(2071, 1, 30), time.Date(2071, 1, 30, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		v, err := tc.Date.Value()
		assert.NoError(err)
		assert.IsType(time.Time{}, v)
		t := v.(time.Time)
		assert.True(t.Equal(tc.Expected))
	}
}

func mustParseDateTime(s string) DateTime {
	dt, err := ParseDateTime(s)
	if err != nil {
		panic(err.Error())
	}
	return dt
}

func TestDateTimeProperties(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		DateTime  DateTime
		IsZero    bool
		Unix      int64
		Year      int
		Week      int
		YearDay   int
		Formatted string
	}{
		{
			DateTime:  DateTime{},
			IsZero:    true,
			Unix:      -62135596800,
			Year:      1,
			Week:      1,
			YearDay:   1,
			Formatted: "1 Jan 0001 00:00:00",
		},
		{
			DateTime:  DateTimeFor(1970, 1, 1, 0, 0, 0),
			IsZero:    false,
			Unix:      0,
			Year:      1970,
			Week:      1,
			YearDay:   1,
			Formatted: "1 Jan 1970 00:00:00",
		},
		{
			DateTime:  DateTimeFor(2048, 1, 30, 12, 34, 56),
			IsZero:    false,
			Unix:      2464000496,
			Year:      2048,
			Week:      5,
			YearDay:   30,
			Formatted: "30 Jan 2048 12:34:56",
		},
	}

	for _, tc := range testCases {
		assert.Equal(tc.IsZero, tc.DateTime.IsZero())
		assert.Equal(tc.Unix, tc.DateTime.Unix())
		year, week := tc.DateTime.ISOWeek()
		assert.Equal(tc.Year, year, "Year")
		assert.Equal(tc.Week, week, "Week")
		assert.Equal(tc.YearDay, tc.DateTime.YearDay())
		assert.Equal(tc.Formatted, tc.DateTime.Format("2 Jan 2006 15:04:05"))
	}
}

// Test for case where attempt made to unmarshal invalid binary data
func TestDateTimeUnmarshalBinaryError(t *testing.T) {
	assert := assert.New(t)
	data := []byte("xxxx")
	var d DateTime
	err := d.UnmarshalBinary(data)
	assert.Error(err)
}

func TestDateTimeAddDate(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		DateTime DateTime
		Years    int
		Months   int
		Days     int
		Expected DateTime
	}{
		{
			DateTime: DateTimeFor(2029, 12, 16, 11, 47, 42),
			Years:    1,
			Expected: DateTimeFor(2030, 12, 16, 11, 47, 42),
		},
		{
			DateTime: DateTimeFor(2029, 12, 16, 12, 34, 56),
			Years:    1,
			Months:   3,
			Expected: DateTimeFor(2031, 3, 16, 12, 34, 56),
		},
		{
			DateTime: DateTimeFor(2029, 12, 16, 1, 2, 3),
			Years:    1,
			Months:   3,
			Days:     30,
			Expected: DateTimeFor(2031, 4, 15, 1, 2, 3),
		},
		{
			DateTime: DateTimeFor(2029, 12, 16, 3, 2, 1),
			Years:    -1,
			Months:   3,
			Days:     30,
			Expected: DateTimeFor(2029, 4, 15, 3, 2, 1),
		},
		{
			DateTime: DateTimeFor(2029, 12, 16, 4, 5, 6),
			Months:   -13,
			Expected: DateTimeFor(2028, 11, 16, 4, 5, 6),
		},
		{
			DateTime: DateTimeFor(2029, 12, 16, 11, 21, 13),
			Days:     15,
			Expected: DateTimeFor(2029, 12, 31, 11, 21, 13),
		},
	}
	for _, tc := range testCases {
		d := tc.DateTime.AddDate(tc.Years, tc.Months, tc.Days)
		assert.Equal(tc.Expected, d, tc.Expected.String()+" vs "+d.String())
	}
}

func TestDateTimeSub(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		DateTime1 DateTime
		DateTime2 DateTime
		Days      int
		Hours     int
		Minutes   int
		Seconds   int
	}{
		{
			DateTime1: DateTimeFor(1994, 11, 14, 2, 3, 4),
			DateTime2: DateTimeFor(1994, 11, 13, 0, 0, 0),
			Days:      1,
			Hours:     2,
			Minutes:   3,
			Seconds:   4,
		},
		{
			DateTime1: DateTimeFor(1994, 11, 14, 0, 0, 0),
			DateTime2: DateTimeFor(1994, 11, 15, 2, 3, 4),
			Days:      -1,
			Hours:     -2,
			Minutes:   -3,
			Seconds:   -4,
		},
		{
			DateTime1: DateTimeFor(1994, 11, 14, 18, 58, 0),
			DateTime2: DateTimeFor(1992, 12, 16, 11, 47, 0),
			Days:      698,
			Hours:     7,
			Minutes:   11,
			Seconds:   0,
		},
	}
	for _, tc := range testCases {
		d := tc.DateTime1.Sub(tc.DateTime2)
		expected := time.Duration(tc.Days)*time.Hour*24 + time.Duration(tc.Hours)*time.Hour + time.Duration(tc.Minutes)*time.Minute + time.Duration(tc.Seconds)*time.Second
		assert.Equal(expected, d, fmt.Sprintf("%s - %s expected %d:%d:%d:%d", tc.DateTime1.String(), tc.DateTime2.String(), tc.Days, tc.Hours, tc.Minutes, tc.Seconds))
	}
}

func TestDateTimeParseLayout(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		Text     string
		Layout   string
		Expected DateTime
		Error    bool
	}{
		{
			Text:     "11 Jan 1994 05:45:23",
			Layout:   "02 Jan 2006 15:04:05",
			Expected: DateTimeFor(1994, 1, 11, 5, 45, 23),
		},
		{
			Text:   "Jan 11 1994",
			Layout: "02 Jan 2006",
			Error:  true,
		},
	}

	for _, tc := range testCases {
		d, err := ParseDateTimeLayout(tc.Layout, tc.Text)
		if tc.Error {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.Equal(tc.Expected, d, dateTimesNotEqual(tc.Expected, d))
		}
	}
}

func dateTimesNotEqual(expected, actual DateTime) string {
	return fmt.Sprintf("%s vs %s", expected.String(), actual.String())
}
