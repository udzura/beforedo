package main

import (
	"flag"
	"fmt"
	"os"
)

var Version = "0.0.2"
var confPath string
var showVersion bool

func realMain() error {
	flag.StringVar(&confPath, "c", "Before.yaml", "Specify config path")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	tasks, err := ParseConfig(confPath)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		err = t.CheckAndRun()
		if err != nil {
			return err
		}
	}

	return nil // all OK
}

func main() {
	if err := realMain(); err != nil {
		panic(err)
	}
}
