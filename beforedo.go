package main

import "flag"

var confPath string

func realMain() error {
	flag.StringVar(&confPath, "c", "Before.yaml", "Specify config path")
	flag.Parse()

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
