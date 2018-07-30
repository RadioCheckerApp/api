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
	dynamoDB  DynamoDB
	tableName string
}

func NewDDBTrackRecordDAO(dynamodb DynamoDB, tableName string) *DDBTrackRecordDAO {
	return &DDBTrackRecordDAO{dynamodb, tableName}
}

func (dao *DDBTrackRecordDAO) GetTrackRecords(station string, startDate, endDate time.Time) ([]model.TrackRecord, error) {
	if startDate.After(endDate) {
		return nil, errors.New("startDate must be before endDate")
	}

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(dao.tableName),
		KeyConditionExpression: aws.String(
			"#sid = :stationId AND airtime BETWEEN :lowerBound AND :upperBound"),
		ExpressionAttributeNames: map[string]*string{
			"#sid": aws.String("stationId"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":stationId":  {S: aws.String(station)},
			":lowerBound": {N: aws.String(strconv.FormatInt(startDate.Unix(), 10))},
			":upperBound": {N: aws.String(strconv.FormatInt(endDate.Unix(), 10))},
		},
	}

	output, err := dao.dynamoDB.Query(queryInput)
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
