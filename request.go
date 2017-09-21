package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"path"
)

type Request struct {
	Id        string            `json:"id"`
	Url       string            `json:"url"`
	Method    string            `json:"method"`
	BasicAuth BasicAuth         `json:"basic_auth"`
	Headers   map[string]string `json:"headers"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r Request) IsBasicAuth() bool {
	return r.BasicAuth.Username != "" && r.BasicAuth.Password != ""
}

func (r Request) Body() io.Reader {
	bodyPath := path.Join(basePath, r.Id, "body.json")
	raw, err := ioutil.ReadFile(bodyPath)
	if err != nil {
		return nil
	}
	return bytes.NewReader(raw)
}
