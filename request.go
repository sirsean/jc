package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"path"
)

type Request struct {
	Id         string            `json:"id"`
	Url        string            `json:"url"`
	Method     string            `json:"method"`
	BasicAuth  BasicAuth         `json:"basic_auth"`
	ClientCert ClientCert        `json:"client_cert"`
	Headers    map[string]string `json:"headers"`
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
	caCert, err := ioutil.ReadFile(c.CaCertPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(c.ClientCertPath, c.ClientKeyPath)
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

func (r Request) Body() io.Reader {
	bodyPath := path.Join(basePath, r.Id, "body.json")
	raw, err := ioutil.ReadFile(bodyPath)
	if err != nil {
		return nil
	}
	return bytes.NewReader(raw)
}
