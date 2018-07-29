package datalayer

import (
	"github.com/RadioCheckerApp/api/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type DDBStationDAO struct {
	dynamoDB  DynamoDB
	tableName string
}

func NewDDBStationDAO(dynamodb DynamoDB, tableName string) *DDBStationDAO {
	return &DDBStationDAO{dynamodb, tableName}
}

func (dao *DDBStationDAO) GetAll() ([]model.Station, error) {
	stations := make([]model.Station, 0)

	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(dao.tableName),
	}

	err := dao.dynamoDB.ScanPages(scanInput, func(page *dynamodb.ScanOutput, last bool) bool {
		var stats []model.Station
		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &stats)
		if err != nil {
			log.Printf("failed to unmarshal DynamoDB scan items: %v", err)
		}
		stations = append(stations, stats...)
		return true
	})

	return stations, err
}
