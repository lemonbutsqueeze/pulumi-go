package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"io/ioutil"
)

var (
	CACHE_DIR_PATH   = "pulumi-go"
	CACHE_FILE_NAME  = "cached.json"
	CACHE_FULL_PATH  = CACHE_DIR_PATH + "/" + CACHE_FILE_NAME
	TIMESTAMP_FORMAT = "20060102150405"
)

type Cache struct {
	AwsCredentials AwsCredentials `json:"AwsCredentials"`
	CreateTime     string         `json:"CreateTime"`
}

func GetCache() (Cache, error) {
	var cache Cache

	path := CACHE_FULL_PATH
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return cache, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &cache)

	return cache, err
}

func WriteTemp(content string) {
	file, err := ioutil.TempFile(CACHE_DIR_PATH, CACHE_FILE_NAME)

	if err != nil {
		fmt.Println(err)
	}

	defer os.Remove(file.Name())

	if _, err := file.Write([]byte(content)); err != nil {
		fmt.Println(err)
	}
}

func CacheToJson(cache Cache) {
	bytes, err := json.Marshal(cache)

	if err != nil {
		fmt.Println(err)
		return
	}

	WriteTemp(string(bytes))
}

func GetCachedAwsCredentials(awsCredentials AwsCredentials) AwsCredentials {
	isCacheExist := true
	if _, err := os.Stat(CACHE_FULL_PATH); errors.Is(err, os.ErrNotExist) {
		isCacheExist = false
	}

	cache, _ := GetCache()
	cachedTime, _ := time.Parse(TIMESTAMP_FORMAT, cache.CreateTime)
	if !isCacheExist || cache.AwsCredentials.RoleArn != awsCredentials.RoleArn || time.Since(cachedTime).Minutes() > 60 {
		println("Cannot find existing cached creds")
		var newCache Cache
		newCache.AwsCredentials = awsCredentials
		newCache.CreateTime = time.Now().Format(TIMESTAMP_FORMAT)
		CacheToJson(newCache)
		return awsCredentials
	}

	println("Re-using cached AWS credentials")
	return cache.AwsCredentials
}
