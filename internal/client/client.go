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

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type Requester interface {
	Do(message interface{}, url string, headers http.Header) ([]byte, int, error)
}

type Client struct {
	http *http.Client
	marshaller func(v interface{}) ([]byte, error)
	bodyReader func(r io.Reader) ([]byte, error)
}

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout: time.Second * 10,
		},
		marshaller: json.Marshal,
		bodyReader: io.ReadAll,
	}
}

func (c *Client) Do(message interface{}, url string, headers http.Header) ([]byte, int, error) {
	data, err := c.marshaller(message)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, err
	}
	req.Header = headers

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Successful response
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Read the response body into a buffer for processing using
		// the bodyReader function.
		buf, err := c.bodyReader(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, err
		}
		return buf, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, errors.New("go-mail client: invalid request")

}
