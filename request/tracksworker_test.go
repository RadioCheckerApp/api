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

func (dao MockTrackRecordDAO) GetMostRecentTrackRecordByStation(stationId string) (model.
	TrackRecord, error) {
	if stationId == "notracksstation" {
		return model.TrackRecord{}, errors.New("no track records in database")
	}
	return model.TrackRecord{stationId, 1234567890, "track", model.Track{"RHCP",
		"Californication"}}, nil
}

func (dao MockTrackRecordDAO) CreateTrackRecord(trackRecord model.TrackRecord) error {
	if trackRecord.Title == "database error" {
		return errors.New("database error")
	}
	return nil
}

type MockTrackRecordDAOLimitTracks struct{}

func (dao MockTrackRecordDAOLimitTracks) GetTrackRecords(start, end time.Time) ([]model.TrackRecord, error) {
	return []model.TrackRecord{}, nil
}

func (dao MockTrackRecordDAOLimitTracks) GetTrackRecordsByStation(stationId string, start time.Time,
	end time.Time) ([]model.TrackRecord, error) {
	if stationId == "withMoreThanTopThree" {
		trackRecords := []model.TrackRecord{
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
			{stationId, time.Now().Unix(), "track", model.Track{"MØ", "Final Song"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		}
		return trackRecords, nil
	}

	if stationId == "withMoreThanTopThreeAndDuplicatedCounters" {
		trackRecords := []model.TrackRecord{
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
			{stationId, time.Now().Unix(), "track", model.Track{"MØ", "Final Song"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		}
		return trackRecords, nil
	}

	if stationId == "withDuplicatedCountersOnly" {
		trackRecords := []model.TrackRecord{
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
			{stationId, time.Now().Unix(), "track", model.Track{"RHCP", "Dani California"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Cardi B", "I Like It"}},
			{stationId, time.Now().Unix(), "track", model.Track{"MØ", "Final Song"}},
			{stationId, time.Now().Unix(), "track", model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		}
		return trackRecords, nil
	}
	return []model.TrackRecord{}, nil
}

func (dao MockTrackRecordDAOLimitTracks) GetMostRecentTrackRecordByStation(stationId string) (model.
	TrackRecord, error) {
	return model.TrackRecord{}, nil
}

func (dao MockTrackRecordDAOLimitTracks) CreateTrackRecord(trackRecord model.TrackRecord) error {
	return nil
}

var countedTracks = model.CountedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.CountedTrack{
		{3, model.Track{"RHCP", "Californication"}},
		{2, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
		{1, model.Track{"Cardi B", "I Like It"}},
	},
}

var tracks = model.Tracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
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
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)

	countedTracks.StartDate = startDate
	countedTracks.EndDate = endDate

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
			model.CountedTracks{
				model.MetaInfo{startDate, endDate},
				[]model.CountedTrack{},
			},
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

		if result.StartDate != test.expectedResult.StartDate ||
			result.EndDate != test.expectedResult.EndDate {
			t.Errorf("(%q).TopTracks(%v, %v): got result startdate: %v / enddate: %v",
				test.worker, test.startDate, test.endDate, result.StartDate, result.EndDate)
		}

		if len(result.CountedTracks) != len(test.expectedResult.CountedTracks) {
			t.Errorf("(%q).TopTracks(%v, %v): got len of result (%q), expected (%q)",
				test.worker, test.startDate, test.endDate, len(result.CountedTracks),
				len(test.expectedResult.CountedTracks))
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

var countedTracksWithMoreThanTopThree = model.CountedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.CountedTrack{
		{5, model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
		{3, model.Track{"RHCP", "Dani California"}},
		{2, model.Track{"Cardi B", "I Like It"}},
	},
}

var countedTracksWithMoreThanTopThreeAndDuplicatedCounters = model.CountedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.CountedTrack{
		{5, model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
		{5, model.Track{"RHCP", "Dani California"}},
		{2, model.Track{"Cardi B", "I Like It"}},
		{1, model.Track{"MØ", "Final Song"}},
		{1, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
	},
}

var countedTracksWithDuplicatedCountersOnly = model.CountedTracks{
	model.MetaInfo{
		time.Now(), // to be defined in the specific tests
		time.Now(), // to be defined in the specific tests
	},
	[]model.CountedTrack{
		{1, model.Track{"RHCP", "The Adventures Of Rain Dance Maggie"}},
		{1, model.Track{"RHCP", "Dani California"}},
		{1, model.Track{"Cardi B", "I Like It"}},
		{1, model.Track{"MØ", "Final Song"}},
		{1, model.Track{"Jonas Blue, Jack & Jack", "Rise"}},
	},
}

func TestTracksWorker_TopTracksLimited(t *testing.T) {
	startDate, endDate := time.Now(), time.Now().AddDate(0, 0, 1)

	countedTracksWithMoreThanTopThree.StartDate = startDate
	countedTracksWithMoreThanTopThree.EndDate = endDate

	countedTracksWithMoreThanTopThreeAndDuplicatedCounters.StartDate = startDate
	countedTracksWithMoreThanTopThreeAndDuplicatedCounters.EndDate = endDate

	countedTracksWithDuplicatedCountersOnly.StartDate = startDate
	countedTracksWithDuplicatedCountersOnly.EndDate = endDate

	var tests = []struct {
		worker         TracksWorker
		expectedResult model.CountedTracks
		expectedErr    bool
	}{
		{
			TracksWorker{MockTrackRecordDAOLimitTracks{}, "withMoreThanTopThree"},
			countedTracksWithMoreThanTopThree,
			false,
		},
		{
			TracksWorker{MockTrackRecordDAOLimitTracks{}, "withMoreThanTopThreeAndDuplicatedCounters"},
			countedTracksWithMoreThanTopThreeAndDuplicatedCounters,
			false,
		},
		{
			TracksWorker{MockTrackRecordDAOLimitTracks{}, "withDuplicatedCountersOnly"},
			countedTracksWithDuplicatedCountersOnly,
			false,
		},
	}

	for _, test := range tests {
		result, err := test.worker.TopTracks(startDate, endDate)
		if (err != nil) != test.expectedErr {
			t.Errorf("(%q).TopTracks(%v, %v): got err (%v), expected err: %v",
				test.worker, startDate, endDate, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if result.StartDate != test.expectedResult.StartDate ||
			result.EndDate != test.expectedResult.EndDate {
			t.Errorf("(%q).TopTracks(%v, %v): got result startdate: %v / enddate: %v",
				test.worker, startDate, endDate, result.StartDate, result.EndDate)
		}

		if len(result.CountedTracks) != len(test.expectedResult.CountedTracks) {
			t.Errorf("(%q).TopTracks(%v, %v): got len of result (%q), expected (%q)",
				test.worker, startDate, endDate, len(result.CountedTracks),
				len(test.expectedResult.CountedTracks))
			continue
		}

		expectedNumberOfTracksPerCounter := make(map[int]int)
		for _, expectedTrack := range test.expectedResult.CountedTracks {
			expectedNumberOfTracksPerCounter[expectedTrack.Counter]++
		}

		gotNumberOfTracksPerCounter := make(map[int]int)
		for _, gotTrack := range result.CountedTracks {
			gotNumberOfTracksPerCounter[gotTrack.Counter]++
		}

		// just check if the number of tracks per counter value are equal
		if !reflect.DeepEqual(expectedNumberOfTracksPerCounter, gotNumberOfTracksPerCounter) {
			t.Errorf("(%q).TopTracks(%v, %v): got number of track per counter: (%q), expected (%q)",
				test.worker, startDate, endDate, gotNumberOfTracksPerCounter, expectedNumberOfTracksPerCounter)
		}
	}
}

func TestTracksWorker_AllTracks(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 1)

	tracks.StartDate = startDate
	tracks.EndDate = endDate

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
			model.Tracks{
				model.MetaInfo{startDate, endDate},
				[]model.Track{},
			},
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

		if result.StartDate != test.expectedResult.StartDate ||
			result.EndDate != test.expectedResult.EndDate {
			t.Errorf("(%q).AllTracks(%v, %v): got result startdate: %v / enddate: %v",
				test.worker, test.startDate, test.endDate, result.StartDate, result.EndDate)
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

func TestTracksWorker_MostRecentTrackRecord(t *testing.T) {
	var tests = []struct {
		worker         TracksWorker
		expectedResult model.TrackRecord
		expectedErr    bool
	}{
		{
			TracksWorker{MockTrackRecordDAO{}, "station-A"},
			model.TrackRecord{
				"station-A",
				1234567890,
				"track",
				model.Track{"RHCP", "Californication"},
			},
			false,
		},
		{
			TracksWorker{MockTrackRecordDAO{}, "notracksstation"},
			model.TrackRecord{},
			true,
		},
	}

	for _, test := range tests {
		result, err := test.worker.MostRecentTrackRecord()
		if (err != nil) != test.expectedErr {
			t.Errorf("(%q).MostRecentTrackRecord(): got err (%v), expected err: %v",
				test.worker, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Errorf("(%q).MostRecentTrackRecord(): result (%q) does not match expected result (%q)",
				test.worker, result, test.expectedResult)
		}
	}
}
