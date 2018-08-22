package request

import (
	"fmt"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"testing"
	"time"
)

func TestNewCreateTrackWorker(t *testing.T) {
	var tests = []struct {
		trDAO       datalayer.TrackRecordDAO
		sDAO        datalayer.StationDAO
		trackRecord model.TrackRecord
		expectedErr bool
	}{
		{
			MockTrackRecordDAO{},
			MockStationDAOSuccess{},
			model.TrackRecord{
				"station-a",
				time.Now().Unix(),
				"track",
				model.Track{"RHCP", "Californication"}},
			false,
		},
		{nil, MockStationDAOSuccess{}, model.TrackRecord{}, true},
		{MockTrackRecordDAO{}, nil, model.TrackRecord{}, true},
	}

	for _, test := range tests {
		result, err := NewCreateTrackWorker(test.trDAO, test.sDAO, test.trackRecord)
		if (err != nil) != test.expectedErr {
			t.Errorf("NewCreateTrackWorker(%q, %q, %q): got err (%v), expected err: %v",
				test.trDAO, test.sDAO, test.trackRecord, err, test.expectedErr)
			continue
		}
		expectedResult := CreateTrackWorker{test.trDAO, test.sDAO, test.trackRecord}
		if err == nil && !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("NewDaySearchWorker(%q, %q, %q): got result (%v), expected (%v)",
				test.trDAO, test.sDAO, test.trackRecord, result, expectedResult)
		}
	}
}

func TestCreateTrackWorker_HandleRequest(t *testing.T) {
	var timestamp = time.Now().Unix()

	var tests = []struct {
		worker         CreateTrackWorker
		expectedResult string
		expectedErr    bool
	}{
		// cache empty
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccessEmpty{},
				model.TrackRecord{
					"hitradio-oe3", timestamp,
					"track",
					model.Track{"RHCP", "Californication"},
				},
			},
			"ignored",
			true, // cache empty & MockStationDAOSuccessEmpty serves no stations
		},
		// cache filled
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3", timestamp,
					"track",
					model.Track{"RHCP", "Californication"},
				},
			},
			"track created: /stations/hitradio-oe3/tracks/" + fmt.Sprintf("%d", timestamp),
			false,
		},
		// database error
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3", timestamp,
					"track",
					model.Track{"CAUTION:", "DATABASE ERROR"},
				},
			},
			"ignored",
			true,
		},
		// invalid station
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"invalid station", timestamp,
					"track",
					model.Track{"RHCP", "Californication"},
				},
			},
			"ignored",
			true,
		},
		// timestamp too far in future
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3", time.Now().Add(31 * time.Minute).Unix(),
					"track",
					model.Track{"RHCP", "Californication"},
				},
			},
			"ignored",
			true,
		},
		// invalid type
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3", timestamp,
					"invalid type",
					model.Track{"RHCP", "Californication"},
				},
			},
			"ignored",
			true,
		},
		// invalid artist
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3", timestamp,
					"track",
					model.Track{"", "Californication"},
				},
			},
			"ignored",
			true,
		},
		// invalid title
		{
			CreateTrackWorker{
				MockTrackRecordDAO{},
				MockStationDAOSuccess{},
				model.TrackRecord{
					"hitradio-oe3", timestamp,
					"track",
					model.Track{"RHCP", ""},
				},
			},
			"ignored",
			true,
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

		resultCasted, _ := result.(string)

		if resultCasted != test.expectedResult {
			t.Errorf("(%v).HandleRequest(): got result (%q), expected (%q)",
				test.worker, resultCasted, test.expectedResult)
		}
	}
}
