package main

import (
	"fmt"
	"log"
	"os"
)

var basePath = ".jc"

type command string

const (
	LS   command = "ls"
	DEL  command = "del"
	NEW  command = "new"
	RUN  command = "run"
	RESP command = "resp"
)

func getArgCommand() command {
	return command(os.Args[1])
}

func getArgId() string {
	return os.Args[2]
}

func main() {
	// make sure the base path is present
	os.MkdirAll(basePath, os.ModePerm)

	switch getArgCommand() {
	case LS:
		ListRequests()
	case NEW:
		filename, err := NewRequest(getArgId())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(filename)
	case DEL:
		err := DeleteRequest(getArgId())
		if err != nil {
			log.Fatal(err)
		}
	case RUN:
		r, err := LoadRequest(getArgId())
		if err != nil {
			log.Fatal(err)
		}
		filename, err := Execute(r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(filename)
	case RESP:
		c, err := LoadResponse(getArgId())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(c))
	}
}
