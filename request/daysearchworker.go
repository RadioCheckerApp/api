package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"time"
)

type DaySearchWorker struct {
	SearchWorker
	date time.Time
}

func NewDaySearchWorker(dao datalayer.TrackRecordDAO, query string,
	date time.Time) (DaySearchWorker, error) {
	searchWorker, err := NewSearchWorker(dao, query)
	if err != nil {
		return DaySearchWorker{}, err
	}
	return DaySearchWorker{searchWorker, date}, nil
}

func (worker DaySearchWorker) HandleRequest() (interface{}, error) {
	startDate, endDate := calculateDayBoundaries(worker.date)
	return worker.Search(startDate, endDate)
}
