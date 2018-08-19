package datalayer

import (
	"github.com/RadioCheckerApp/api/model"
	"time"
)

type TrackRecordDAO interface {
	GetTrackRecords(startDate, endDate time.Time) ([]model.TrackRecord, error)
	GetTrackRecordsByStation(station string, startDate, endDate time.Time) ([]model.TrackRecord,
		error)
	CreateTrackRecord(trackRecord model.TrackRecord) error
}
