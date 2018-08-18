package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewWeekSearchWorker(t *testing.T) {
	var tests = []struct {
		dao         datalayer.TrackRecordDAO
		query       string
		date        time.Time
		expectedErr bool
	}{
		{
			MockTrackRecordDAO{},
			"RHCP+maggie",
			time.Now(),
			false,
		},
		{nil, "RHCP+maggie", time.Now(), true},
		{MockTrackRecordDAO{}, "", time.Now(), true},
	}

	for _, test := range tests {
		result, err := NewWeekSearchWorker(test.dao, test.query, test.date)
		if (err != nil) != test.expectedErr {
			t.Errorf("NewWeekSearchWorker(%q, %q, %q): got err (%v), expected err: %v",
				test.dao, test.query, test.date, err, test.expectedErr)
			continue
		}
		expectedResult := WeekSearchWorker{
			SearchWorker{test.dao, strings.Split(strings.ToLower(test.query), queryStrKeywordsSeparator)},
			test.date,
		}
		if err == nil && !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("NewWeekSearchWorker(%q, %q, %q): got result (%v), expected (%v)",
				test.dao, test.query, test.date, result, expectedResult)
		}
	}
}

func TestWeekSearchWorker_HandleRequest(t *testing.T) {
	var date = time.Now()

	var tests = []struct {
		worker         WeekSearchWorker
		expectedResult model.MatchedTracks
		expectedErr    bool
	}{
		{
			WeekSearchWorker{SearchWorker{MockTrackRecordDAO{}, []string{"californication"}}, date},
			matchedTracks0,
			false,
		},
		{
			WeekSearchWorker{SearchWorker{MockTrackRecordDAO{}, []string{"cali"}}, date},
			matchedTracks1,
			false,
		},
		{
			WeekSearchWorker{SearchWorker{MockTrackRecordDAO{}, []string{"maggie", "rhcp"}}, date},
			matchedTracks2,
			false,
		},
		{
			WeekSearchWorker{SearchWorker{MockTrackRecordDAO{}, []string{"ø"}}, date},
			matchedTracks3,
			false,
		},
		{
			WeekSearchWorker{SearchWorker{MockTrackRecordDAO{}, []string{"no", "tracks", "query"}}, date},
			model.MatchedTracks{},
			false,
		},
		{
			WeekSearchWorker{SearchWorker{MockTrackRecordDAOWeekVerifier{}, []string{"nevermind"}},
				date},
			model.MatchedTracks{},
			false,
		},
	}

	for _, test := range tests {
		result, err := test.worker.HandleRequest()
		if (err != nil) != test.expectedErr {
			t.Errorf("(%v).HandleRequest(): got err (%v), expected err: %v",
				test.worker, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("(%v).HandleRequest(): got return type (%s), expected type (%s)",
				test.worker, reflect.TypeOf(result).String(), reflect.TypeOf(test.expectedResult).String())
			continue
		}

		resultCasted, _ := result.(model.MatchedTracks)

		if len(resultCasted.MatchedTracks) != len(test.expectedResult.MatchedTracks) {
			t.Errorf("(%v).HandleRequest(): got result (%q), expected (%q)",
				test.worker, resultCasted, test.expectedResult)
		}

		for _, track := range test.expectedResult.MatchedTracks {
			match := false
			for _, resultTrack := range resultCasted.MatchedTracks {
				if reflect.DeepEqual(resultTrack, track) {
					match = true
					break
				}
			}
			if !match {
				t.Errorf("(%v).HandleRequest(): expected item (%q) is not element of result (%q)",
					test.worker, track, test.expectedResult)
			}
		}
	}
}
