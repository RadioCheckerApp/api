package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"reflect"
	"testing"
	"time"
)

func TestCreateTracksWorker(t *testing.T) {
	date := time.Now()
	dateStr := date.Format("2006-01-02")
	loc, _ := time.LoadLocation("Europe/Berlin")

	var tests = []struct {
		dao               datalayer.TrackRecordDAO
		pathParams        map[string]string
		queryStringParams map[string]string
		expectedResult    Worker
		expectedErr       bool
	}{
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"date": dateStr, "filter": "top"},
			DayTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "station-a"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				Top,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": ""},
			map[string]string{"date": dateStr, "filter": "top"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"week": dateStr, "filter": "top"},
			WeekTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "station-a"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				Top,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"date": dateStr, "filter": "all"},
			DayTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "station-a"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				All,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"week": dateStr, "filter": "all"},
			WeekTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "station-a"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				All,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "CamelCaseStation"},
			map[string]string{"week": dateStr, "filter": "all"},
			WeekTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "camelcasestation"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				All,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "invalidDateDay"},
			map[string]string{"date": "2018-07-32", "filter": "all"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "invalidDateWeek"},
			map[string]string{"week": "2018-07-32", "filter": "all"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "invalidFilter"},
			map[string]string{"date": dateStr, "filter": "invalidFilter"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "noTracksStation"},
			map[string]string{"week": dateStr, "filter": ""},
			WeekTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "notracksstation"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				Top,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "missingFilter"},
			map[string]string{"date": dateStr},
			DayTracksWorker{
				TracksWorker{MockTrackRecordDAO{}, "missingfilter"},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
				Top,
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "missingDateOrWeek"},
			map[string]string{"filter": "top"},
			nil,
			true,
		},
	}

	for _, test := range tests {
		result, err := CreateTracksWorker(test.dao, test.pathParams, test.queryStringParams)
		if (err != nil) != test.expectedErr {
			t.Errorf("CreateTracksWorker(%q, %q, %q): got (%q, %v), expected error: %v",
				test.dao, test.pathParams, test.queryStringParams, result, err,
				test.expectedErr)
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("CreateTracksWorker(%q, %q, %q): got return type (%v), expected (%v)",
				test.dao, test.pathParams, test.queryStringParams,
				reflect.TypeOf(result), reflect.TypeOf(test.expectedResult))
			continue
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Errorf("CreateTracksWorker(%q, %q, %q): got \n(%q), expected \n(%q)",
				test.dao, test.pathParams, test.queryStringParams, result, test.expectedResult)
		}
	}
}

func TestCreateSearchWorker(t *testing.T) {
	date := time.Now()
	dateStr := date.Format("2006-01-02")
	loc, _ := time.LoadLocation("Europe/Berlin")

	var tests = []struct {
		dao               datalayer.TrackRecordDAO
		queryStringParams map[string]string
		expectedResult    Worker
		expectedErr       bool
	}{
		{
			MockTrackRecordDAO{},
			map[string]string{"date": dateStr, "q": "dani+california"},
			DaySearchWorker{
				SearchWorker{MockTrackRecordDAO{}, []string{"dani", "california"}},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"week": dateStr, "q": "dani+california"},
			WeekSearchWorker{
				SearchWorker{MockTrackRecordDAO{}, []string{"dani", "california"}},
				time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc),
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"date": "2018-07-32", "q": "dani+california"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"week": "2018-07-32", "q": "dani+california"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"date": dateStr},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"week": dateStr, "q": ""},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"q": "dani+california"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{},
			nil,
			true,
		},
	}

	for _, test := range tests {
		result, err := CreateSearchWorker(test.dao, test.queryStringParams)
		if (err != nil) != test.expectedErr {
			t.Errorf("CreateSearchWorker(%q, %q): got (%q, %v), expected error: %v",
				test.dao, test.queryStringParams, result, err,
				test.expectedErr)
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("CreateSearchWorker(%q, %q): got return type (%v), expected (%v)",
				test.dao, test.queryStringParams, reflect.TypeOf(result),
				reflect.TypeOf(test.expectedResult))
			continue
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Errorf("CreateSearchWorker(%q, %q): got \n(%q), expected \n(%q)",
				test.dao, test.queryStringParams, result, test.expectedResult)
		}
	}
}
