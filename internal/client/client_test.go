// Copyright 2022 Ainsley Clark. All rights reserved.
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
	"context"
	"github.com/ainsleyclark/go-mail/internal/errors"
	"github.com/ainsleyclark/go-mail/internal/httputil"
	mocks "github.com/ainsleyclark/go-mail/internal/mocks/httputil"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	got := New(nil)
	assert.NotNil(t, got.bodyReader)
	c := &http.Client{}
	withClient := New(c)
	assert.Equal(t, withClient.Client, c)
}

func TestClient_Do(t *testing.T) {
	successHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("buf"))
		assert.NoError(t, err)
	}

	tt := map[string]struct {
		input      *httputil.Request
		handler    http.HandlerFunc
		responder  func(m *mocks.Responder)
		bodyReader func(r io.Reader) ([]byte, error)
		want       interface{}
	}{
		"Success": {
			handler: successHandler,
			responder: func(m *mocks.Responder) {
				m.On("Unmarshal", mock.Anything).
					Return(nil)
				m.On("CheckError", mock.Anything, []byte("buf")).
					Return(nil)
				m.On("Meta", mock.Anything).
					Return(httputil.Meta{Message: "message", ID: "10"})
			},
			bodyReader: io.ReadAll,
			want: mail.Response{
				StatusCode: http.StatusOK,
				Body:       []byte("buf"),
				Headers:    nil,
				ID:         "10",
				Message:    "message",
			},
		},
		"Bad Request": {
			input: &httputil.Request{URL: "@#@#$$%$"},
			want:  "Error creating http request",
		},
		"Do Error": {
			input: &httputil.Request{URL: "wrong"},
			want:  "Error doing request",
		},
		"Body Read Error": {
			handler: successHandler,
			bodyReader: func(r io.Reader) ([]byte, error) {
				return nil, errors.New("body read error")
			},
			want: "Error reading response body",
		},
		"Unmarshal Error": {
			handler: successHandler,
			responder: func(m *mocks.Responder) {
				m.On("Unmarshal", mock.Anything).
					Return(errors.New("unmarshal error"))
			},
			bodyReader: io.ReadAll,
			want:       "Error unmarshalling response error",
		},
		"Responder Error": {
			handler: successHandler,
			responder: func(m *mocks.Responder) {
				m.On("Unmarshal", mock.Anything).
					Return(nil)
				m.On("CheckError", mock.Anything, []byte("buf")).
					Return(errors.New("response error"))
			},
			bodyReader: io.ReadAll,
			want:       "Error performing mail request",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(test.handler)
			defer server.Close()

			if test.input == nil {
				test.input = &httputil.Request{}
			}

			if test.input.URL == "" {
				test.input.URL = server.URL
			}

			responder := &mocks.Responder{}
			if test.responder != nil {
				test.responder(responder)
			}

			c := Client{
				Client:     server.Client(),
				bodyReader: test.bodyReader,
			}

			got, err := c.Do(context.Background(), test.input, nil, responder)
			if err != nil {
				assert.Contains(t, errors.Message(err), test.want)
				return
			}

			want := test.want.(mail.Response)

			assert.Equal(t, want.StatusCode, got.StatusCode)
			assert.Equal(t, want.Body, got.Body)
			assert.NotEmpty(t, got.Headers)
			assert.Equal(t, want.ID, got.ID)
			assert.Equal(t, want.Message, got.Message)
		})
	}
}

func TestClient_MakeRequest(t *testing.T) {
	uri, err := url.Parse("https://gomail.example.com")
	assert.NoError(t, err)

	tt := map[string]struct {
		request *httputil.Request
		payload func(m *mocks.Payload)
		want    interface{}
	}{
		"Success": {
			&httputil.Request{
				Method:            http.MethodPost,
				URL:               "https://gomail.example.com",
				BasicAuthUser:     "user",
				BasicAuthPassword: "password",
				Headers:           map[string]string{"header": "Value"},
			},
			func(m *mocks.Payload) {
				m.On("Buffer").
					Return(&bytes.Buffer{}, nil)
				m.On("ContentType").
					Return(httputil.JSONContentType)
				m.On("Values").
					Return(map[string]string{"key": "value"})
			},
			&http.Request{
				Method: http.MethodPost,
				URL:    uri,
				Header: map[string][]string{
					"Authorization": {"Basic dXNlcjpwYXNzd29yZA=="},
					"Content-Type":  {httputil.JSONContentType},
					"Header":        {"Value"},
				},
			},
		},
		"Buffer Error": {
			&httputil.Request{},
			func(m *mocks.Payload) {
				m.On("Buffer").
					Return(nil, &errors.Error{Message: "buffer error"})
			},
			"buffer error",
		},
		"Request Error": {
			&httputil.Request{
				URL: "@#@#$$%$",
			},
			func(m *mocks.Payload) {
				m.On("Buffer").
					Return(&bytes.Buffer{}, nil)
			},
			"Error creating http request",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			defer func() { mail.Debug = false }()
			mail.Debug = true

			c := New(nil)

			mock := &mocks.Payload{}
			if test.payload != nil {
				test.payload(mock)
			}

			request, err := c.makeRequest(context.Background(), test.request, mock)
			if err != nil {
				assert.Contains(t, errors.Message(err), test.want)
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

	payload := mocks.Payload{}
	payload.On("Values").
		Return(map[string]string{"key": "value"})

	c := Client{}
	got := c.curlString(req, &payload)
	want := "curl -i -X GET https://gomail.example.com -H \"header: value\"  -F key='value'"
	assert.Equal(t, want, got)
}
