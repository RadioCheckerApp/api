package datalayer

import (
	"github.com/RadioCheckerApp/api/model"
	"time"
)

type TrackRecordDAO interface {
	GetTrackRecords(station string, startDate, endDate time.Time) ([]model.TrackRecord, error)
}
