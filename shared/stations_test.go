package shared

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

func TestStations(t *testing.T) {
	var tests = []struct {
		inputDAO         datalayer.StationDAO
		expectedStations model.Stations
		expectedErr      error
	}{
		{
			MockStationDAOSuccess{},
			model.Stations{
				[]model.Station{
					{"kronehit", "Kronehit", "We are the most music", true},
					{"hitradio-oe3", "Hitradio Ö3", "", false},
				},
			},
			nil,
		},
		{
			MockStationDAOSuccessEmpty{},
			model.Stations{[]model.Station{}},
			nil,
		},
		{
			MockStationDAOFail{},
			model.Stations{},
			nil,
		},
	}

	for _, test := range tests {
		stations, err := Stations(test.inputDAO)
		if !reflect.DeepEqual(stations, test.expectedStations) || err != test.expectedErr {
			t.Errorf("Stations(%s): got (%v, %v), expect (%v, %v)",
				reflect.TypeOf(test.inputDAO).Name(), stations, err, test.expectedStations, test.expectedErr)
		}
	}
}
