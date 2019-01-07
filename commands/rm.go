package commands

import (
    "os"
    "log"
    "github.com/sirsean/jc/path"
)

func Rm() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	err = os.RemoveAll(path.RequestDirPath(id))
	if err != nil {
		log.Fatal(err)
	}
}
