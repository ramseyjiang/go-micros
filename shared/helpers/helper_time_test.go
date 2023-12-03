package helpers

import (
	"testing"
	"time"
)

var ()

func TestBeginOfDay(t *testing.T) {
	locationNZ, tzErr := time.LoadLocation("Pacific/Auckland")
	if tzErr != nil {
		t.Fatalf("Unable to load timezone: %v", tzErr)
	}

	// 2019-11-19 22:44:41.464662577 +0000 UTC m=+15.884646276,
	// timeCheck := time.Unix(0, 1574203481464662577).In(locationNZ)
	timeUTC := time.Date(2019, 11, 19, 22, 44, 41, 464662577, time.UTC)
	// timeCheck := time.Date(2019, 11, 20, 11, 44, 41, 464662577, locationNZ)
	timeCheck := timeUTC.In(locationNZ)
	timeExpected := time.Date(2019, 11, 20, 0, 0, 0, 0, locationNZ)

	timeGot := BeginOfDay(timeCheck, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%s) expected: (%s), got: (%s), diff(%s)", timeCheck.String(), timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}

	timeGot = BeginOfDay(timeUTC, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%s) expected: (%s), got: (%s), diff(%s)", timeUTC.String(), timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}
}

func TestEndOfDay(t *testing.T) {
	locationNZ, tzErr := time.LoadLocation("Pacific/Auckland")
	if tzErr != nil {
		t.Fatalf("Unable to load timezone: %v", tzErr)
	}

	// 2019-11-19 22:44:41.464662577 +0000 UTC m=+15.884646276,
	// timeCheck := time.Unix(0, 1574203481464662577).In(locationNZ)
	timeUTC := time.Date(2019, 11, 19, 22, 44, 41, 464662577, time.UTC)
	// timeCheck := time.Date(2019, 11, 20, 11, 44, 41, 464662577, locationNZ)
	timeCheck := timeUTC.In(locationNZ)
	timeExpected := time.Date(2019, 11, 20, 23, 59, 59, 999999999, locationNZ)

	timeGot := EndOfDay(timeCheck, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%s) expected: (%s), got: (%s), diff(%s)", timeCheck.String(), timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}

	timeGot = EndOfDay(timeUTC, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%s) expected: (%s), got: (%s), diff(%s)", timeUTC.String(), timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}
}

func TestDateFromTimeStamp(t *testing.T) {
	locationNZ, tzErr := time.LoadLocation("Pacific/Auckland")
	if tzErr != nil {
		t.Fatalf("Unable to load timezone: %v", tzErr)
	}

	// 1654541600 - 07/Jun/22
	timeGot := DateFromTimeStamp(1654541600, nil)
	timeExpected := time.Date(2022, 6, 7, 6, 53, 20, 0, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%d) expected: (%s), got: (%s), diff(%s)", 1654541600, timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}

	// 1654241600000 - 03/Jun/22
	timeGot = DateFromTimeStamp(1654241600000, nil)
	timeExpected = time.Date(2022, 6, 3, 19, 33, 20, 0, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%d) expected: (%s), got: (%s), diff(%s)", 1654241600000, timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}

	// 1654041600000000 - 01/Jun/22
	timeGot = DateFromTimeStamp(1654041600000000, nil)
	timeExpected = time.Date(2022, 6, 1, 12, 0, 0, 0, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%d) expected: (%s), got: (%s), diff(%s)", 1654041600000000, timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}

	// 1654941600000000000 - 11/Jun/22
	timeGot = DateFromTimeStamp(1654941600000000000, nil)
	timeExpected = time.Date(2022, 6, 11, 22, 0, 0, 0, locationNZ)
	if !timeGot.Equal(timeExpected) {
		t.Fatalf("from (%d) expected: (%s), got: (%s), diff(%s)", 1654941600000000000, timeExpected.String(), timeGot.String(), timeGot.Sub(timeExpected).String())
	}

}
