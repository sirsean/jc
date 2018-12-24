package json

import (
	"bytes"
	"encoding/json"
    "io/ioutil"
)

func Pretty(in []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "    ")
	if err != nil {
		return in
	}
	return out.Bytes()
}

func Write(filename string, bytes []byte) error {
	return ioutil.WriteFile(filename, Pretty(bytes), 0666)
}
