package main

import (
	"pulumi-go/config"
	"pulumi-go/network"
)

func main() {
	config.GetConfig()
	network.Execute()
}
