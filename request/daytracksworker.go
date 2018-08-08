package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"time"
)

type DayTracksWorker struct {
	TracksWorker
	date   time.Time
	filter Filter
}

func NewDayTracksWorker(dao datalayer.TrackRecordDAO, station string, date time.Time,
	filter Filter) (DayTracksWorker, error) {
	tracksWorker, err := NewTracksWorker(dao, station)
	if err != nil {
		return DayTracksWorker{}, err
	}
	return DayTracksWorker{tracksWorker, date, filter}, nil
}

func (worker DayTracksWorker) HandleRequest() (interface{}, error) {
	startDate, endDate := calculateDayBoundaries(worker.date)
	if worker.filter == Top {
		return worker.TopTracks(startDate, endDate)
	}
	return worker.AllTracks(startDate, endDate)
}

func calculateDayBoundaries(date time.Time) (time.Time, time.Time) {
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0,
		date.Location())
	endDate := startDate.AddDate(0, 0, 1).Add(-1 * time.Second)
	return startDate, endDate
}
