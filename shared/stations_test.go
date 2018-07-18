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
		input   datalayer.StationDAO
		jsonStr string
		err     error
	}{
		{
			MockStationDAOSuccess{},
			"[{\"stationId\":\"kronehit\",\"name\":\"Kronehit\",\"description\":" +
				"\"We are the most music\",\"active\":true},{\"stationId\":\"hitradio-oe3\"," +
				"\"name\":\"Hitradio Ö3\",\"description\":\"\",\"active\":false}]",
			nil,
		},
		{
			MockStationDAOSuccessEmpty{},
			"[]",
			nil,
		},
		{
			MockStationDAOFail{},
			"[]",
			errors.New("error"),
		},
	}

	for _, test := range tests {
		jsonStr, err := Stations(test.input)
		if jsonStr != test.jsonStr || (err != err && err.Error() != test.err.Error()) {
			t.Errorf("Stations(%s): got (%q, %q), expect (%q, %q)",
				reflect.TypeOf(test.input).Name(), jsonStr, err, test.jsonStr, test.err)
		}
	}
}
