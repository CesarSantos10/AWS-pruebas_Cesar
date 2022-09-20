package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type deps struct {
}

type CognitoClient interface {
	SignUp(email string, password string) (string, error)
	AdminCreateUser(email string) (string, error)
	SignIn(email string, password string) (string, error)
	AdminDisableUser(username string) (string, error)
	AdminEnableUser(username string) (string, error)
}

type awsCognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	appClientId   string
	userPoolId    string
}

type Event struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Case     int    `json:"case"`
}

func (d *deps) handler(ctx context.Context, event Event) (string, error) {
	// CONECTAR SESSION CON AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*aws.String("us-east-1"))},
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect session, %v", err))
	}
	// INICIAR SESSION EN COGNITO
	svc := cognito.New(sess)

	client := awsCognitoClient{
		cognitoClient: svc,
		appClientId:   "1brn5dsq5sbom0ba9bckeqlmve",
		userPoolId:    "us-east-1_gDzPxPab7",
	}
	fmt.Printf("Email :%s Password: %s Name: %s, UserName: %s \n", event.Email, event.Password, event.Name, event.Username)
	fmt.Println("cliente: ", client)

	switch event.Case {
	case 0: // SignUp
		client.SignUp(event.Email, event.Password)
	case 1: // AdminCreateUser
		client.AdminCreateUser(event.Email, event.Name)
	case 2: // SignIn
		client.SignIn(event.Email, event.Password)
	case 3: // AdminDisableUser
		client.AdminDisableUser(event.Username)
	case 4: // AdminDisableUser
		client.AdminEnableUser(event.Username)
	}

	fmt.Print(client)
	return "", nil
}

func main() {
	d := deps{}
	lambda.Start(d.handler)
}

func (ctx *awsCognitoClient) SignUp(email string, password string) (string, error) {

	user := &cognito.SignUpInput{
		ClientId: aws.String(ctx.appClientId),
		Username: aws.String(email),
		Password: aws.String(password),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
	fmt.Println("USER:  ", user)

	result, err := ctx.cognitoClient.SignUp(user)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return result.String(), nil
}

func (ctx *awsCognitoClient) AdminCreateUser(email string, name string) (string, error) {

	user := &cognito.AdminCreateUserInput{
		UserPoolId: aws.String(ctx.userPoolId),
		Username:   aws.String(email),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("custom:name"),
				Value: aws.String(name),
			},
		},
	}
	fmt.Println("USER: ", user)
	fmt.Println("Name: ", name)

	result, err := ctx.cognitoClient.AdminCreateUser(user)
	if err != nil {
		fmt.Println("Error :", err)
		return "", err
	}
	return result.String(), nil
}

func (ctx *awsCognitoClient) SignIn(email string, password string) (string, error) {
	initiateAuthInput := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: aws.StringMap(map[string]string{
			"USERNAME": email,
			"PASSWORD": password,
		}),
		ClientId: aws.String(ctx.appClientId),
	}

	result, err := ctx.cognitoClient.InitiateAuth(initiateAuthInput)

	if err != nil {
		fmt.Println("Error  : InitiateAuth", err)
		return "", err
	}

	return result.String(), nil
}

func (ctx *awsCognitoClient) AdminDisableUser(username string) (string, error) {

	adminDisableUserInput := &cognito.AdminDisableUserInput{
		UserPoolId: aws.String(ctx.userPoolId),
		Username:   aws.String(username),
	}

	result, _ := ctx.cognitoClient.AdminDisableUser(adminDisableUserInput)

	return result.String(), nil
}

func (ctx *awsCognitoClient) AdminEnableUser(username string) (string, error) {

	adminEnableUserInput := &cognito.AdminEnableUserInput{
		UserPoolId: aws.String(ctx.userPoolId),
		Username:   aws.String(username),
	}

	result, _ := ctx.cognitoClient.AdminEnableUser(adminEnableUserInput)

	return result.String(), nil
}
