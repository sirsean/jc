package commands

import (
    "log"
    "github.com/sirsean/jc/request"
    "github.com/sirsean/jc/path"
    "os"
    "io/ioutil"
)

func Copy() {
	fromId, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	toId, err := getArgAt(3)
	if err != nil {
		log.Fatal(err)
	}

	err = copyRequest(fromId, toId)
	if err != nil {
		log.Fatal(err)
	}
}

func copyRequest(fromId, toId string) error {
	from, err := request.LoadRequest(fromId)
	from.Id = toId

	requestDirPath := path.RequestDirPath(toId)
	os.MkdirAll(requestDirPath, os.ModePerm)

	err = from.Write(path.RequestPath(toId))
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(path.RequestBodyPath(fromId))
	if err == nil {
		// no error means there was a body file, so copy it
		err = ioutil.WriteFile(path.RequestBodyPath(toId), raw, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}
