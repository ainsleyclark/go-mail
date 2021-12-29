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

package clientold

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Requester defines the method used to send data to a mail
// driver API endpoint.
type Requester interface {
	// Do accepts a message, URL endpoint and optional headers to POST data
	// to a drivers API.
	// Returns an error if data could not be marshalled/unmarshalled
	// or if the request could not be processed.
	Do(message interface{}, url string, headers http.Header) ([]byte, *http.Response, error)
}

// Client defines a http.Client to interact with the mail
// drivers API's. It acts as a reusable helper to send
// data to the endpoints.
type Client struct {
	http       *http.Client
	baseURL    string
	marshaller func(v interface{}) ([]byte, error)
	bodyReader func(r io.Reader) ([]byte, error)
}

const (
	// Timeout is the amount of time to wait before
	// a mail request is cancelled.
	Timeout = time.Second * 10
)

// New creates a new Client accepting a baseURL for
// request endpoints.
func New(baseURL string) *Client {
	return &Client{
		http: &http.Client{
			Timeout: Timeout,
		},
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		marshaller: json.Marshal,
		bodyReader: io.ReadAll,
	}
}

func (c *Client) Do(message interface{}, url string, headers http.Header) ([]byte, *http.Response, error) {
	data, err := c.marshaller(message)
	if err != nil {
		return nil, nil, err
	}

	// Setup request with URL, ensures URL's are
	// trimmed.
	url = fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(url, "/"))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	if len(headers) == 0 {
		headers = http.Header{}
	}

	req.Header = headers
	req.Header.Set("User-Agent", "Go Mail v0.1")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	// Read the response body into a buffer for processing using
	// the bodyReader function.
	buf, err := c.bodyReader(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	return buf, resp, nil
}

// Is2XX returns true if the provided HTTP response code is
// in the range 200-299.
func Is2XX(code int) bool {
	if code < 300 && code >= 200 {
		return true
	}
	return false
}
