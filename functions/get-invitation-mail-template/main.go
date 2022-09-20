package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Event struct {
	UserType string `json:"userType"`
}

func handler(ctx context.Context, event Event) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	svc := s3.New(sess)

	userType := event.UserType
	var htmlPage string

	switch userType {
	case "Usuario Protecta":
		htmlPage = "protectaUser"
	case "Administrador Protecta":
		htmlPage = "protectaUser"
	case "Usuario Externo":
		htmlPage = "externalUser"
	}

	rawObject, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("BucketName")),
			Key:    aws.String(fmt.Sprintf("%s.html", htmlPage)),
		})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(rawObject.Body)
	result := buf.String()
	return result, nil

}

func main() {
	lambda.Start(handler)
}
