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

package drivers

import (
	"encoding/json"
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"io"
	"net/http"
	"net/http/httptest"
)

func (t *DriversTestSuite) TestNewPostal() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
				URL:         "https://postal.example.com",
				APIKey:      "key",
				FromAddress: "addr",
				FromName:    "name",
			},
			nil,
		},
		"Validation Failed": {
			mail.Config{},
			"mailer requires from address",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := NewPostal(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestPostalResponse_HasError() {
	tt := map[string]struct {
		input postalResponse
		want  bool
	}{
		"Error": {
			postalResponse{Status: "success"},
			false,
		},
		"No Error": {
			postalResponse{Status: "error"},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.HasError()
			t.Equal(test.want, got)
		})
	}
}

func (t *DriversTestSuite) TestPostalResponse_Error() {
	tt := map[string]struct {
		input postalResponse
		want  string
	}{
		"Default": {
			postalResponse{},
			postalErrorMessage,
		},
		"Code": {
			postalResponse{Data: map[string]interface{}{"code": "ValidationFailed"}},
			fmt.Sprintf("%s - code: ValidationFailed", postalErrorMessage),
		},
		"All": {
			postalResponse{Data: map[string]interface{}{"code": "ValidationFailed", "message": "Postal Message"}},
			fmt.Sprintf("%s - code: ValidationFailed, message: Postal Message", postalErrorMessage),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.Error()
			t.Contains(got.Error(), test.want)
		})
	}
}

func (t *DriversTestSuite) TestPostalResponse_ToResponse() {
	tt := map[string]struct {
		input []byte
		resp  postalResponse
		want  mail.Response
	}{
		"Default": {
			[]byte("body"),
			postalResponse{},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       "body",
				Message:    "Successfully sent Postal email",
			},
		},
		"With ID": {
			[]byte("body"),
			postalResponse{Data: map[string]interface{}{"message_id": "1"}},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       "body",
				Message:    "Successfully sent Postal email",
				ID:         "1",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.resp.ToResponse(test.input)
			t.Equal(test.want, got)
		})
	}
}

func (t *DriversTestSuite) TestPostal_Send() {
	t.T().Skip()

	tt := map[string]struct {
		input      *mail.Transmission
		handler    http.HandlerFunc
		url        string
		marshaller func(v interface{}) ([]byte, error)
		bodyReader func(r io.Reader) ([]byte, error)
		want       interface{}
	}{
		"Success": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "success",
				}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"status":"success","time":0,"flags":null,"data":null}`,
				Message:    "Successfully sent Postal email",
			},
		},
		"With ID": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "success",
					Data:   map[string]interface{}{"message_id": "1"},
				}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"status":"success","time":0,"flags":null,"data":{"message_id":"1"}}`,
				Message:    "Successfully sent Postal email",
				ID:         "1",
			},
		},
		"With Attachment": {
			TransWithAttachment,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "success",
				}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"status":"success","time":0,"flags":null,"data":null}`,
				Message:    "Successfully sent Postal email",
			},
		},
		"Validation Failed": {
			nil,
			nil,
			"",
			json.Marshal,
			io.ReadAll,
			"can't validate a nil transmission",
		},
		"Marshal Error": {
			Trans,
			nil,
			"",
			func(v interface{}) ([]byte, error) {
				return nil, fmt.Errorf("marshal error")
			},
			io.ReadAll,
			"marshal error",
		},
		"Bad Request": {
			Trans,
			nil,
			"@#@#$$%$",
			json.Marshal,
			io.ReadAll,
			"invalid URL",
		},
		"Do Error": {
			Trans,
			nil,
			"wrong",
			json.Marshal,
			io.ReadAll,
			"unsupported protocol scheme",
		},
		"Read Error": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("wrong"))
				t.NoError(err)
			},
			"",
			json.Marshal,
			func(r io.Reader) ([]byte, error) {
				return nil, fmt.Errorf("read error")
			},
			"read error",
		},
		"Decode Error": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte("wrong"))
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			"invalid character",
		},
		"Server Error": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				res := postalResponse{}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			postalErrorMessage,
		},
		"Postal Error": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "error",
					Data:   map[string]interface{}{"code": "ValidationFailed", "message": "Postal Message"},
				}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			fmt.Sprintf("%s - code: ValidationFailed, message: Postal Message", postalErrorMessage),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			server := httptest.NewServer(test.handler)
			defer server.Close()

			url := server.URL
			if test.url != "" {
				url = test.url
			}

			ptl := postal{
				cfg: mail.Config{
					URL:         url,
					FromAddress: "from",
				},
				// TODO, replace with interface
				//client:     server.Client(),
				marshaller: test.marshaller,
				bodyReader: test.bodyReader,
			}

			resp, err := ptl.Send(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, resp)
		})
	}
}
