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
	"bytes"
	"context"
	"errors"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/ainsleyclark/go-mail/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	got := NewClient()
	assert.NotNil(t, got.bodyReader)
}

func TestClient_Do(t *testing.T) {
	tt := map[string]struct {
		input      *Request
		handler    http.HandlerFunc
		bodyReader func(r io.Reader) ([]byte, error)
		want       interface{}
	}{
		"Success": {
			input: &Request{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("buf"))
				assert.NoError(t, err)
			},
			bodyReader: io.ReadAll,
			want:       "buf",
		},
		"Bad Request": {
			input:      &Request{Url: "@#@#$$%$"},
			handler:    nil,
			bodyReader: io.ReadAll,
			want:       "invalid URL escape",
		},
		"Do Error": {
			input:      &Request{Url: "wrong"},
			handler:    nil,
			bodyReader: io.ReadAll,
			want:       "unsupported protocol scheme",
		},
		"Body Read Error": {
			input: &Request{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("buf"))
				assert.NoError(t, err)
			},
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

			if test.input.Url == "" {
				test.input.Url = server.URL
			}

			c := Client{
				Client:     server.Client(),
				bodyReader: test.bodyReader,
			}

			buf, resp, err := c.Do(context.Background(), test.input, nil)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}

			assert.Equal(t, test.want, string(buf))
			assert.Equal(t, resp.StatusCode, http.StatusOK)
		})
	}
}

func TestClient_MakeRequest(t *testing.T) {
	uri, err := url.Parse("https://gomail.example.com")
	assert.NoError(t, err)

	tt := map[string]struct {
		request *Request
		payload func(m *mocks.Payload)
		want    interface{}
	}{
		"Success": {
			&Request{
				Method:            http.MethodPost,
				Url:               "https://gomail.example.com",
				BasicAuthUser:     "user",
				BasicAuthPassword: "password",
				Headers:           map[string]string{"header": "Value"},
			},
			func(m *mocks.Payload) {
				m.On("Buffer").Return(&bytes.Buffer{}, nil)
				m.On("ContentType").Return(JSONContentType)
				m.On("Values").Return(map[string]string{"key": "value"})
			},
			&http.Request{
				Method: http.MethodPost,
				URL:    uri,
				Header: map[string][]string{
					"Authorization": {"Basic dXNlcjpwYXNzd29yZA=="},
					"Content-Type":  {JSONContentType},
					"Header":        {"Value"},
				},
			},
		},
		"Buffer Error": {
			&Request{},
			func(m *mocks.Payload) {
				m.On("Buffer").Return(nil, errors.New("buffer error"))
			},
			"buffer error",
		},
		"Request Error": {
			&Request{
				Url: "@#@#$$%$",
			},
			func(m *mocks.Payload) {
				m.On("Buffer").
					Return(&bytes.Buffer{}, nil)
			},
			"invalid URL escape",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			defer func() { mail.Debug = false }()
			mail.Debug = true

			c := NewClient()

			mock := &mocks.Payload{}
			if test.payload != nil {
				test.payload(mock)
			}

			request, err := c.makeRequest(context.Background(), test.request, mock)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}

			want := test.want.(*http.Request)
			assert.Equal(t, want.Method, request.Method)
			assert.Equal(t, want.URL, request.URL)
			assert.Equal(t, want.Header, request.Header)
		})
	}
}

func TestClient_CurlString(t *testing.T) {
	uri, err := url.Parse("https://gomail.example.com")
	assert.NoError(t, err)

	req := &http.Request{
		Method: http.MethodGet,
		URL:    uri,
		Header: map[string][]string{"header": {"value"}},
	}

	mock := mocks.Payload{}
	mock.On("Values").
		Return(map[string]string{"key": "value"})

	c := Client{}
	got := c.curlString(req, &mock)
	want := "curl -i -X GET https://gomail.example.com -H \"header: value\"  -F key='value'"
	assert.Equal(t, want, got)
}
