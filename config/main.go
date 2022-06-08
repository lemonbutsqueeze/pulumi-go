package main

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	err := AuthenticateAws()
	if err != nil {
		log.WithError(err).Error("This is an error\n", string(debug.Stack()))
		panic(err)
	}
}
