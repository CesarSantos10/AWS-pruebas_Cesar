package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type StudentObject struct {
	Id   string `json:"id"`
	Sk   string `json:"sk"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Event struct {
	Student StudentObject `json:"student"`
}

func handler(ctx context.Context, event Event) (StudentObject, error) {
	TABLE_NAME := os.Getenv("TableName")

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		fmt.Printf("1")
		return StudentObject{}, err
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	fmt.Printf("2")
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(event.Student.Name),
			},
			":age": {
				N: aws.String(strconv.Itoa(event.Student.Age)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#name": aws.String("name"),
		},
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(event.Student.Id),
			},
			"sk": {
				S: aws.String(event.Student.Sk),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set #name = :name, age = :age"),
	}
	fmt.Printf("3")
	result, err := svc.UpdateItem(input)

	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}
	fmt.Printf("4")
	studentObject := StudentObject{}

	err1 := dynamodbattribute.UnmarshalMap(result.Attributes, &studentObject)
	if err1 != nil {
		log.Fatalf("Got error calling UnmarshalMap: %s", err)
	}
	fmt.Printf("5")
	return studentObject, nil
}

func main() {
	lambda.Start(handler)
}
