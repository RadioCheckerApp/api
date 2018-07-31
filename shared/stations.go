package shared

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
)

func Stations(dao datalayer.StationDAO) (model.Stations, error) {
	stations, err := dao.GetAll()
	if err != nil {
		return model.Stations{}, nil
	}
	return model.Stations{stations}, nil
}
