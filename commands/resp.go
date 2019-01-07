package commands

import (
    "fmt"
    "log"
	"github.com/sirsean/jc/path"
	"io/ioutil"
)


func Resp() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	c, err := loadResponse(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(c))
}

func loadResponse(id string) ([]byte, error) {
	return ioutil.ReadFile(path.ResponsePath(id))
}
