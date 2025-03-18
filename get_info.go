package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Manager struct {
	processes map[string][]string
	tcp       map[string]Process
}

type Process struct {
	pid          string
	port         string
	process_name string
	inode        []byte
}

// todo get process_name

const PROC_NET_PATH string = "/proc/net/"
const PROC_PATH string = "/proc/"
const NEW_LINE_MASK byte = 10

const (
	QUOTE_PORT  int = 2
	QUOTE_FINAL int = 6
)

func (m *Manager) cacheProc() {
	m.tcp = make(map[string]Process)
	files, err := os.ReadDir(PROC_PATH)

	if err != nil {
		log.Fatalln("ERROR: Reading dir", err)
	}

	for _, file := range files {
		n, ok := strconv.Atoi(file.Name())

		if ok != nil {
			continue
		}
		path := fmt.Sprintf("%s%d/", PROC_PATH, n)
		fds, err := os.ReadDir(path + "fd/")

		if err != nil {
			continue
		}
		for _, fd := range fds {
			linkPath := fmt.Sprintf("%sfd/%s", path, fd.Name())
			result, err := os.Readlink(linkPath)

			if err != nil {
				continue
			}

			if strings.Contains(result, "socket") {
				cmdline, _ := os.ReadFile(path + "cmdline")
				m.processes[result] = []string{fmt.Sprintf("%d", n), string(cmdline)}
			}
		}
	}
}

func (m *Manager) ParseTCP() {
	file, err := os.ReadFile(PROC_NET_PATH + "tcp6")

	if err != nil {
		log.Fatalf("ERROR: Could not read %s because %s", PROC_NET_PATH, err)
	}

	for idx, line := range bytes.Split(file, []byte("\n")) {

		if idx == 0 {
			continue
		}
		chunks := bytes.Split(line, []byte(" "))

		if len(chunks) == 1 {
			break
		}
		p := Process{}
		p.inode = chunks[len(chunks)-8]
		search := fmt.Sprintf("socket:[%s]", string(p.inode))
		result := m.processes[search]
		if len(result) == 0 {
			continue
		}

		// convert byte to string then to number then to string again
		portLiteral := string(chunks[4][len(chunks[4])-4 : len(chunks[4])])
		port, _ := strconv.ParseInt(string(portLiteral), 16, 64)
		p.port = fmt.Sprintf("%d", port)

		p.pid, p.process_name = result[0], result[1]
		m.tcp[search] = p
	}

}
