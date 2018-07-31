package shared

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"testing"
	"time"
)

type MockTrackRecordDAO struct{}

func (dao MockTrackRecordDAO) GetTrackRecords(stationId string, start time.Time,
	end time.Time) ([]model.TrackRecord, error) {
	if stationId == "notracksstation" {
		return []model.TrackRecord{}, nil
	}

	trackRecords := []model.TrackRecord{
		{stationId, time.Now().Unix(), model.Track{"RHCP", "Californication"}},
		{stationId, time.Now().Unix(), model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		{stationId, time.Now().Unix(), model.Track{"Cardi B", "I Like It"}},
		{stationId, time.Now().Unix(), model.Track{"RHCP", "Californication"}},
		{stationId, time.Now().Unix(), model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		{stationId, time.Now().Unix(), model.Track{"RHCP", "Californication"}},
	}
	return trackRecords, nil
}

func TestTracks(t *testing.T) {
	date := time.Now().Format("2006-01-02")

	var tests = []struct {
		inputDAO         datalayer.TrackRecordDAO
		inputPathParams  map[string]string
		inputQueryParams map[string]string
		expectedResult   interface{}
		expectedErr      bool
	}{
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"date": date, "filter": "top"},
			model.CountedTracks{
				[]model.CountedTrack{
					{3, model.Track{"RHCP", "Californication"}},
					{2, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
					{1, model.Track{"Cardi B", "I Like It"}},
				},
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"week": date, "filter": "top"},
			model.CountedTracks{
				[]model.CountedTrack{
					{3, model.Track{"RHCP", "Californication"}},
					{2, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
					{1, model.Track{"Cardi B", "I Like It"}},
				},
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"date": date, "filter": "all"},
			model.Tracks{
				[]model.Track{
					{"RHCP", "Californication"},
					{"Jonas Blue, Jack & Jack", "Rise"},
					{"Cardi B", "I Like It"},
				},
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "station-a"},
			map[string]string{"week": date, "filter": "all"},
			model.Tracks{
				[]model.Track{
					{"RHCP", "Californication"},
					{"Jonas Blue, Jack & Jack", "Rise"},
					{"Cardi B", "I Like It"},
				},
			},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "noTracksStation"},
			map[string]string{"date": date, "filter": "top"},
			model.CountedTracks{},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "noTracksStation"},
			map[string]string{"date": date, "filter": "all"},
			model.Tracks{},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "whatever"},
			map[string]string{"date": "2018-07-32", "filter": "all"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "whatever"},
			map[string]string{"week": "2018-07-32", "filter": "all"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "whatever"},
			map[string]string{"date": date, "filter": "invalidFilter"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "whatever"},
			map[string]string{"week": date, "filter": "invalidFilter"},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "noTracksStation"},
			map[string]string{"week": date, "filter": ""},
			model.CountedTracks{},
			false,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "whatever"},
			map[string]string{"week": date},
			nil,
			true,
		},
		{
			MockTrackRecordDAO{},
			map[string]string{"station": "whatever"},
			map[string]string{"filter": "all"},
			nil,
			true,
		},
	}

	for _, test := range tests {
		result, err := Tracks(test.inputDAO, test.inputPathParams, test.inputQueryParams)
		if (err != nil) != test.expectedErr {
			t.Errorf("Tracks(%q, %q, %q): got (%q, %v), expected error: %v",
				test.inputDAO, test.inputPathParams, test.inputQueryParams, result, err,
				test.expectedErr)
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("Tracks(%q, %q, %q): got return type (%s), expected (%s)",
				test.inputDAO, test.inputPathParams, test.inputQueryParams,
				reflect.TypeOf(result).String(), reflect.TypeOf(test.expectedResult).String())
			continue
		}

		if reflect.TypeOf(result) == reflect.TypeOf(model.CountedTracks{}) {
			resultCasted, _ := result.(model.CountedTracks)
			expectedResultCasted, _ := test.expectedResult.(model.CountedTracks)
			for i, expectedCountedTrack := range expectedResultCasted.CountedTracks {
				if resultCasted.CountedTracks[i] != expectedCountedTrack {
					t.Errorf("Tracks(%q, %q, %q): got (%q), expected (%q)",
						test.inputDAO, test.inputPathParams, test.inputQueryParams, result,
						test.expectedResult)
				}
			}
		}

		if reflect.TypeOf(result) == reflect.TypeOf(model.Tracks{}) {
			resultCasted, _ := result.(model.Tracks)
			expectedResultCasted, _ := test.expectedResult.(model.Tracks)
			for _, track := range expectedResultCasted.Tracks {
				match := false
				for _, resultTrack := range resultCasted.Tracks {
					if resultTrack == track {
						match = true
						break
					}
				}
				if !match {
					t.Errorf("Tracks(%q, %q, %q): got (%q), expected (%q)",
						test.inputDAO, test.inputPathParams, test.inputQueryParams, result,
						test.expectedResult)
				}
			}
		}
	}
}

