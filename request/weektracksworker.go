package request

import (
	"github.com/RadioCheckerApp/api/datalayer"
	"log"
	"time"
)

type WeekTracksWorker struct {
	TracksWorker
	date   time.Time
	filter Filter
}

func NewWeekTracksWorker(dao datalayer.TrackRecordDAO, station string, date time.Time,
	filter Filter) (WeekTracksWorker, error) {
	tracksWorker, err := NewTracksWorker(dao, station)
	if err != nil {
		return WeekTracksWorker{}, err
	}
	return WeekTracksWorker{tracksWorker, date, filter}, nil
}

func (worker WeekTracksWorker) HandleRequest() (interface{}, error) {
	startDate, endDate := calculateWeekBoundaries(worker.date)
	if worker.filter == Top {
		return worker.TopTracks(startDate, endDate)
	}
	return worker.AllTracks(startDate, endDate)
}

func calculateWeekBoundaries(date time.Time) (time.Time, time.Time) {
	startDate := calculateFirstDateOfWeek(date)
	endDate := startDate.AddDate(0, 0, 7).Add(-1 * time.Second)
	return startDate, endDate
}

func calculateFirstDateOfWeek(date time.Time) time.Time {
	dateWithoutTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, getLocation())
	return dateWithoutTime.AddDate(0, 0, -normalizeWeekdayNumber(dateWithoutTime))
}

func getLocation() *time.Location {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal("unable to load timezone location `Europe/Berlin`")
	}
	return location
}

func normalizeWeekdayNumber(date time.Time) int {
	// Sunday = 0, ..., Saturday = 6
	usWeekdayNumber := date.Weekday()
	return int((usWeekdayNumber + 6) % 7)
}
