package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

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
		return t.keepRunInFront()
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
	case t.SuccessCommand != "":
		return checkSuccessCommand(t)
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
	go startStreamingOutputs(stdout, stderr)

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	var exitCode int
	if err == nil {
		exitCode = 0
	} else {
		if exit, ok := err.(*exec.ExitError); ok {
			if s, ok := exit.Sys().(syscall.WaitStatus); ok {
				exitCode = s.ExitStatus()
			}
		}
		if exitCode != 0 {
			return fmt.Errorf("Command exited with code %d.\nOriginal Error: %s", exitCode, err.Error())
		} else {
			return err
		}
	}

	fmt.Printf("OK: %s\n", t.Command)

	return nil
}

func (t *Task) keepRunInFront() error {
	fmt.Printf("Start blocking process: %s\n", t.Command)

	cmd := exec.Command("sh", "-c", t.Command)
	sin, _ := cmd.StdinPipe()
	sout, _ := cmd.StdoutPipe()
	serr, _ := cmd.StderrPipe()
	go startStreamingOutputs(sout, serr)

	err := cmd.Start()
	if err != nil {
		return err
	}
	err = sin.Close()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("!! Interrupted!\n")
	}()

	err = cmd.Wait()
	var exitCode int
	if err == nil {
		exitCode = 0
		fmt.Printf("Command exited with code %d.\n", exitCode)
	} else {
		exit, ok := err.(*exec.ExitError)
		switch {
		case ok:
			if s, ok := exit.Sys().(syscall.WaitStatus); ok {
				exitCode = s.ExitStatus()
			}
		case strings.Contains(err.Error(), "interrupt"):
			fmt.Printf("Command is interrupted.\nOriginal Error: %s\n", err.Error())
		default:
			return err
		}

		if exitCode > 0 {
			fmt.Printf("Command exited with code %d.\nOriginal Error: %s\n", exitCode, err.Error())
		}
	}

	fmt.Printf("Finished!: %s\n", t.Command)

	// A front: true task would not continue following tasks.
	os.Exit(0)
	return nil
}

func startStreamingOutputs(stdout, stderr io.Reader) {
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
