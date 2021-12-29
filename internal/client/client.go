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
	"context"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/httputil"
	"github.com/ainsleyclark/go-mail/mail"
	"io"
	"net/http"
	"strings"
	"time"
)

type Requester interface {
	// Do accepts a message, url endpoint and optional Headers to POST data
	// to a drivers API.
	// Returns an error if data could not be marshalled/unmarshalled
	// or if the request could not be processed.
	Do(ctx context.Context, r *httputil.Request, payload httputil.Payload, responder httputil.Responder) (mail.Response, error)
}

// New creates a new Client with a stdlib http.Client.
func New() *Client {
	return &Client{
		Client: &http.Client{
			Timeout: Timeout,
		},
		bodyReader: io.ReadAll,
	}
}

const (
	// Timeout is the amount of time to wait before
	// a mail request is cancelled.
	Timeout = time.Second * 10
)

// Client defines a http.Client to interact with the mail
// drivers API's. It acts as a reusable helper to send
// data to the drivers endpoints.
type Client struct {
	Client     *http.Client
	bodyReader func(r io.Reader) ([]byte, error)
}

// Do accepts a message, Request and a Payload to POST data
// to a drivers API.
// Logs Curl output if mail.debug is set to true.
//
// Returns an error if data could not be marshalled/unmarshalled
// or if the request could not be processed.
func (c *Client) Do(ctx context.Context, r *httputil.Request, payload httputil.Payload, responder httputil.Responder) (mail.Response, error) {
	req, err := c.makeRequest(ctx, r, payload)
	if err != nil {
		return mail.Response{}, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return mail.Response{}, err
	}
	defer resp.Body.Close()

	response := mail.Response{
		StatusCode: resp.StatusCode,
	}

	buf, err := c.bodyReader(resp.Body)
	if err != nil {
		return response, err
	}
	response.Body = buf

	err = responder.Unmarshal(buf)
	if err != nil {
		return response, err
	}

	err = responder.CheckError(resp, buf)
	if err != nil {
		return response, err
	}

	meta := responder.Meta()
	return mail.Response{
		StatusCode: resp.StatusCode,
		Body:       buf,
		Headers:    resp.Header,
		ID:         meta.ID,
		Message:    meta.Message,
	}, nil
}

// makeRequest creates a stdlib http.Request.
// Content-Type, BasicAuth and headers are attached to the request.
// Returns an error if the request could not be created.
func (c *Client) makeRequest(ctx context.Context, r *httputil.Request, payload httputil.Payload) (*http.Request, error) {
	var body io.Reader = nil
	if payload != nil {
		b, err := payload.Buffer()
		if err != nil {
			return nil, err
		}
		body = b
	}

	req, err := http.NewRequest(r.Method, r.Url, body)
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

	if r.BasicAuthUser != "" && r.BasicAuthPassword != "" {
		req.SetBasicAuth(r.BasicAuthUser, r.BasicAuthPassword)
	}

	for header, value := range r.Headers {
		req.Header.Add(header, value)
	}

	return req, nil
}

// curlString constructs a string used for posting the
// request via Curl.
func (c *Client) curlString(req *http.Request, p httputil.Payload) string {
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
