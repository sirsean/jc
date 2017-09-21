package main

import (
	"bytes"
	"encoding/json"
)

func prettyJson(in []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, in, "", "    ")
	if err != nil {
		return in
	}
	return out.Bytes()
}
