package main

import (
	"fmt"
	"os"
)

func checkFile(t *Task) (bool, error) {
	s, err := os.Stat(t.File)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		} else {
			return false, err
		}
	}
	if s.IsDir() {
		return false, fmt.Errorf(
			"The file %s seems to be a directory, not a normal file. Maybe a mistake?", t.File)
	}
	return false, nil
}
