package datalayer

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoDB interface {
	ScanPages(input *dynamodb.ScanInput, fn func(*dynamodb.ScanOutput, bool) bool) error
}