func TestGetStation(t *testing.T) {
	var tests = []struct {
		input       map[string]string
		expectedStr string
		expectedErr bool
	}{
		{map[string]string{"station": "Station-A"}, "station-a", false},
		{map[string]string{}, "", true},
		{map[string]string{"station": ""}, "", true},
	}

	for _, test := range tests {
		result, err := getStation(test.input)
		if (err != nil) != test.expectedErr {
			t.Errorf("getStation(%q): got (%q, %v), expected error: %v",
				test.input, result, err, test.expectedErr)
		}
		if result != test.expectedStr {
			t.Errorf("getStation(%q): got (%q, %v), expected (%q, error: %v)",
				test.input, result, err, test.expectedStr, test.expectedErr)
		}
	}
}

func TestGetFilter(t *testing.T) {
	var tests = []struct {
		input          map[string]string
		expectedFilter Filter
		expectedErr    bool
	}{
		{
			map[string]string{"date": "2018-07-28", "filter": "TOP"},
			Top,
			false,
		},
		{
			map[string]string{"week": "2018-07-28", "filter": "top"},
			Top,
			false,
		},
		{
			map[string]string{"filter": "top"},
			Top,
			false,
		},
		{
			map[string]string{"week": "2018-07-28", "filter": "All"},
			All,
			false,
		},
		{
			map[string]string{"filter": "all"},
			All,
			false,
		},
		{
			map[string]string{"filter": ""},
			Top,
			false,
		},
		{
			map[string]string{"date": "2018-07-28"},
			Top,
			false,
		},
		{
			map[string]string{"filter": "invalidFilter"},
			Err,
			true,
		},
		{
			map[string]string{},
			Top,
			false,
		},
		{
			map[string]string{"unknownFilterParam": "top"},
			Top,
			false,
		},
	}

	for _, test := range tests {
		filter, err := getFilter(test.input)
		if (err != nil) != test.expectedErr {
			t.Errorf("getFilter(%q): got (%q, %v), expected error: %v",
				test.input, filter, err, test.expectedErr)
		}
		if filter != test.expectedFilter {
			t.Errorf("getFilter(%q): got (%q, %v), expected (%q, error: %v)",
				test.input, filter, err, test.expectedFilter, test.expectedErr)
		}
	}
}

func TestCreateDate(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestGetDate: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		input        string
		expectedTime time.Time
		expectedErr  bool
	}{
		{
			"2018-07-28",
			time.Date(2018, 07, 28, 0, 0, 0, 0, loc),
			false,
		},
		{
			"2018-07-32",
			time.Time{},
			true,
		},
		{
			"",
			time.Time{},
			true,
		},
	}

	for _, test := range tests {
		result, err := createDate(test.input)
		if (err != nil) != test.expectedErr {
			t.Errorf("getDate(%q): got (%q, %v), expected error: %v",
				test.input, result, err, test.expectedErr)
		}
		if !result.Equal(test.expectedTime) {
			t.Errorf("getDate(%q): got (%q, %v), expected (%q, error: %v)",
				test.input, result, err, test.expectedTime, test.expectedErr)
		}
	}
}

func TestCreateFirstDateOfWeek(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestGetFirstDayOfWeek: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		input        string
		expectedTime time.Time
		expectedErr  bool
	}{
		{
			"2018-07-28",
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			false,
		},
		{
			"2018-07-23",
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			false,
		},
		{
			"2018-07-32",
			time.Time{},
			true,
		},
		{
			"",
			time.Time{},
			true,
		},
	}

	for _, test := range tests {
		result, err := createFirstDateOfWeek(test.input)
		if (err != nil) != test.expectedErr {
			t.Errorf("getFirstDayOfWeek(%q): got (%q, %v), expected error: %v",
				test.input, result, err, test.expectedErr)
		}
		if !result.Equal(test.expectedTime) {
			t.Errorf("getFirstDayOfWeek(%q): got (%q, %v), expected (%q, error: %v)",
				test.input, result, err, test.expectedTime, test.expectedErr)
		}
	}
}

func TestNormalizeWeekdayNumber(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestNormalizeWeekdayNumber: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		input    time.Time
		expected int
	}{
		{time.Date(2018, 07, 23, 0, 0, 0, 0, loc), 0}, // Mon
		{time.Date(2018, 07, 24, 0, 0, 0, 0, loc), 1}, // Tue
		{time.Date(2018, 07, 25, 0, 0, 0, 0, loc), 2}, // Wed
		{time.Date(2018, 07, 26, 0, 0, 0, 0, loc), 3}, // Thu
		{time.Date(2018, 07, 27, 0, 0, 0, 0, loc), 4}, // Fri
		{time.Date(2018, 07, 28, 0, 0, 0, 0, loc), 5}, // Sat
		{time.Date(2018, 07, 29, 0, 0, 0, 0, loc), 6}, // Sun
	}

	for _, test := range tests {
		normalizedWeekdayNumber := normalizeWeekdayNumber(test.input)
		if normalizedWeekdayNumber != test.expected {
			t.Errorf("normalizeWeekdayNumber(%q): got (%v), expected (%v)",
				test.input, normalizedWeekdayNumber, test.expected)
		}
	}
}

