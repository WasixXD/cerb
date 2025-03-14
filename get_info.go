package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Process struct {
	pid          int
	port         []byte
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

// this is dumb but lets go this way
// probably the best thing should be to cache the /proc and search by the inode
// TODO: cache /proc
func findProcessByInode(inode []byte) (int, string) {
	files, err := os.ReadDir(PROC_PATH)

	if err != nil {
		log.Fatalf("ERROR: Could not read %s because %s", PROC_PATH, err)
	}

	expected := fmt.Sprintf("socket:[%s]", string(inode))

	for _, file := range files {
		n, ok := strconv.Atoi(file.Name())

		if ok != nil {
			continue
		}

		path := fmt.Sprintf("%s%d/", PROC_PATH, n)
		fds, err := os.ReadDir(path + "fd/")

		if err != nil {
			// fmt.Printf("ERROR: Could not read fds because %s\n", err)
			continue
		}
		for _, fd := range fds {
			linkPath := fmt.Sprintf("%sfd/%s", path, fd.Name())
			result, err := os.Readlink(linkPath)

			if err != nil {
				continue
			}

			if result == expected {
				cmdline, _ := os.ReadFile(path + "cmdline")
				fmt.Println(cmdline)
				return n, string(cmdline)
			}

		}
	}
	return -1, "not found"

}

func ParseTCP() [][]byte {
	file, err := os.ReadFile(PROC_NET_PATH + "tcp6")
	port := [][]byte{}

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

		// // those who know
		p.port = chunks[4][len(chunks[4])-4 : len(chunks[4])]
		p.inode = chunks[len(chunks)-8]
		p.pid, p.process_name = findProcessByInode(p.inode)

		fmt.Println(p)
	}

	return port
}
