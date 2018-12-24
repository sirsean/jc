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

func writeRequestJson(filename string, req Request) error {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, prettyJson(jsonBytes), 0666)
}

func NewRequest(id string) (string, error) {
	requestDirPath := path.Join(basePath, id)
	os.MkdirAll(requestDirPath, os.ModePerm)
	requestJsonPath := RequestJsonPath(id)
	r := Request{
		Id:      id,
		Timeout: Duration{60 * time.Second},
	}
	err := writeRequestJson(requestJsonPath, r)
	return requestJsonPath, err
}

func CopyRequest(fromId, toId string) error {
	from, err := LoadRequest(fromId)
	from.Id = toId

	requestDirPath := path.Join(basePath, toId)
	os.MkdirAll(requestDirPath, os.ModePerm)

	err = writeRequestJson(RequestJsonPath(toId), from)
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(RequestBodyJsonPath(fromId))
	if err == nil {
		// no error means there was a body file, so copy it
		err = ioutil.WriteFile(RequestBodyJsonPath(toId), raw, 0666)
		if err != nil {
			return err
		}
	}

	return nil
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
	if r.Timeout.Seconds() <= 0 {
		r.Timeout = Duration{60 * time.Second}
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
		Timeout:   r.Timeout.Duration,
		Transport: transport,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()

	fmt.Println(resp.Status)

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
