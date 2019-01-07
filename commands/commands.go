package commands

import (
    "os"
    "fmt"
)

type Command func()

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
