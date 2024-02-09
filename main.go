package main

import (
	"fmt"
	"github.com/AndresBott/videoconv/app/cmd"
	"time"
)

var version = "DEV"
var commit = "SNAPSHOT"
var date string

func main() {
	if date == "" {
		now := time.Now()
		date = now.Format("2006-01-02")
	}
	cmd.Run(fmt.Sprintf("%s - %s (%s)", version, commit, date))
}
