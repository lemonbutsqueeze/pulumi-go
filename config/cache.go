package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

const (
	CACHE_DIR_PATH   = "pulumi-go"
	CACHE_FILE_NAME  = "cached.json"
	TIMESTAMP_FORMAT = "20060102150405"
)

var (
	CACHE_FULL_PATH = filepath.Join(CACHE_DIR_PATH, CACHE_FILE_NAME)
)

type Cache struct {
	IamProfile     IamProfile     `json:"IamProfile"`
	CreateTime     string         `json:"CreateTime"`
	AwsCredentials AwsCredentials `json:"AwsCredentials"`
}

func IsFileExist(path string) bool {
	log.Debug("Checking if file '" + path + "' exist ..")
	if _, err := os.Stat(CACHE_FULL_PATH); errors.Is(err, os.ErrNotExist) {
		log.Debug("File does not exist")
		return false
	}

	log.Debug("File exist")
	return true
}

func GetCache() (*Cache, error) {
	var cache *Cache = &Cache{}

	path := CACHE_FULL_PATH
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &cache)

	return cache, err
}

func WriteTempFile(content string) error {

	log.Debug("Creating file at " + CACHE_FULL_PATH)
	file, err := os.OpenFile(CACHE_FULL_PATH, os.O_RDONLY|os.O_CREATE, os.FileMode(os.O_CREATE))
	// file, err := os.Create(CACHE_FULL_PATH)
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	file.Sync()
	defer file.Close()

	return err
}

func CreateCache(iamProfile IamProfile, awsCredentials AwsCredentials) error {
	var newCache Cache
	newCache.IamProfile = iamProfile
	newCache.CreateTime = time.Now().Format(TIMESTAMP_FORMAT)
	newCache.AwsCredentials = awsCredentials

	bytes, err := json.Marshal(newCache)

	if err != nil {
		return err
	}

	return WriteTempFile(string(bytes))
}

func GetCachedAwsCredentials(iamProfile *IamProfile) (*AwsCredentials, error) {

	if !IsFileExist(CACHE_FULL_PATH) {
		return nil, nil
	}

	cache, err := GetCache()
	if err != nil {
		log.Warn("Failed while fetching cache", err)
	}

	cachedTime, _ := time.Parse(TIMESTAMP_FORMAT, cache.CreateTime)
	if cache.IamProfile.RoleArn != iamProfile.RoleArn || time.Since(cachedTime).Minutes() > 60 {
		return nil, err
	}

	log.Info("Re-using cached AWS credentials")
	return &cache.AwsCredentials, err
}
