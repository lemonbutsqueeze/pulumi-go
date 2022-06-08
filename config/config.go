package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var (
	CONFIG_FILE_PATH = "../config/config.json"
)

type Config struct {
	IamProfile IamProfile `json:"IamProfile"`
}

type IamProfile struct {
	RoleArn      string `json:"RoleArn"`
	MfaSerialArn string `json:"MfaSerialArn"`
}

type AwsCredentials struct {
	AccessKey    string `json:"AccessKey"`
	SecretKey    string `json:"SecretKey"`
	SessionToken string `json:"SessionToken"`
}

func GetConfig(path string) (Config, error) {
	log.Info("Grabbing config file from " + path)
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

func AssumeRoleWithMfa(roleArn string, mfaSerialArn string) error {
	log.Info("Authenticating AWS with MFA using role '" + roleArn + "' ..")

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Error("Failed to load configuration")
		return err
	}

	stsClient := sts.NewFromConfig(cfg)

	provider := stscreds.NewAssumeRoleProvider(stsClient, roleArn, func(o *stscreds.AssumeRoleOptions) {
		o.SerialNumber = aws.String(mfaSerialArn)
		o.TokenProvider = stscreds.StdinTokenProvider
	})

	cfg.Credentials = aws.NewCredentialsCache(provider)
	creds, err := cfg.Credentials.Retrieve(context.Background())

	if err != nil {
		log.Error("Failed to load configuration")
		return err
	}

	SetAwsCredentials(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken)

	log.Info("Successfully configured AWS credentials")
	return err
}

func SetAwsCredentials(accessKey string, secretKey string, sessionToken string) {
	log.Debug("Setting AWS credentials using env vars ..")
	os.Setenv("AWS_ACCESS_KEY_ID", accessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretKey)
	os.Setenv("AWS_SESSION_TOKEN", sessionToken)
}

func GetCurrentAwsCredentials() *AwsCredentials {
	var awsCredentials *AwsCredentials = &AwsCredentials{
		AccessKey:    os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SessionToken: os.Getenv("AWS_SESSION_TOKEN"),
	}

	return awsCredentials
}

func AuthenticateAws() error {
	cfg, err := GetConfig(CONFIG_FILE_PATH)
	if err != nil {
		return err
	}

	cachedCreds, err := GetCachedAwsCredentials(&cfg.IamProfile)
	if err != nil {
		return err
	}

	if cachedCreds == nil {
		err := AssumeRoleWithMfa(cfg.IamProfile.RoleArn, cfg.IamProfile.MfaSerialArn)
		if err != nil {
			return err
		}

		CreateCache(cfg.IamProfile, *GetCurrentAwsCredentials())
	} else {
		SetAwsCredentials(cachedCreds.AccessKey, cachedCreds.SecretKey, cachedCreds.SessionToken)
	}

	return err
}
