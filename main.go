/*
Copyright © 2024 NAME HERE ra.skalbania@gmail.com
*/
package main

import (
	"os"

	"github.com/radekska/outbox/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.App.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
