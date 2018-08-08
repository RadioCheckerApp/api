package request

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
	"reflect"
	"testing"
)

type MockStationDAOSuccess struct{}

func (dao MockStationDAOSuccess) GetAll() ([]model.Station, error) {
	return []model.Station{
		{"kronehit", "Kronehit", "We are the most music", true},
		{"hitradio-oe3", "Hitradio Ö3", "", false},
	}, nil
}

type MockStationDAOSuccessEmpty struct{}

func (dao MockStationDAOSuccessEmpty) GetAll() ([]model.Station, error) {
	return []model.Station{}, nil
}

type MockStationDAOFail struct{}

func (dao MockStationDAOFail) GetAll() ([]model.Station, error) {
	return nil, errors.New("error")
}

func TestNewStationsWorker(t *testing.T) {
	var tests = []struct {
		dao         datalayer.StationDAO
		expectedErr bool
	}{
		{MockStationDAOSuccess{}, false},
		{nil, true},
	}

	for _, test := range tests {
		result, err := NewStationsWorker(test.dao)
		if (err != nil) != test.expectedErr {
			t.Errorf("TestNewStationsWorker(%q): got err (%v), expected err: %v",
				test.dao, err, test.expectedErr)
			continue
		}
		expectedResult := StationsWorker{test.dao}
		if err == nil && !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("TestNewStationsWorker(%q): got result (%v), expected (%v)",
				test.dao, result, expectedResult)
		}
	}
}

func TestStationsWorker_HandleRequest(t *testing.T) {
	var tests = []struct {
		worker         StationsWorker
		expectedResult model.Stations
		expectedErr    bool
	}{
		{
			StationsWorker{MockStationDAOSuccess{}},
			model.Stations{
				[]model.Station{
					{"kronehit", "Kronehit", "We are the most music", true},
					{"hitradio-oe3", "Hitradio Ö3", "", false},
				},
			},
			false,
		},
		{
			StationsWorker{MockStationDAOSuccessEmpty{}},
			model.Stations{[]model.Station{}},
			false,
		},
		{
			StationsWorker{MockStationDAOFail{}},
			model.Stations{},
			true,
		},
	}

	for _, test := range tests {
		result, err := test.worker.HandleRequest()
		if (err != nil) != test.expectedErr {
			t.Errorf("StationsWorker (%v).HandleRequest(): got err (%v), expected err: %v",
				test.worker, err, test.expectedErr)
			continue
		}

		if err != nil {
			// the following tests require a valid result,
			// continue with next test if result was created along with an error
			continue
		}

		if reflect.TypeOf(result) != reflect.TypeOf(test.expectedResult) {
			t.Errorf("StationsWorker (%v).HandleRequest(): got return type (%s), "+
				"expected type (%s)",
				test.worker, reflect.TypeOf(result).String(), reflect.TypeOf(test.expectedResult).String())
			continue
		}

		if !reflect.DeepEqual(result, test.expectedResult) {
			t.Errorf("StationsWorker (%v).HandleRequest(): got (%v, %v), expect (%v, %v)",
				test.worker, result, err, test.expectedResult, test.expectedErr)
		}
	}
}
