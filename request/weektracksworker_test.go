package request

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"testing"
	"time"
)

type MockTrackRecordDAOWeekVerifier struct{}

func (dao MockTrackRecordDAOWeekVerifier) GetTrackRecords(start,
	end time.Time) ([]model.TrackRecord, error) {
	return nil, nil
}

func (dao MockTrackRecordDAOWeekVerifier) GetTrackRecordsByStation(stationId string, start,
	end time.Time) ([]model.TrackRecord, error) {
	expectedStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	if !start.Equal(expectedStart) {
		return nil, errors.New("invalid start time")
	}

	if start.Weekday() != time.Monday {
		return nil, errors.New("start date is not a monday")
	}

	expectedEndDay := start.AddDate(0, 0, 6)
	expectedEnd := time.Date(expectedEndDay.Year(), expectedEndDay.Month(), expectedEndDay.Day(), 23, 59,
		59, 0, start.Location())
	if !end.Equal(expectedEnd) {
		return nil, errors.New("invalid end time")
	}

	if end.Weekday() != time.Sunday {
		return nil, errors.New("end date is not a sunday")
	}

	return []model.TrackRecord{}, nil
}

func (dao MockTrackRecordDAOWeekVerifier) CreateTrackRecord(trackRecord model.TrackRecord) error {
	return nil
}

func TestNewWeekTracksWorker(t *testing.T) {
	var tests = []struct {
		dao         datalayer.TrackRecordDAO
		station     string
		date        time.Time
		filter      Filter
		expectedErr bool
	}{
		{
			MockTrackRecordDAO{},
			"teststation",
			time.Now(),
			Top,
			false,
		},
		{nil, "teststation", time.Now(), All, true},
		{MockTrackRecordDAO{}, "", time.Now(), Top, true},
	}

	for _, test := range tests {
		result, err := NewWeekTracksWorker(test.dao, test.station, test.date, test.filter)
		if (err != nil) != test.expectedErr {
			t.Errorf("TestWeekDayTracksWorker(%q, %q, %q, %q): got err (%v), expected err: %v",
				test.dao, test.station, test.date, test.filter, err, test.expectedErr)
			continue
		}
		expectedResult := WeekTracksWorker{TracksWorker{test.dao, test.station}, test.date,
			test.filter}
		if err == nil && !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("TestWeekDayTracksWorker(%q, %q, %q, %q): got result (%v), expected (%v)",
				test.dao, test.station, test.date, test.filter, result, expectedResult)
		}
	}
}

func TestWeekTracksWorker_HandleRequest_TopTracks(t *testing.T) {
	var date = time.Now()

	var tests = []struct {
		worker         WeekTracksWorker
		expectedResult model.CountedTracks
		expectedErr    bool
	}{
		{
			WeekTracksWorker{TracksWorker{MockTrackRecordDAO{}, "station-A"}, date, Top},
			countedTracks,
			false,
		},
		{
			WeekTracksWorker{TracksWorker{MockTrackRecordDAO{}, "notracksstation"}, date,
				Top},
			model.CountedTracks{},
			false,
		},
		{
			WeekTracksWorker{TracksWorker{MockTrackRecordDAOWeekVerifier{}, "nevermind"}, date,
				Top},
			model.CountedTracks{},
			false,
		},
	}

	for _, test := range tests {
		result, err := test.worker.HandleRequest()
		if (err != nil) != test.expectedErr {
			t.Errorf("WeekTracksWorker (%v).HandleRequest(): got err (%v), expected err: %v",
				test.worker, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("WeekTracksWorker (%v).HandleRequest(): got return type (%s), "+
				"expected type (%s)",
				test.worker, reflect.TypeOf(result).String(), reflect.TypeOf(test.expectedResult).String())
			continue
		}

		resultCasted, _ := result.(model.CountedTracks)
		for i, expectedCountedTrack := range test.expectedResult.CountedTracks {
			if resultCasted.CountedTracks[i] != expectedCountedTrack {
				t.Errorf("WeekTracksWorker (%v).HandleRequest(): expected (%q) at pos #%d, "+
					"got (%q)",
					test.worker, expectedCountedTrack, i, resultCasted.CountedTracks[i])
			}
		}
	}
}

func TestWeekTracksWorker_HandleRequest_AllTracks(t *testing.T) {
	var date = time.Now()

	var tests = []struct {
		worker         WeekTracksWorker
		expectedResult model.Tracks
		expectedErr    bool
	}{
		{
			WeekTracksWorker{TracksWorker{MockTrackRecordDAO{}, "station-A"}, date, All},
			tracks,
			false,
		},
		{
			WeekTracksWorker{TracksWorker{MockTrackRecordDAO{}, "notracksstation"}, date, All},
			model.Tracks{},
			false,
		},
		{
			WeekTracksWorker{TracksWorker{MockTrackRecordDAOWeekVerifier{}, "nevermind"}, date,
				All},
			model.Tracks{},
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

		resultCasted, _ := result.(model.Tracks)

		if len(resultCasted.Tracks) != len(test.expectedResult.Tracks) {
			t.Errorf("(%v).HandleRequest(): got result (%q), expected (%q)",
				test.worker, resultCasted, test.expectedResult)
		}

		for _, track := range test.expectedResult.Tracks {
			match := false
			for _, resultTrack := range resultCasted.Tracks {
				if resultTrack == track {
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
