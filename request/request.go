package request

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/sirsean/jc/path"
	j "github.com/sirsean/jc/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		// if given as a number, treat as seconds
		d.Duration = time.Duration(value) * time.Second
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func NewRequest(id string) (string, error) {
	requestDirPath := path.RequestDirPath(id)
	os.MkdirAll(requestDirPath, os.ModePerm)
	requestJsonPath := path.RequestPath(id)
	r := Request{
		Id:      id,
		Timeout: Duration{60 * time.Second},
	}
	err := r.Write(requestJsonPath)
	return requestJsonPath, err
}

func LoadRequest(id string) (Request, error) {
	var r Request
	raw, err := ioutil.ReadFile(path.RequestPath(id))
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(raw, &r)
	if r.Timeout.Seconds() <= 0 {
		r.Timeout = Duration{60 * time.Second}
	}
	return r, err
}

type Request struct {
	Id         string            `json:"id"`
	Url        string            `json:"url"`
	Method     string            `json:"method"`
	BasicAuth  BasicAuth         `json:"basic_auth"`
	ClientCert ClientCert        `json:"client_cert"`
	Headers    map[string]string `json:"headers"`
	Timeout    Duration          `json:"timeout"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientCert struct {
	CaCertPath     string `json:"ca_cert"`
	ClientCertPath string `json:"client_cert"`
	ClientKeyPath  string `json:"client_key"`
}

func (r Request) IsBasicAuth() bool {
	return r.BasicAuth.Username != "" && r.BasicAuth.Password != ""
}

func (r Request) IsClientCert() bool {
	return r.ClientCert.CaCertPath != "" &&
		r.ClientCert.ClientCertPath != "" &&
		r.ClientCert.ClientKeyPath != ""
}

func (r Request) GetTlsConfig() (*tls.Config, error) {
	c := r.ClientCert
	caCert, err := ioutil.ReadFile(expandHome(c.CaCertPath))
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(expandHome(c.ClientCertPath), expandHome(c.ClientKeyPath))
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	config := tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
	}

	return &config, nil
}

func expandHome(in string) string {
	return strings.Replace(in, "~", os.Getenv("HOME"), -1)
}

func (r Request) Body() io.Reader {
	raw, err := ioutil.ReadFile(path.RequestBodyPath(r.Id))
	if err != nil {
		return nil
	}
	return bytes.NewReader(raw)
}

func (r Request) Write(filename string) error {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return j.Write(filename, jsonBytes)
}
