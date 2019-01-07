package commands

import (
    "log"
	"github.com/codeskyblue/go-sh"
	"bytes"
	"fmt"
	"github.com/sirsean/jc/json"
	"github.com/sirsean/jc/path"
    "github.com/sirsean/jc/request"
	"net/http"
	"time"
)

func Run() {
	id, err := getArgId()
	if err != nil {
		log.Fatal(err)
	}
	r, err := request.LoadRequest(id)
	if err != nil {
		log.Fatal(err)
	}
	filename, err := execute(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(filename)

	sh.Command("less", filename).Run()
}

func execute(r request.Request) (string, error) {
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

	responseJsonPath := path.ResponsePath(r.Id)
	err = json.Write(responseJsonPath, buf.Bytes())
	return responseJsonPath, err
}
