package datalayer

import (
	"github.com/RadioCheckerApp/api/model"
	"time"
)

type DDBTrackRecordDAO struct {
	dynamoDB  *DynamoDB
	tableName string
}

func NewDDBTrackRecordDAO(dynamodb *DynamoDB, tableName string) *DDBTrackRecordDAO {
	return &DDBTrackRecordDAO{dynamodb, tableName}
}

func (dao *DDBTrackRecordDAO) GetTrackRecords(station string, startDate, endDate time.Time) ([]model.TrackRecord, error) {
	return []model.TrackRecord{}, nil
}
