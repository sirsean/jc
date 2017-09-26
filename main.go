package main

import (
	"fmt"
	"log"
	"os"
)

var basePath = ".jc"

var commandFuncs = map[string]func(){
	"help": helpCommand,
	"ls":   listCommand,
	"list": listCommand,
	"new":  newCommand,
	"del":  delCommand,
	"rm":   delCommand,
	"run":  runCommand,
	"resp": respCommand,
}

func helpCommand() {
	fmt.Println(`
	Usage: jc <command> (<id>)

	Commands:
		- ls/list
		- new
		- del/rm
		- run
		- resp
	`)
}

func listCommand() {
	ListRequests()
}

func newCommand() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	filename, err := NewRequest(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filename)
}

func delCommand() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	err = DeleteRequest(id)
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	r, err := LoadRequest(id)
	if err != nil {
		log.Fatal(err)
	}
	filename, err := Execute(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filename)
}

func respCommand() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	c, err := LoadResponse(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(c))
}

func getArgCommand() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	} else {
		return ""
	}
}

func getArgId() (string, error) {
	if len(os.Args) >= 3 {
		return os.Args[2], nil
	} else {
		return "", fmt.Errorf("invalid id")
	}
}

func main() {
	// make sure the base path is present
	os.MkdirAll(basePath, os.ModePerm)

	if f, ok := commandFuncs[getArgCommand()]; ok {
		f()
	} else {
		fmt.Println("unknown command")
		helpCommand()
	}
}
