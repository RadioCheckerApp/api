package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"time"
)

type WeekSearchWorker struct {
	SearchWorker
	date time.Time
}

func NewWeekSearchWorker(dao datalayer.TrackRecordDAO, query string,
	date time.Time) (WeekSearchWorker, error) {
	searchWorker, err := NewSearchWorker(dao, query)
	if err != nil {
		return WeekSearchWorker{}, err
	}
	return WeekSearchWorker{searchWorker, date}, nil
}

func (worker WeekSearchWorker) HandleRequest() (interface{}, error) {
	startDate, endDate := calculateWeekBoundaries(worker.date)
	return worker.Search(startDate, endDate)
}
