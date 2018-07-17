package datalayer

import "github.com/RadioCheckerApp/api/model"

type StationDAO interface {
	GetAll() ([]model.Station, error)
}
