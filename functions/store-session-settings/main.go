package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Option struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Icon   string `json:"icon"`
	Active bool   `json:"active"`
}

type UserObject struct {
	Id                  string   `json:"id"`
	Sort                string   `json:"sort"`
	Name                string   `json:"name"`
	DocType             string   `json:"docType"`
	Dni                 string   `json:"dni"`
	Gender              string   `json:"gender"`
	BirthDate           string   `json:"birthDate"`
	CountryOfBirth      string   `json:"countryOfBirth"`
	PersonalEmail       string   `json:"personalEmail"`
	MaritalStatus       string   `json:"maritalStatus"`
	PersonalPhone       string   `json:"personalPhone"`
	CountryOfResidence  string   `json:"countryOfResidence"`
	ResidenceDepartment string   `json:"residenceDepartment"`
	Address             string   `json:"address"`
	Area                string   `json:"area"`
	SubArea             string   `json:"subArea"`
	WorkerType          string   `json:"workerType"`
	Email               string   `json:"email"`
	CreationDate        string   `json:"creationDate"`
	EntryDate           string   `json:"entryDate"`
	Phone               string   `json:"phone"`
	Apps                []Option `json:"apps"`
	Menu                []Option `json:"menu"`
	Processes           []Option `json:"processes"`
	UserType            string   `json:"userType"`
	UserStatus          string   `json:"userStatus"`
	Role                string   `json:"role"`
	OfficeRole          string   `json:"officeRole"`
	Days                int      `json:"days"`
	HomeOffice          int      `json:"homeOffice"`
	Photo               string   `json:"photo"`
	Boss                string   `json:"boss,omitempty"`
	BossName            string   `json:"bossName,omitempty"`
	User                string   `json:"user"`
	Backup              string   `json:"backup"`
	BackupName          string   `json:"backupName"`
}

type Event struct {
	EmailId string `json:"emailId"`
	UserId  string `json:"userId"`
	Photo   string `json:"photo"`
}

func handler(ctx context.Context, event Event) (string, error) {
	TABLE_NAME := os.Getenv("TableName")
	SORT := "SETTINGS"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		return "", err
	}

	svc := dynamodb.New(sess)

	user, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(event.EmailId),
			},
			"sort": {
				S: aws.String(SORT),
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to Dynamodb Get Item, %v", err))
	}

	userData := UserObject{}

	err = dynamodbattribute.UnmarshalMap(user.Item, &userData)
	if err != nil {
		return "", err
	}

	in := UserObject{
		Id:                  event.UserId,
		Sort:                SORT,
		Name:                userData.Name,
		DocType:             userData.DocType,
		Dni:                 userData.Dni,
		Gender:              userData.Gender,
		BirthDate:           userData.BirthDate,
		CountryOfBirth:      userData.CountryOfBirth,
		PersonalEmail:       userData.PersonalEmail,
		MaritalStatus:       userData.MaritalStatus,
		PersonalPhone:       userData.PersonalPhone,
		CountryOfResidence:  userData.CountryOfResidence,
		ResidenceDepartment: userData.ResidenceDepartment,
		Address:             userData.Address,
		Area:                userData.Area,
		SubArea:             userData.SubArea,
		WorkerType:          userData.WorkerType,
		Email:               userData.Email,
		CreationDate:        userData.CreationDate,
		EntryDate:           userData.EntryDate,
		Phone:               userData.Phone,
		Apps:                userData.Apps,
		Menu:                userData.Menu,
		Processes:           userData.Processes,
		UserType:            userData.UserType,
		UserStatus:          "ACTIVE",
		Role:                userData.Role,
		OfficeRole:          userData.OfficeRole,
		Days:                userData.Days,
		HomeOffice:          userData.HomeOffice,
		Photo:               userData.Photo,
		Boss:                userData.Boss,
		BossName:            userData.BossName,
		User:                event.UserId,
		Backup:              event.UserId,
		BackupName:          userData.Name,
	}

	item, err := MarshalMap(in)
	if err != nil {
		return "", err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(TABLE_NAME),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return "", err
	}

	return "Success", nil
}

func main() {
	lambda.Start(handler)
}

func MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	av, err := getEncoder().Encode(in)
	if err != nil || av == nil || av.M == nil {
		return map[string]*dynamodb.AttributeValue{}, err
	}

	return av.M, nil
}

func getEncoder() *dynamodbattribute.Encoder {
	encoder := dynamodbattribute.NewEncoder()
	encoder.NullEmptyString = false
	return encoder
}
