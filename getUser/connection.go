package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Table from DynamoDB
const tableName = "rcs-serverless-users"

// Controller dedicated to fetch the user, should receive the id as a parameter.
func FetchUser(id string) (u User, err error) {
	// Prepare DynamoDB attribute value
	av := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(id),
		},
	}

	// Stablish the connection with DynamoDB
	dynamoSession := createDynamoSession()

	ctx := context.Background()

	// Get the user from the database
	fetchedUser, err := dynamoSession.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       av,
	})

	if err != nil {
		fmt.Println("Failed to get from dynamo: ", err.Error())
		return
	}

	// Save User record on User struct
	err = dynamodbattribute.UnmarshalMap(fetchedUser.Item, &u)
	if err != nil {
		fmt.Println("Failed unmarshaling record: ", err.Error())
		return
	}

	return
}

// Stablish the connection to DynamoDB
func createDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	return dynamodb.New(sess)
}
