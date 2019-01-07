package main

import (
	"fmt"
	"github.com/sirsean/jc/path"
	"github.com/sirsean/jc/commands"
	"os"
)

var commandFuncs = map[string]commands.Command{
	"help": commands.Help,
	"ls":   commands.Ls,
	"list": commands.Ls,
	"new":  commands.New,
	"cp":   commands.Copy,
	"del":  commands.Rm,
	"rm":   commands.Rm,
	"run":  commands.Run,
	"resp": commands.Resp,
	"req":  commands.Req,
	"body": commands.Body,
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
	path.MakeBasePath()

	if f, ok := commandFuncs[getArgCommand()]; ok {
		f()
	} else {
		fmt.Println("unknown command")
		commands.Help()
	}
}
