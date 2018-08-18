package datalayer

import (
	"errors"
	"github.com/RadioCheckerApp/api/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
	"time"
)

type DDBTrackRecordDAO struct {
	dynamoDB       DynamoDB
	tableName      string
	gsiTypeAirtime string
}

func NewDDBTrackRecordDAO(dynamodb DynamoDB, tableName, gsiTypeAirtime string) *DDBTrackRecordDAO {
	return &DDBTrackRecordDAO{dynamodb, tableName, gsiTypeAirtime}
}

func (dao *DDBTrackRecordDAO) GetTrackRecords(startDate, endDate time.Time) ([]model.TrackRecord, error) {
	if err := valiDate(startDate, endDate); err != nil {
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(dao.tableName),
		IndexName: aws.String(dao.gsiTypeAirtime),
		KeyConditionExpression: aws.String(
			"#t = :type AND airtime BETWEEN :lowerBound AND :upperBound"),
		ExpressionAttributeNames: map[string]*string{
			"#t": aws.String("type"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":type":       {S: aws.String("track")},
			":lowerBound": {N: aws.String(strconv.FormatInt(startDate.Unix(), 10))},
			":upperBound": {N: aws.String(strconv.FormatInt(endDate.Unix(), 10))},
		},
	}

	return dao.executeQuery(queryInput)
}

func (dao *DDBTrackRecordDAO) GetTrackRecordsByStation(station string, startDate,
	endDate time.Time) ([]model.TrackRecord, error) {
	if err := valiDate(startDate, endDate); err != nil {
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(dao.tableName),
		KeyConditionExpression: aws.String(
			"#sid = :stationId AND airtime BETWEEN :lowerBound AND :upperBound"),
		FilterExpression: aws.String("#t = :type"),
		ExpressionAttributeNames: map[string]*string{
			"#sid": aws.String("stationId"),
			"#t":   aws.String("type"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":stationId":  {S: aws.String(station)},
			":lowerBound": {N: aws.String(strconv.FormatInt(startDate.Unix(), 10))},
			":upperBound": {N: aws.String(strconv.FormatInt(endDate.Unix(), 10))},
			":type":       {S: aws.String("track")},
		},
	}

	return dao.executeQuery(queryInput)
}

func (dao *DDBTrackRecordDAO) executeQuery(input *dynamodb.QueryInput) ([]model.TrackRecord,
	error) {
	output, err := dao.dynamoDB.Query(input)
	if err != nil {
		return nil, err
	}

	trackRecords := make([]model.TrackRecord, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(output.Items, &trackRecords)
	if err != nil {
		return nil, err
	}

	return trackRecords, nil
}

func valiDate(startDate, endDate time.Time) error {
	if startDate.After(endDate) {
		return errors.New("startDate must be before endDate")
	}
	return nil
}
