package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	AwsCredentials AwsCredentials `json:"AwsCredentials"`
}

type AwsCredentials struct {
	RoleArn      string `json:"RoleArn"`
	MfaSerialArn string `json:"MfaSerialArn"`
}

func GetConfig(path string) (Config, error) {
	var config Config

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return config, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	return config, err
}
