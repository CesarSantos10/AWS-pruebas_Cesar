package main

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Age string `json:"age"`
}

func handler(event Event) (int, error) {

	intAge, _ := strconv.Atoi(event.Age)

	fmt.Println(intAge)

	return intAge + 10, nil
}

func main() {
	lambda.Start(handler)
}
