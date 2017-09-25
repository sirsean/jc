package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func NewRequest(id string) (string, error) {
	requestDirPath := path.Join(basePath, id)
	os.MkdirAll(requestDirPath, os.ModePerm)
	requestJsonPath := path.Join(basePath, id, "request.json")
	r := Request{
		Id: id,
	}
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(requestJsonPath, prettyJson(jsonBytes), 0666)
	return requestJsonPath, err
}

func DeleteRequest(id string) error {
	requestDirPath := path.Join(basePath, id)
	return os.RemoveAll(requestDirPath)
}

func LoadRequest(id string) (Request, error) {
	var r Request
	requestJsonPath := path.Join(basePath, id, "request.json")
	raw, err := ioutil.ReadFile(requestJsonPath)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(raw, &r)
	return r, err
}

func Execute(r Request) (string, error) {
	log.Println("executing", r)
	start := time.Now()
	req, err := http.NewRequest(r.Method, r.Url, r.Body())
	if err != nil {
		return "", err
	}
	if r.IsBasicAuth() {
		req.SetBasicAuth(r.BasicAuth.Username, r.BasicAuth.Password)
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()

	d := time.Since(start)
	fmt.Println(d)

	responseJsonPath := path.Join(basePath, r.Id, "response.json")
	err = ioutil.WriteFile(responseJsonPath, prettyJson(buf.Bytes()), 0666)
	return responseJsonPath, err
}

func ListRequests() error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, f := range files {
		id := f.Name()
		r, err := LoadRequest(id)
		if err != nil {
			fmt.Printf("%s: ?\n", id)
		} else {
			fmt.Printf("%s: %s %s\n", id, r.Method, r.Url)
		}
	}

	return nil
}

func LoadResponse(id string) ([]byte, error) {
	responseJsonPath := path.Join(basePath, id, "response.json")
	return ioutil.ReadFile(responseJsonPath)
}
