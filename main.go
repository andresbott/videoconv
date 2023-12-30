package main

import (
	"fmt"
	"github.com/AndresBott/videoconv/app/cmd"
	"time"
)

var version = "DEV"
var commit = "SNAPSHOT"
var date = time.Now()

func main() {
	cmd.Run(fmt.Sprintf("%s - %s (%s)", version, commit, date.Format("2006-01-02")))
}
