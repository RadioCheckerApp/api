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
		return nil, errors.New("KeyConditionExpression must not be nil")
	}

	if len(strings.Split(*input.KeyConditionExpression, " ")) != 9 {
		return nil, errors.New("KeyConditionExpression must exist of 9 words")
	}

	if input.ExpressionAttributeNames == nil {
		return nil, errors.New("ExpressionAttributeNames must not be nil")
	}

	if name, ok := input.ExpressionAttributeNames["#sid"]; input.IndexName == nil && (!ok || *name != "stationId") {
		return nil, errors.New("ExpressionAttributeNames must contain mapping `#sid`: `stationId`")
	}

	if name, ok := input.ExpressionAttributeNames["#t"]; !ok || *name != "type" {
		return nil, errors.New("ExpressionAttributeNames must contain mapping `#t`: `type`")
	}

	if input.ExpressionAttributeValues == nil {
		return nil, errors.New("ExpressionAttributeValues must not be nil")
	}

	if _, ok := input.ExpressionAttributeValues[":stationId"]; input.IndexName == nil && !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:stationId`")
	}

	if _, ok := input.ExpressionAttributeValues[":lowerBound"]; !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:lowerBound`")
	}

	if _, ok := input.ExpressionAttributeValues[":upperBound"]; !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:upperBound`")
	}

	if value, ok := input.ExpressionAttributeValues[":type"]; !ok || *value.S != "track" {
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

func (ddb MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	if input == nil {
		return nil, errors.New("input must not be nil")
	}

	if input.TableName == nil {
		return nil, errors.New("TableName must not be nil")
	}

	if input.Item == nil || len(input.Item) != 5 {
		return nil, errors.New("Item must contain 5 mappings")
	}

	if input.ConditionExpression == nil ||
		*input.ConditionExpression != "attribute_not_exists(stationId)" {
		return nil, errors.New("ConditionExpression must be `attribute_not_exists(stationId)`")
	}
	return nil, nil
}

type MockDynamoDBLimitedQuery struct{}

func (ddb MockDynamoDBLimitedQuery) ScanPages(input *dynamodb.ScanInput,
	fn func(*dynamodb.ScanOutput, bool) bool) error {
	return nil
}

func (ddb MockDynamoDBLimitedQuery) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput,
	error) {
	if input == nil {
		return nil, errors.New("input must not be nil")
	}

	if input.TableName == nil {
		return nil, errors.New("TableName must not be nil")
	}

	if input.KeyConditionExpression == nil {
		return nil, errors.New("KeyConditionExpression must not be nil")
	}

	if *input.KeyConditionExpression != "stationId = :stationId" {
		return nil, errors.New("KeyConditionExpression must be `stationId = :stationId`")
	}

	if input.ExpressionAttributeNames != nil {
		return nil, errors.New("ExpressionAttributeNames must be nil")
	}

	if input.ExpressionAttributeValues == nil {
		return nil, errors.New("ExpressionAttributeValues must not be nil")
	}

	if _, ok := input.ExpressionAttributeValues[":stationId"]; input.IndexName == nil && !ok {
		return nil, errors.New("ExpressionAttributeValues must contain a key `:stationId`")
	}

	if *input.ExpressionAttributeValues[":stationId"].S == "notracksstation" {
		return &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{},
		}, nil
	}

	if *input.ExpressionAttributeValues[":stationId"].S == "error" {
		return nil, errors.New("database error")
	}

	stationId := "station-a"
	airtime := "1234567890"
	trackType := "track"
	artist := "rhcp"
	title := "californication"

	output := &dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"stationId": &dynamodb.AttributeValue{S: &stationId},
				"airtime":   &dynamodb.AttributeValue{N: &airtime},
				"type":      &dynamodb.AttributeValue{S: &trackType},
				"artist":    &dynamodb.AttributeValue{S: &artist},
				"title":     &dynamodb.AttributeValue{S: &title},
			},
		},
	}

	return output, nil
}

func (ddb MockDynamoDBLimitedQuery) PutItem(input *dynamodb.PutItemInput) (*dynamodb.
	PutItemOutput, error) {
	return nil, nil
}

func TestDDBTrackRecordDAO_GetTrackRecordsSuccess(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(
		MockDynamoDB{},
		"testTable",
		"trackrecords-table-dev-gsi-type-airtime")
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now()

	expectedTrackRecords := []model.TrackRecord{
		{"station-a", 1532897851, "track", model.Track{"Mø", "Final Song"}},
		{"station-b", 1532897892, "track", model.Track{"Jack Ü ft. Skrillex & Diplo",
			"Where Are Ü Now"}},
	}

	trackRecords, err := trackRecordDAO.GetTrackRecords(startDate, endDate)

	if err != nil {
		t.Errorf("(%q) GetTrackRecords(%q, %q): got (%q, %v), expected (%q, nil)",
			trackRecordDAO, startDate, endDate, trackRecords, err, expectedTrackRecords)
		return
	}

	for i, expectedTrackRecord := range expectedTrackRecords {
		if trackRecords[i] != expectedTrackRecord {
			t.Errorf("(%q) GetTrackRecords(%q, %q): got (%q, %v), expected (%q, nil)",
				trackRecordDAO, startDate, endDate, trackRecords, err, expectedTrackRecords)
		}
	}
}

func TestDDBTrackRecordDAO_GetTrackRecordsFail(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(
		MockDynamoDB{},
		"testTable",
		"trackrecords-table-dev-gsi-type-airtime")

	var tests = []struct {
		inputStartDate time.Time
		inputEndDate   time.Time
	}{
		{time.Now(), time.Now().AddDate(0, 0, -1)},
	}

	for _, test := range tests {
		trackRecords, err := trackRecordDAO.GetTrackRecords(test.inputStartDate, test.inputEndDate)
		if err == nil {
			t.Errorf("(%q) GetTrackRecords(%q, %q): got (%q, %v), expected (nil, error)",
				trackRecordDAO, test.inputStartDate, test.inputEndDate, trackRecords, err)
		}
	}
}

func TestDDBTrackRecordDAO_GetTrackRecordsByStationSuccess(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDB{}, "testTable", "gsi")
	station := "ignoredDueToMock"
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now()

	expectedTrackRecords := []model.TrackRecord{
		{"station-a", 1532897851, "track", model.Track{"Mø", "Final Song"}},
		{"station-b", 1532897892, "track", model.Track{"Jack Ü ft. Skrillex & Diplo",
			"Where Are Ü Now"}},
	}

	trackRecords, err := trackRecordDAO.GetTrackRecordsByStation(station, startDate, endDate)

	if err != nil {
		t.Errorf("(%q) GetTrackRecordsByStation(%q, %q, %q): got (%q, %v), expected (%q, nil)",
			trackRecordDAO, station, startDate, endDate, trackRecords, err, expectedTrackRecords)
		return
	}

	for i, expectedTrackRecord := range expectedTrackRecords {
		if trackRecords[i] != expectedTrackRecord {
			t.Errorf("(%q) GetTrackRecordsByStation(%q, %q, %q): got (%q, %v), expected (%q, nil)",
				trackRecordDAO, station, startDate, endDate, trackRecords, err, expectedTrackRecords)
		}
	}
}

func TestDDBTrackRecordDAO_GetTrackRecordsByStationFail(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDB{}, "testTable", "gsi")

	var tests = []struct {
		inputStation   string
		inputStartDate time.Time
		inputEndDate   time.Time
	}{
		{"station-a", time.Now(), time.Now().AddDate(0, 0, -1)},
	}

	for _, test := range tests {
		trackRecords, err := trackRecordDAO.GetTrackRecordsByStation(test.inputStation,
			test.inputStartDate, test.inputEndDate)
		if err == nil {
			t.Errorf("(%q) GetTrackRecordsByStation(%q, %q, %q): got (%q, %v), expected (nil, error)",
				trackRecordDAO, test.inputStation, test.inputStartDate, test.inputEndDate, trackRecords, err)
		}
	}
}

func TestDDBTrackRecordDAO_GetMostRecentTrackRecordByStationSuccess(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDBLimitedQuery{}, "testTable", "gsi")
	station := "station-a"

	expectedTrackRecord := model.TrackRecord{
		"station-a",
		1234567890,
		"track",
		model.Track{"rhcp", "californication"},
	}

	trackRecord, err := trackRecordDAO.GetMostRecentTrackRecordByStation(station)

	if err != nil {
		t.Errorf("(%q) GetMostRecentTrackRecordByStationSuccess(%q): got (%q, %v), expected (%q, nil)",
			trackRecordDAO, station, trackRecord, err, expectedTrackRecord)
		return
	}

	if trackRecord != expectedTrackRecord {
		t.Errorf("(%q) GetMostRecentTrackRecordByStationSuccess(%q): got (%q, %v), expected (%q, nil)",
			trackRecordDAO, station, trackRecord, err, expectedTrackRecord)
	}
}

func TestDDBTrackRecordDAO_GetMostRecentTrackRecordByStationFail(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDBLimitedQuery{}, "testTable", "gsi")

	var tests = []string{"notracksstation", "error"}

	for _, test := range tests {
		trackRecord, err := trackRecordDAO.GetMostRecentTrackRecordByStation(test)
		if err == nil {
			t.Errorf("(%q) GetMostRecentTrackRecordByStation(%q): got (%q, %v), expected (nil, error)",
				trackRecordDAO, test, trackRecord, err)
		}
	}
}

func TestDDBTrackRecordDAO_CreateTrackRecord(t *testing.T) {
	trackRecordDAO := NewDDBTrackRecordDAO(MockDynamoDB{}, "testTable", "gsi")

	var tests = []model.TrackRecord{
		{"station-a", time.Now().Unix(), "track", model.Track{"RHCP", "Californication"}},
	}

	for _, testRecord := range tests {
		err := trackRecordDAO.CreateTrackRecord(testRecord)
		if err != nil {
			t.Errorf("(%q) CreateTrackRecord(%q): got err (%v), expected err: false",
				trackRecordDAO, testRecord, err)
		}
	}
}
