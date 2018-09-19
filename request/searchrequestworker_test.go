package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewSearchWorker(t *testing.T) {
	var tests = []struct {
		dao         datalayer.TrackRecordDAO
		queryStr    string
		expectedErr bool
	}{
		{MockTrackRecordDAO{}, "The+Adventures+Of+Rain+Dance+Maggie", false},
		{nil, "The+Adventures+Of+Rain+Dance+Maggie", true},
		{MockTrackRecordDAO{}, "", true},
	}

	for _, test := range tests {
		result, err := NewSearchWorker(test.dao, test.queryStr)
		if (err != nil) != test.expectedErr {
			t.Errorf("TestNewSearchWorker(%q, %q): got err (%v), expected err: %v",
				test.dao, test.queryStr, err, test.expectedErr)
			continue
		}
		expectedResult := SearchWorker{
			test.dao,
			strings.Split(strings.ToLower(test.queryStr), queryStrKeywordsSeparator),
		}
		if err == nil && !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("TestNewSearchWorker(%q, %q): got result (%v), expected (%v)",
				test.dao, test.queryStr, result, expectedResult)
		}
	}
}

var matchedTracks0 = model.MatchedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.MatchedTrack{
		{
			map[string]int{"station-a": 3, "station-b": 1},
			model.Track{"RHCP", "Californication"},
		},
	},
}

var matchedTracks1 = model.MatchedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.MatchedTrack{
		{
			map[string]int{"station-a": 3, "station-b": 1},
			model.Track{"RHCP", "Californication"},
		},
		{
			map[string]int{"station-a": 0, "station-b": 1},
			model.Track{"RHCP", "Dani California"},
		},
	},
}

var matchedTracks2 = model.MatchedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.MatchedTrack{
		{
			map[string]int{"station-a": 3, "station-b": 1, "station-c": 0},
			model.Track{"RHCP", "Californication"},
		},
		{
			map[string]int{"station-a": 0, "station-b": 1, "station-c": 0},
			model.Track{"RHCP", "Dani California"},
		},
		{
			map[string]int{"station-a": 0, "station-b": 1, "station-c": 2},
			model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"},
		},
	},
}

var matchedTracks3 = model.MatchedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.MatchedTrack{
		{
			map[string]int{"station-b": 1},
			model.Track{"MØ", "Final Song"},
		},
	},
}

func TestSearchWorker_Search(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)

	matchedTracks0.StartDate = startDate
	matchedTracks0.EndDate = endDate

	matchedTracks1.StartDate = startDate
	matchedTracks1.EndDate = endDate

	matchedTracks2.StartDate = startDate
	matchedTracks2.EndDate = endDate

	matchedTracks3.StartDate = startDate
	matchedTracks3.EndDate = endDate

	var tests = []struct {
		worker         SearchWorker
		startDate      time.Time
		endDate        time.Time
		expectedResult model.MatchedTracks
		expectedErr    bool
	}{
		{
			SearchWorker{MockTrackRecordDAO{}, []string{"californication"}},
			startDate,
			endDate,
			matchedTracks0,
			false,
		},
		{
			SearchWorker{MockTrackRecordDAO{}, []string{"cali"}},
			startDate,
			endDate,
			matchedTracks1,
			false,
		},
		{
			SearchWorker{MockTrackRecordDAO{}, []string{"maggie", "rhcp"}},
			startDate,
			endDate,
			matchedTracks2,
			false,
		},
		{
			SearchWorker{MockTrackRecordDAO{}, []string{"ø"}},
			startDate,
			endDate,
			matchedTracks3,
			false,
		},
		{
			SearchWorker{MockTrackRecordDAO{}, []string{"no", "tracks", "query"}},
			startDate,
			endDate,
			model.MatchedTracks{
				model.MetaInfo{startDate, endDate},
				[]model.MatchedTrack{},
			},
			false,
		},
		{
			SearchWorker{MockTrackRecordDAO{}, []string{""}},
			endDate,
			startDate,
			model.MatchedTracks{
				model.MetaInfo{startDate, endDate},
				[]model.MatchedTrack{},
			},
			true,
		},
	}

	for _, test := range tests {
		result, err := test.worker.Search(test.startDate, test.endDate)
		if (err != nil) != test.expectedErr {
			t.Errorf("(%q).Search(%v, %v): got err (%v), expected err: %v",
				test.worker, test.startDate, test.endDate, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if result.StartDate != test.expectedResult.StartDate ||
			result.EndDate != test.expectedResult.EndDate {
			t.Errorf("(%q).Search(%v, %v): got result startdate: %v / enddate: %v",
				test.worker, test.startDate, test.endDate, result.StartDate, result.EndDate)
		}

		if len(result.MatchedTracks) != len(test.expectedResult.MatchedTracks) {
			t.Errorf("(%q).Search(%v, %v): got result length (%d), expected (%d)",
				test.worker, test.startDate, test.endDate, len(result.MatchedTracks),
				len(test.expectedResult.MatchedTracks))
		}

		for _, expectedMatchedTrack := range test.expectedResult.MatchedTracks {
			match := false
			for _, resultTrack := range result.MatchedTracks {
				if reflect.DeepEqual(resultTrack, expectedMatchedTrack) {
					match = true
					break
				}
			}
			if !match {
				t.Errorf("(%q).Search(%v, %v): expected item (%q) is not element of result (%q)",
					test.worker, test.startDate, test.endDate, expectedMatchedTrack,
					test.expectedResult.MatchedTracks)
			}
		}
	}
}
