package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
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
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"filter": "latest"},
			TracksWorker{MockTrackRecordDAO{}, "station-a"},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"date": "2018-08-26", "filter": "latest"},
			TracksWorker{MockTrackRecordDAO{}, "station-a"},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": ""},
			map[string]string{},
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

func TestCreateCreateTrackWorker(t *testing.T) {
	var tests = []struct {
		trDAO          datalayer.TrackRecordDAO
		sDAO           datalayer.StationDAO
		pathParams     map[string]string
		body           []byte
		expectedResult Worker
		expectedErr    bool
	}{
		// success
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "hitradio-oe3", "timestamp": "1234567890"},
			[]byte("{\"artist\":\"RHCP\",\"title\":\"Californication\"}"),
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3",
					1234567890,
					"track",
					model.Track{"RHCP", "Californication"},
				},
			},
			false,
		},
		// empty station
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "", "timestamp": "1234567890"},
			[]byte("{\"artist\":\"RHCP\",\"title\":\"Californication\"}"),
			nil,
			true,
		},
		// missing station
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"timestamp": "1234567890"},
			[]byte("{\"artist\":\"RHCP\",\"title\":\"Californication\"}"),
			nil,
			true,
		},
		// empty timestamp
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "hitradio-oe3", "timestamp": ""},
			[]byte("{\"artist\":\"RHCP\",\"title\":\"Californication\"}"),
			nil,
			true,
		},
		// missing timestamp
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "hitradio-oe3"},
			[]byte("{\"artist\":\"RHCP\",\"title\":\"Californication\"}"),
			nil,
			true,
		},
		// empty body
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "hitradio-oe3", "timestamp": "1234567890"},
			[]byte(""),
			nil,
			true,
		},
		// invalid body
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "hitradio-oe3", "timestamp": "1234567890"},
			[]byte("I am a sentence and cannot be unmarshalled into a track object."),
			nil,
			true,
		},
		// invalid JSON
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			map[string]string{"station": "hitradio-oe3", "timestamp": "1234567890"},
			[]byte("\"arist\":\"RHCP\",\"invalid\"-\"field\"}"),
			nil,
			true,
		},
	}

	for _, test := range tests {
		result, err := CreateCreateTrackWorker(test.trDAO, test.sDAO, test.pathParams, test.body)
		if (err != nil) != test.expectedErr {
			t.Errorf("CreateCreateTracksWorker(%q, %q, %q, %q): got (%q, %v), expected error: %v",
				test.trDAO, test.sDAO, test.pathParams, test.body, result,
				err, test.expectedErr)
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("CreateCreateTracksWorker(%q, %q, %q, %q): got return type (%v), expected (%q)",
				test.trDAO, test.sDAO, test.pathParams, test.body,
				reflect.TypeOf(result), reflect.TypeOf(test.expectedResult))
			continue
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Errorf("CreateCreateTracksWorker(%q, %q, %q, %q): got \n(%q), expected \n(%q)",
				test.trDAO, test.sDAO, test.pathParams, test.body,
				result, test.expectedResult)
		}
	}
}
