package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirsean/go-pool"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"time"
)

func RequestJsonPath(id string) string {
	return path.Join(basePath, id, "request.json")
}

func RequestBodyJsonPath(id string) string {
	return path.Join(basePath, id, "body.json")
}

func NewRequest(id string) (string, error) {
	requestDirPath := path.Join(basePath, id)
	os.MkdirAll(requestDirPath, os.ModePerm)
	requestJsonPath := path.Join(basePath, id, "request.json")
	r := Request{
		Id:      id,
		Timeout: 60,
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
	raw, err := ioutil.ReadFile(RequestJsonPath(id))
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(raw, &r)
	if r.Timeout <= 0 {
		r.Timeout = 60
	}
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

	transport := &http.Transport{}

	if r.IsClientCert() {
		tlsConfig, err := r.GetTlsConfig()
		if err != nil {
			return "", err
		}
		transport.TLSClientConfig = tlsConfig
	}

	httpClient := http.Client{
		Timeout:   r.Timeout * time.Second,
		Transport: transport,
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

type LoadRequestWorkUnit struct {
	Id    string
	OutCh chan Request
}

func (u LoadRequestWorkUnit) Perform() {
	r, err := LoadRequest(u.Id)
	if err == nil {
		u.OutCh <- r
	}
}

func ListRequests() error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	p := pool.NewPool(100, 100) // we will read many files at a time
	p.Start()

	ch := make(chan Request, len(files)) // channel for each request
	rCh := make(chan []Request, 1)       // result channel for full requests list

	go func() {
		// we will build up the list of requests in this goroutine
		requests := make([]Request, 0)
		for r := range ch {
			requests = append(requests, r)
		}
		// all the requests are in the list, sort it
		sort.Slice(requests[:], func(i, j int) bool {
			return requests[i].Id < requests[j].Id
		})
		// send the list back to the main thread
		rCh <- requests
	}()

	for _, f := range files {
		p.Add(LoadRequestWorkUnit{
			Id:    f.Name(),
			OutCh: ch,
		})
	}

	p.Close()
	close(ch)

	requests := <-rCh
	for _, r := range requests {
		fmt.Printf("%s: %s %s\n", r.Id, r.Method, r.Url)
	}

	return nil
}

func LoadResponse(id string) ([]byte, error) {
	responseJsonPath := path.Join(basePath, id, "response.json")
	return ioutil.ReadFile(responseJsonPath)
}
