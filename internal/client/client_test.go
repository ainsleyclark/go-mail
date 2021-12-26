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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	base := "https://gomail.example.com"
	got := New(base)
	assert.NotNil(t, got.http)
	assert.Equal(t, base, got.baseURL)
	assert.NotNil(t, got.marshaller)
	assert.NotNil(t, got.bodyReader)
}

func TestClient_Do(t *testing.T) {
	tt := map[string]struct {
		input      interface{}
		url        string
		handler    http.HandlerFunc
		marshaller func(v interface{}) ([]byte, error)
		bodyReader func(r io.Reader) ([]byte, error)
		want       interface{}
	}{
		"Success": {
			input: "input",
			url:   "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("buf"))
				assert.NoError(t, err)
			},
			marshaller: json.Marshal,
			bodyReader: io.ReadAll,
			want:       "buf",
		},
		"Marshal Error": {
			input:   "input",
			url:     "",
			handler: nil,
			marshaller: func(v interface{}) ([]byte, error) {
				return nil, fmt.Errorf("marshal error")
			},
			bodyReader: io.ReadAll,
			want:       "marshal error",
		},
		"Bad Request": {
			input:      "input",
			url:        "@#@#$$%$",
			handler:    nil,
			marshaller: json.Marshal,
			bodyReader: io.ReadAll,
			want:       "invalid URL escape",
		},
		"Do Error": {
			input:      "input",
			url:        "wrong",
			handler:    nil,
			marshaller: json.Marshal,
			bodyReader: io.ReadAll,
			want:       "unsupported protocol scheme",
		},
		"Request Error": {
			input: "input",
			url:   "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("buf"))
				assert.NoError(t, err)
			},
			marshaller: json.Marshal,
			bodyReader: io.ReadAll,
			want:       "invalid request",
		},
		"Body Read Error": {
			input: "input",
			url:   "",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("buf"))
				assert.NoError(t, err)
			},
			marshaller: json.Marshal,
			bodyReader: func(r io.Reader) ([]byte, error) {
				return nil, errors.New("body read error")
			},
			want: "body read error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(test.handler)
			defer server.Close()

			url := server.URL
			if test.url != "" {
				url = test.url
			}

			c := Client{
				http:       server.Client(),
				baseURL:    url,
				marshaller: test.marshaller,
				bodyReader: test.bodyReader,
			}

			buf, resp, err := c.Do("input", "", nil)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}

			assert.Equal(t, test.want, string(buf))
			assert.Equal(t, resp.StatusCode, http.StatusOK)
		})
	}
}
