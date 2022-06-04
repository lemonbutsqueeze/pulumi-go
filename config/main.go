package main

import (
	"fmt"
	"os"
)

type Config struct {
	AwsCredentials AwsCredentials `json:"AwsCredentials"`
}

type AwsCredentials struct {
	RoleArn      string `json:"RoleArn"`
	MfaSerialArn string `json:"MfaSerialArn"`
}

func main() {
	jsonFile, err := os.Open("../config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
}
