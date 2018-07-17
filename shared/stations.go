package shared

import (
	"encoding/json"
	"github.com/RadioCheckerApp/api/datalayer"
)

func Stations(dao datalayer.StationDAO) (string, error) {
	stations, err := dao.GetAll()
	if err != nil {
		return "[]", err
	}
	jsonBytes, err := json.Marshal(stations)
	return string(jsonBytes), err
}
