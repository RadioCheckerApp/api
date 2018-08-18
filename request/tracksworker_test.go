package request

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"testing"
	"time"
)

type MockTrackRecordDAO struct{}

func (dao MockTrackRecordDAO) GetTrackRecords(start, end time.Time) ([]model.TrackRecord, error) {
	return dao.GetTrackRecordsByStation("getTrackRecords", start, end)
}

func (dao MockTrackRecordDAO) GetTrackRecordsByStation(stationId string, start time.Time,
	end time.Time) ([]model.TrackRecord, error) {
	if stationId == "notracksstation" {
		return []model.TrackRecord{}, nil
	}

	if start.After(end) {
		return []model.TrackRecord{}, errors.New("error")
	}

	if stationId == "getTrackRecords" {
		trackRecords := []model.TrackRecord{
			{"station-a", time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
			{"station-a", time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
			{"station-a", time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
			{"station-a", time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
			{"station-a", time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
			{"station-a", time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
			{"station-b", time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
			{"station-b", time.Now().Unix(), "track", model.Track{"MØ", "Final Song"}},
			{"station-b", time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{"station-b", time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{"station-c", time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{"station-c", time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
		}
		return trackRecords, nil
	}

	trackRecords := []model.TrackRecord{
		{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
		{stationId, time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		{stationId, time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
		{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
		{stationId, time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
	}
	return trackRecords, nil
}

var countedTracks = model.CountedTracks{
	[]model.CountedTrack{
		{3, model.Track{"RHCP", "Californication"}},
		{2, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		{1, model.Track{"Cardi B", "I Like It"}},
	},
}

var tracks = model.Tracks{
	[]model.Track{
		{"RHCP", "Californication"},
		{"Jonas Blue, Jack & Jack", "Rise"},
		{"Cardi B", "I Like It"},
	},
}

func TestNewTracksWorker(t *testing.T) {
	var tests = []struct {
		dao         datalayer.TrackRecordDAO
		station     string
		expectedErr bool
	}{
		{MockTrackRecordDAO{}, "teststation", false},
		{nil, "teststation", true},
		{MockTrackRecordDAO{}, "", true},
	}

	for _, test := range tests {
		result, err := NewTracksWorker(test.dao, test.station)
		if (err != nil) != test.expectedErr {
			t.Errorf("TestNewTracksWorker(%q, %q): got err (%v), expected err: %v",
				test.dao, test.station, err, test.expectedErr)
			continue
		}
		expectedResult := TracksWorker{test.dao, test.station}
		if err == nil && !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("TestNewTracksWorker(%q, %q): got result (%v), expected (%v)",
				test.dao, test.station, result, expectedResult)
		}
	}
}

func TestTracksWorker_TopTracks(t *testing.T) {
	var startDate = time.Now()
	var endDate = startDate.AddDate(0, 0, 1)

	var tests = []struct {
		worker         TracksWorker
		startDate      time.Time
		endDate        time.Time
		expectedResult model.CountedTracks
		expectedErr    bool
	}{
		{
			TracksWorker{MockTrackRecordDAO{}, "station-A"},
			startDate,
			endDate,
			countedTracks,
			false,
		},
		{
			TracksWorker{MockTrackRecordDAO{}, "notracksstation"},
			startDate,
			endDate,
			model.CountedTracks{},
			false,
		},
		{
			TracksWorker{MockTrackRecordDAO{}, "errorstation"},
			endDate,
			startDate,
			model.CountedTracks{},
			true,
		},
	}

	for _, test := range tests {
		result, err := test.worker.TopTracks(test.startDate, test.endDate)
		if (err != nil) != test.expectedErr {
			t.Errorf("(%q).TopTracks(%v, %v): got err (%v), expected err: %v",
				test.worker, test.startDate, test.endDate, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		for i, expectedCountedTrack := range test.expectedResult.CountedTracks {
			if !reflect.DeepEqual(result.CountedTracks[i], expectedCountedTrack) {
				t.Errorf("(%q).TopTracks(%v, %v): got result (%q), expected (%q)",
					test.worker, test.startDate, test.endDate, result, test.expectedResult)
			}
		}
	}
}

func TestTracksWorker_AllTracks(t *testing.T) {
	var startDate = time.Now()
	var endDate = startDate.AddDate(0, 0, 1)

	var tests = []struct {
		worker         TracksWorker
		startDate      time.Time
		endDate        time.Time
		expectedResult model.Tracks
		expectedErr    bool
	}{
		{
			TracksWorker{MockTrackRecordDAO{}, "station-A"},
			startDate,
			endDate,
			tracks,
			false,
		},
		{
			TracksWorker{MockTrackRecordDAO{}, "notracksstation"},
			startDate,
			endDate,
			model.Tracks{},
			false,
		},
		{
			TracksWorker{MockTrackRecordDAO{}, "errorstation"},
			endDate,
			startDate,
			tracks,
			true,
		},
	}

	for _, test := range tests {
		result, err := test.worker.AllTracks(test.startDate, test.endDate)
		if (err != nil) != test.expectedErr {
			t.Errorf("(%q).AllTracks(%v, %v): got err (%v), expected err: %v",
				test.worker, test.startDate, test.endDate, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if len(result.Tracks) != len(test.expectedResult.Tracks) {
			t.Errorf("(%q).AllTracks(%v, %v): got result (%q), expected (%q)",
				test.worker, test.startDate, test.endDate, result, test.expectedResult)
		}

		for _, track := range test.expectedResult.Tracks {
			match := false
			for _, resultTrack := range result.Tracks {
				if resultTrack == track {
					match = true
					break
				}
			}
			if !match {
				t.Errorf("(%q).AllTracks(%v, %v): expected item (%q) is not element of result (%q)",
					test.worker, test.startDate, test.endDate, track, result)
			}
		}
	}
}
