package strhelpers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/ramseyjiang/go-micros/shared/apierror"
	"github.com/ramseyjiang/go-micros/shared/srvlogs/v2"
)

var (
	// DefaultLocalTimeZoneStr is the local timezone
	DefaultLocalTimeZoneStr = "Pacific/Auckland"
	// GlobalLocalTimeZone stores the preloaded local timezone object
	GlobalLocalTimeZone *time.Location
	DateDebug           = false
)

func init() {
	prev := DateDebug
	DateDebug = true
	TimeZoneInit()
	DateDebug = prev
}

func TimeZoneInit() error {
	var tzerr error
	GlobalLocalTimeZone, tzerr = time.LoadLocation(DefaultLocalTimeZoneStr)
	if tzerr != nil {
		if DateDebug {
			srvlogs.Errorf("ERROR: Unable to load timezone (%s): %v", GlobalLocalTimeZone, tzerr)
		}
		return tzerr
	}
	return nil
}

type timeDecoder struct {
	Format string
	Name   string
}

var TimeFormats = []timeDecoder{
	timeDecoder{Name: "RFC3339Nano", Format: time.RFC3339Nano},
	timeDecoder{Name: "ISO8601Format", Format: "2006-01-02T15:04:05-0700"},
	timeDecoder{Name: "xsd.DateTime", Format: "2006-01-02T15:04:05-07:00"},
	timeDecoder{Name: "elasticSearch", Format: "2006-01-02 15:04:05.999 -0700 UTC"},
	timeDecoder{Name: "GoDefault", Format: "2006-01-02 15:04:05 -0700 MST"},
	timeDecoder{Name: "chorusDate", Format: "2006-01-02T15:04:05Z"},
	timeDecoder{Name: "DashedDateSpaceTime", Format: "2006-01-02 15:04:05.999"},
	timeDecoder{Name: "dashedDateTime", Format: "2006-01-02T15:04:05.999"},
	timeDecoder{Name: "DashedDateTTimeNoNanos", Format: "2006-01-02T15:04:05"},
	timeDecoder{Name: "SlashedDateSpaceTime", Format: "2006/01/02 15:04:05.999"},
	timeDecoder{Name: "SlashedDateSpaceTimeNoNanos", Format: "2006/01/02 15:04:05"},
	timeDecoder{Name: "DashedDate", Format: "2006-01-02"},
	timeDecoder{Name: "SlashedDate", Format: "2006/01/02"},
}

// DecodeDateTimeString attempts to decode a date/time string in various formats. The date must be past the year 1970
func DecodeDateTimeString(timeString string, OverrideTimezone *time.Location) (time.Time, string, error) {
	return decodeDateTimeString(timeString, OverrideTimezone, false)
}

// DecodeDateTimeStringAllow1970 attempts to decode a date/time string in various formats. Any date is supported
func DecodeDateTimeStringAllow1970(timeString string, OverrideTimezone *time.Location) (time.Time, string, error) {
	return decodeDateTimeString(timeString, OverrideTimezone, true)
}

// DecodeDateTimeString attempts to decode a date/time string in various formats. The date must be past the year 1970
func decodeDateTimeString(timeString string, OverrideTimezone *time.Location, allow1970 bool) (time.Time, string, error) {

	timeString = strings.TrimSpace(timeString)
	if timeString == "" {
		if DateDebug {
			srvlogs.Debugf("DecodeDateTimeString: Empty date/time string")
		}
		return time.Time{}, "", fmt.Errorf("Empty date/time string")
	}

	var LocalTimeZone *time.Location
	if OverrideTimezone != nil {
		LocalTimeZone = OverrideTimezone
	} else {
		LocalTimeZone = GlobalLocalTimeZone
	}

	if LocalTimeZone == nil {
		TimeZoneInit()
		LocalTimeZone = GlobalLocalTimeZone
		if LocalTimeZone == nil {
			return time.Time{}, "", apierror.NewAPIError(nil, http.StatusInternalServerError, "datetime", "Date time parser not initialised")
		}
	}

	for i := range TimeFormats {
		t, err := time.Parse(TimeFormats[i].Format, timeString)
		if err == nil && (allow1970 || t.Year() > 1971) {
			if DateDebug {
				srvlogs.Debugf("DecodeDateTimeString(%s): %s==%s\n", TimeFormats[i].Name, timeString, t.Format(time.RFC3339Nano))
			}
			return t, TimeFormats[i].Format, nil
		}
	}

	return time.Time{}, "", fmt.Errorf("DecodeDateTimeString ERROR: Tried %d methods to decode (%s)", len(TimeFormats), timeString)
}

// AddXMonthsMinus1Day adds X number of months, minus 1 day.
// eg: 2017-01-31 + 1 month  = 2017-02-28
// eg: 2017-01-31 + 2 months = 2017-03-31
// eg: 2017-01-31 + 3 months = 2017-04-30
func AddXMonthsMinus1Day(oldDate time.Time, months int) time.Time {
	newDate := oldDate.AddDate(0, months, 0)
	if oldDate.Day() >= 28 && newDate.Day() < oldDate.Day() {
		newDate = newDate.AddDate(0, 0, -newDate.Day())
	}
	return newDate
}

// NumPartMonthsDiff returns the number of months (or part thereof) between
// the supplied a and b time.Time variables
func NumPartMonthsDiff(a time.Time, b time.Time, max int) int {

	if a.After(b) {
		return 0
	}

	for monthCount := 0; monthCount < max; monthCount++ {
		testDate := AddXMonthsMinus1Day(a, monthCount)
		if testDate.After(b) {
			return monthCount
		}

	} // for
	return max
}

func BeginOfDay(t time.Time, timezone *time.Location) time.Time {
	if timezone == nil {
		year, month, day := t.Date()
		return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	}
	year, month, day := t.In(timezone).Date()
	return time.Date(year, month, day, 0, 0, 0, 0, timezone)
}

func EndOfDay(t time.Time, timezone *time.Location) time.Time {
	if timezone == nil {
		year, month, day := t.Date()
		return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
	}
	year, month, day := t.In(timezone).Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, timezone)
}

func GetTimeCode(slicer time.Duration) int64 {
	t := time.Now().UnixNano()
	s := slicer.Nanoseconds()

	d := t / s
	return d
}

func GetTimeWindow(slicer time.Duration) int64 {
	bod := BeginOfDay(time.Now(), nil)

	durationSinceBOD := time.Now().Sub(bod)

	s := slicer
	d := durationSinceBOD / s
	return d.Nanoseconds()
}

func DateFromTimeStamp(timestamp int64, loc *time.Location) time.Time {
	location := loc
	if location == nil {
		location = GlobalLocalTimeZone
	}

	return time.Unix(0, timestamp*int64(math.Pow10(18-int(math.Log10(float64(timestamp)))))).In(location)
}
