package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	CONFIG_FILE_PATH = "../config/config.json"
)

type Config struct {
	AwsCredentials AwsCredentials `json:"AwsCredentials"`
}

type AwsCredentials struct {
	RoleArn      string `json:"RoleArn"`
	MfaSerialArn string `json:"MfaSerialArn"`
}

func GetConfig(path string) (Config, error) {
	println("Grabbing config file from " + path)
	var config Config

	jsonFile, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	return config, err
}

func AssumeRoleWithMfa(roleArn string, mfaSerialArn string) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Failed to load configuration, %v", err)
	}

	stsClient := sts.NewFromConfig(cfg)

	provider := stscreds.NewAssumeRoleProvider(stsClient, roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.SerialNumber = aws.String(mfaSerialArn)
		o.TokenProvider = stscreds.StdinTokenProvider
	})

	cfg.Credentials = aws.NewCredentialsCache(provider)
	creds, err := cfg.Credentials.Retrieve(context.Background())

	if err != nil {
		log.Fatalf("Failed to retrieve role configuration, %v", err)
	}

	os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)

	println("Succesfully configured AWS credentials using role ARN: ", roleArn)
}

func IsAuthenticated() bool {
	_, exist := os.LookupEnv("AWS_SESSION_TOKEN")
	return exist
}

func AuthenticateAws() {
	cfg, _ := GetConfig(CONFIG_FILE_PATH)
	creds := GetCachedAwsCredentials(cfg.AwsCredentials)

	println("Authenticating with AWS ..")
	AssumeRoleWithMfa(creds.RoleArn, creds.MfaSerialArn)
}
