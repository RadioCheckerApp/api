package request

import (
	"errors"
	"github.com/RadioCheckerApp/api/datalayer"
	"github.com/RadioCheckerApp/api/model"
)

type StationsWorker struct {
	dao datalayer.StationDAO
}

func NewStationsWorker(dao datalayer.StationDAO) (StationsWorker, error) {
	if dao == nil {
		return StationsWorker{}, errors.New("dao must not be nil")
	}
	return StationsWorker{dao}, nil
}

func (worker StationsWorker) HandleRequest() (interface{}, error) {
	stations, err := worker.dao.GetAll()
	return model.Stations{stations}, err
}
