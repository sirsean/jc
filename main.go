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
	"cp":   copyCommand,
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
		- cp <from-id> <to-id>
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

func copyCommand() {
	fromId, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	toId, err := getArgAt(3)
	if err != nil {
		log.Fatal(err)
	}

	err = CopyRequest(fromId, toId)
	if err != nil {
		log.Fatal(err)
	}
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

func getArgAt(index int) (string, error) {
	if len(os.Args) >= (index + 1) {
		return os.Args[index], nil
	} else {
		return "", fmt.Errorf("invalid arg index")
	}
}

func getArgId() (string, error) {
	if id, err := getArgAt(2); err != nil {
		return "", fmt.Errorf("invalid id")
	} else {
		return id, err
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
