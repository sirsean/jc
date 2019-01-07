package commands

import (
    "log"
    "fmt"
    "github.com/sirsean/jc/request"
)

func New() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	filename, err := request.NewRequest(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filename)
}
