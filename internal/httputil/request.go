// Copyright 2020 The Go Mail Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httputil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"io"
	"net/http"
	"strings"
	"time"
)

type Requester interface {
	// Do accepts a message, url endpoint and optional headers to POST data
	// to a drivers API.
	// Returns an error if data could not be marshalled/unmarshalled
	// or if the request could not be processed.
	Do(ctx context.Context, r *Request, payload Payload) ([]byte, *http.Response, error)
}

type Client struct {
	Client     *http.Client
	marshaller func(v interface{}) ([]byte, error)
	bodyReader func(r io.Reader) ([]byte, error)
}

type Request struct {
	method            string
	url               string
	headers           map[string]string
	basicAuthUser     string
	basicAuthPassword string
}

const (
	// Timeout is the amount of time to wait before
	// a mail request is cancelled.
	Timeout = time.Second * 10
)

func NewClient() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: Timeout,
		},
		marshaller: json.Marshal,
		bodyReader: io.ReadAll,
	}
}

func NewHTTPRequest(method, url string) *Request {
	return &Request{
		method: method,
		url:    url,
	}
}

func (r *Request) AddHeader(name, value string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[name] = value
}

func (r *Request) SetBasicAuth(user, password string) {
	r.basicAuthUser = user
	r.basicAuthPassword = password
}

func (c *Client) Do(ctx context.Context, r *Request, payload Payload) ([]byte, *http.Response, error) {
	req, err := c.makeRequest(ctx, r, payload)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	buf, err := c.bodyReader(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	return buf, resp, nil
}

func (c *Client) makeRequest(ctx context.Context, r *Request, payload Payload) (*http.Request, error) {
	var body io.Reader = nil
	if payload != nil {
		b, err := payload.Buffer()
		if err != nil {
			return nil, err
		}
		body = b
	}

	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		return nil, err
	}

	if mail.Debug {
		fmt.Println(c.curlString(req, payload))
	}

	req = req.WithContext(ctx)

	if payload != nil && payload.ContentType() != "" {
		req.Header.Add("Content-Type", payload.ContentType())
	}

	if r.basicAuthUser != "" && r.basicAuthPassword != "" {
		req.SetBasicAuth(r.basicAuthUser, r.basicAuthPassword)
	}

	for header, value := range r.headers {
		req.Header.Add(header, value)
	}

	return req, nil
}

func (c *Client) curlString(req *http.Request, p Payload) string {
	parts := []string{"curl", "-i", "-X", req.Method, req.URL.String()}
	for key, value := range req.Header {
		parts = append(parts, fmt.Sprintf("-H \"%s: %s\"", key, value[0]))
	}

	if p != nil {
		for key, value := range p.Values() {
			parts = append(parts, fmt.Sprintf(" -F %s='%s'", key, value))
		}
	}

	return strings.Join(parts, " ")
}
