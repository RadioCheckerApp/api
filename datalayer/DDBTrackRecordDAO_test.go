package datalayer

import (
	"errors"
	"github.com/RadioCheckerApp/api/model"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
	"testing"
	"time"
)

type MockDynamoDB struct{}

func (ddb MockDynamoDB) ScanPages(input *dynamodb.ScanInput, fn func(*dynamodb.ScanOutput,
	bool) bool) error {
	return nil
}

func (ddb MockDynamoDB) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if input == nil {
		return nil, errors.New("input must not be nil")
	}

	if input.TableName == nil {
		return nil, errors.New("TableName must not be nil")
	}

	if input.KeyConditionExpression == nil {
		return nil, errors.New("KeyConditionException must not be nil")
	}

	if len(strings.Split(*input.KeyConditionExpression, " ")) != 9 {
		return nil, errors.New("KeyConditionExpression must exist of 9 words")
	}

	if input.ExpressionAttributeNames == nil {
		return nil, errors.New("ExpressionAttributeNames must not be nil")
	}

	if name, ok := input.ExpressionAttributeNames["#sid"]; !ok || *name != "stationId" {
		return nil, errors.New("ExpressionAttributeNames must contain mapping `#sid`: `stationId`")
	}

	if name, ok := input.ExpressionAttributeNames["#t"]; !ok || *name != "type" {
		return nil, errors.New("ExpressionAttributeNames must contain mapping `#t`: `type`")
	}

	if input.ExpressionAttributeValues == nil {
		return nil, errors.New("ExpressionAttributeValues must not be nil")
	}

	if _, ok := input.ExpressionAttributeValues[":stationId"]; !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:stationId`")
	}

	if _, ok := input.ExpressionAttributeValues[":lowerBound"]; !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:lowerBound`")
	}

	if _, ok := input.ExpressionAttributeValues[":upperBound"]; !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:upperBound`")
	}

	if _, ok := input.ExpressionAttributeValues[":type"]; !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:type`")
	}

	stationIds := []string{"station-a", "station-b"}
	timestamps := []string{"1532897851", "1532897892"}
	recordType := "track"
	artists := []string{"Mø", "Jack Ü ft. Skrillex & Diplo"}
	titles := []string{"Final Song", "Where Are Ü Now"}

	output := &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"stationId": &dynamodb.AttributeValue{S: &stationIds[0]},
				"airtime":   &dynamodb.AttributeValue{N: &timestamps[0]},
				"type":      &dynamodb.AttributeValue{S: &recordType},
				"artist":    &dynamodb.AttributeValue{S: &artists[0]},
				"title":     &dynamodb.AttributeValue{S: &titles[0]},
			},
			{
				"stationId": &dynamodb.AttributeValue{S: &stationIds[1]},
				"airtime":   &dynamodb.AttributeValue{N: &timestamps[1]},
				"type":      &dynamodb.AttributeValue{S: &recordType},
				"artist":    &dynamodb.AttributeValue{S: &artists[1]},
				"title":     &dynamodb.AttributeValue{S: &titles[1]},
			},
		},
	}

	return output, nil
}

func TestDDBTrackRecordDAO_GetTrackRecordsSuccess(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDB{}, "testTable")
	station := "station-a"
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now()

	expectedTrackRecords := []model.TrackRecord{
		{"station-a", 1532897851, "track", model.Track{"Mø", "Final Song"}},
		{"station-b", 1532897892, "track", model.Track{"Jack Ü ft. Skrillex & Diplo",
			"Where Are Ü Now"}},
	}

	trackRecords, err := trackRecordDAO.GetTrackRecords(station, startDate, endDate)

	if err != nil {
		t.Errorf("GetTrackRecords(%q, %q, %q): got (%q, %v), expected (%q, nil)",
			station, startDate, endDate, trackRecords, err, expectedTrackRecords)
	}

	for i, expectedTrackRecord := range expectedTrackRecords {
		if trackRecords[i] != expectedTrackRecord {
			t.Errorf("GetTrackRecords(%q, %q, %q): got (%q, %v), expected (%q, nil)",
				station, startDate, endDate, trackRecords, err, expectedTrackRecords)
		}
	}
}

func TestDDBTrackRecordDAO_GetTrackRecordsFail(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDB{}, "testTable")

	var tests = []struct {
		inputStation   string
		inputStartDate time.Time
		inputEndDate   time.Time
	}{
		{"station-a", time.Now(), time.Now().AddDate(0, 0, -1)},
	}

	for _, test := range tests {
		trackRecords, err := trackRecordDAO.GetTrackRecords(test.inputStation, test.inputStartDate,
			test.inputEndDate)
		if err == nil {
			t.Errorf("GetTrackRecords(%q, %q, %q): got (%q, %v), expected (nil, error)",
				test.inputStation, test.inputStartDate, test.inputEndDate, trackRecords, err)
		}
	}
}
