package commands

import (
    "fmt"
)

func Help() {
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
