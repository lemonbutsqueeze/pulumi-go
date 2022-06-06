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

func GetConfig() (Config, error) {
	var config Config

	path, _ := os.Getwd()
	jsonFile, err := os.Open(path + "/config/config.json")
	if err != nil {
		fmt.Println(err)
		return config, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	return config, err
}
