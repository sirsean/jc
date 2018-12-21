package main

import (
	"fmt"
	"github.com/codeskyblue/go-sh"
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
	"req":  reqCommand,
	"body": bodyCommand,
}

func helpCommand() {
	fmt.Println(`
	Usage: jc <command> (<id>)

	Commands:
		- ls/list
		- new
		- del/rm
		- run
		- resp (print response)
		- req (edit request)
		- body (edit request body)
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

	sh.Command("less", filename).Run()
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

func reqCommand() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	path := RequestJsonPath(id)
	sh.Command("vim", path).SetStdin(os.Stdin).Run()
}

func bodyCommand() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	path := RequestBodyJsonPath(id)
	sh.Command("vim", path).SetStdin(os.Stdin).Run()
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
