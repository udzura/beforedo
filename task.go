package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

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

func (t *Task) CheckAndRun() error {
	if t.Front {
		// return t.keepRunInFront()
		return fmt.Errorf("Not yet implemented")
	} else {
		toRun, err := t.toRun()
		if err != nil {
			return err
		}
		if toRun {
			return t.run()
		}
		fmt.Printf("Skipping Run: %s\n", t.Command)
		return nil
	}
}

func (t *Task) toRun() (bool, error) {
	switch {
	case t.Port != 0:
		return checkPort(t)
	case t.File != "":
		return checkFile(t)
	case t.Always:
		return true, nil
	default:
		return false, fmt.Errorf("Invalid or unsupported config: %v", t)
	}
}

func (t *Task) run() error {
	fmt.Printf("Run: %s\n", t.Command)

	cmd := exec.Command("sh", "-c", t.Command)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	scanOut := bufio.NewScanner(stdout)
	scanErr := bufio.NewScanner(stderr)
	go func() {
		for scanOut.Scan() {
			fmt.Printf("---> %s\n", scanOut.Text())
		}
		if err := scanOut.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Reading standard input:", err)
		}
	}()
	go func() {
		for scanErr.Scan() {
			fmt.Printf("---! %s\n", scanErr.Text())
		}
		if err := scanErr.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Reading standard input:", err)
		}
	}()

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	fmt.Printf("OK: %s\n", t.Command)

	return nil
}

type Tasks []*Task

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
