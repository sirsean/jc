package main

import (
	"fmt"
	"log"
	"os"
)

var basePath = ".jc"

var commandFuncs = map[string]func(){
	"ls":   listCommand,
	"list": listCommand,
	"new":  newCommand,
	"del":  delCommand,
	"rm":   delCommand,
	"run":  runCommand,
	"resp": respCommand,
}

func listCommand() {
	ListRequests()
}

func newCommand() {
	filename, err := NewRequest(getArgId())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filename)
}

func delCommand() {
	err := DeleteRequest(getArgId())
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand() {
	r, err := LoadRequest(getArgId())
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
	c, err := LoadResponse(getArgId())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(c))
}

func getArgCommand() string {
	return os.Args[1]
}

func getArgId() string {
	return os.Args[2]
}

func main() {
	// make sure the base path is present
	os.MkdirAll(basePath, os.ModePerm)

	if f, ok := commandFuncs[getArgCommand()]; ok {
		f()
	} else {
		fmt.Println("unknown command")
	}
}
