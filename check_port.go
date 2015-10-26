package main

import (
	"fmt"
	"net"
)

func checkPort(t *Task) (bool, error) {
	addr := fmt.Sprintf("127.0.0.1:%d", t.Port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// Is OK to ignore
		return true, nil
	} else {
		conn.Close()
		return false, nil
	}
}
