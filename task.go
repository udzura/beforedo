package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Task struct {
	Command        string `yaml:"task"`
	Port           int
	File           string
	SuccessCommand string `yaml:"success"`
	Always         bool
	Front          bool
}

type Tasks []Task

func ParseConfig(path string) (Tasks, error) {
	if path == "" {
		path = "./Before.yaml"
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tasks Tasks
	if err = yaml.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
