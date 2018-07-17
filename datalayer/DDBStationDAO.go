package datalayer

import (
	"github.com/RadioCheckerApp/api/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

type DDBStationDAO struct {
	dynamoDB  *dynamodb.DynamoDB
	tableName string
}

func NewDDBStationDAO(dynamodb *dynamodb.DynamoDB, tableName string) *DDBStationDAO {
	return &DDBStationDAO{dynamodb, tableName}
}

func (dao *DDBStationDAO) GetAll() ([]model.Station, error) {
	var stations []model.Station

	input := &dynamodb.ScanInput{
		TableName: aws.String(dao.tableName),
	}

	err := dao.dynamoDB.ScanPages(input, func(page *dynamodb.ScanOutput, last bool) bool {
		var stats []model.Station
		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &stats)
		if err != nil {
			log.Printf("failed to unmarshal DynamoDB scan items: %v", err)
		}
		stations = append(stations, stats...)
		return true
	})
	if err != nil {
		return nil, err
	}

	return stations, nil
}