func TestTopTracks(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestTopTracks: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		inputDAO       datalayer.TrackRecordDAO
		inputStation   string
		inputStartDate time.Time
		inputEndDate   time.Time
		expectedTracks []model.CountedTrack
	}{
		{
			MockTrackRecordDAO{},
			"station-a",
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 23, 23, 59, 59, 0, loc),
			[]model.CountedTrack{
				{3, model.Track{"RHCP", "Californication"}},
				{2, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
				{1, model.Track{"Cardi B", "I Like It"}},
			},
		},
		{
			MockTrackRecordDAO{},
			"notracksstation",
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 23, 23, 59, 59, 0, loc),
			[]model.CountedTrack{},
		},
	}

	for _, test := range tests {
		result := topTracks(test.inputDAO, test.inputStation, test.inputStartDate, test.inputEndDate)
		if len(result.CountedTracks) != len(test.expectedTracks) {
			t.Errorf("topTracks(%q, %q, %q, %q): got (%v), expected (%v)",
				test.inputDAO, test.inputStation, test.inputStartDate, test.inputEndDate, result,
				test.expectedTracks)
			continue
		}
		for i, expectedTrack := range test.expectedTracks {
			if result.CountedTracks[i] != expectedTrack {
				t.Errorf("topTracks(%q, %q, %q, %q): got (%v), expected (%v)",
					test.inputDAO, test.inputStation, test.inputStartDate, test.inputEndDate, result,
					test.expectedTracks)
			}
		}
	}
}

func TestAllTracks(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestAllTracks: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		inputDAO       datalayer.TrackRecordDAO
		inputStation   string
		inputStartDate time.Time
		inputEndDate   time.Time
		expectedTracks []model.Track
	}{
		{
			MockTrackRecordDAO{},
			"station-a",
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 29, 23, 59, 59, 0, loc),
			[]model.Track{
				{"RHCP", "Californication"},
				{"Jonas Blue, Jack & Jack", "Rise"},
				{"Cardi B", "I Like It"},
			},
		},
		{
			MockTrackRecordDAO{},
			"notracksstation",
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 29, 23, 59, 59, 0, loc),
			[]model.Track{},
		},
	}

	for _, test := range tests {
		result := allTracks(test.inputDAO, test.inputStation, test.inputStartDate,
			test.inputEndDate)
		if len(result.Tracks) != len(test.expectedTracks) {
			t.Errorf("topTracks(%q, %q, %q, %q): got (%v), expected (%v)",
				test.inputDAO, test.inputStation, test.inputStartDate, test.inputEndDate, result,
				test.expectedTracks)
			continue
		}
		for _, expectedTrack := range test.expectedTracks {
			match := false
			for _, resultTrack := range result.Tracks {
				if resultTrack == expectedTrack {
					match = true
					break
				}
			}
			if !match {
				t.Errorf("topTracks(%q, %q, %q, %q): got (%v), expected (%v)",
					test.inputDAO, test.inputStation, test.inputStartDate, test.inputEndDate, result,
					test.expectedTracks)
			}
		}
	}
}

func TestCalculateDayBoundaries(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestCalcualteDayBoundaries: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		inputDate         time.Time
		expectedStartDate time.Time
		expectedEndDate   time.Time
	}{
		{
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 23, 23, 59, 59, 0, loc),
		},
		{
			time.Date(2018, 07, 25, 12, 44, 23, 1, loc),
			time.Date(2018, 07, 25, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 25, 23, 59, 59, 0, loc),
		},
	}

	for _, test := range tests {
		startDate, endDate := calculateDayBoundaries(test.inputDate)
		if !startDate.Equal(test.expectedStartDate) || !endDate.Equal(test.expectedEndDate) {
			t.Errorf("calculateDayBoundaries(%q): got (%q, %q), expected (%q, %q)",
				test.inputDate, startDate, endDate, test.expectedStartDate, test.expectedEndDate)
		}
	}
}

func TestCalculateWeekBoundaries(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("TestCalcualteDayBoundaries: unable to load location `Europe/Berlin`")
	}

	var tests = []struct {
		inputDate         time.Time
		expectedStartDate time.Time
		expectedEndDate   time.Time
	}{
		{
			time.Date(2018, 07, 23, 18, 57, 23, 77, loc),
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 29, 23, 59, 59, 0, loc),
		},
		{
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 29, 23, 59, 59, 0, loc),
		},
		{
			time.Date(2018, 07, 23, 23, 59, 59, 0, loc),
			time.Date(2018, 07, 23, 0, 0, 0, 0, loc),
			time.Date(2018, 07, 29, 23, 59, 59, 0, loc),
		},
	}

	for _, test := range tests {
		startDate, endDate := calculateWeekBoundaries(test.inputDate)
		if !startDate.Equal(test.expectedStartDate) || !endDate.Equal(test.expectedEndDate) {
			t.Errorf("calculateDayBoundaries(%q): got (%q, %q), expected (%q, %q)",
				test.inputDate, startDate, endDate, test.expectedStartDate, test.expectedEndDate)
		}
	}
}