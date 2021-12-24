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

package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func (t *MailTestSuite) TestNewPostal() {
	tt := map[string]struct {
		input Config
		want  interface{}
	}{
		"Success": {
			Config{
				URL:         "https://postal.example.com",
				APIKey:      "key",
				FromAddress: "addr",
				FromName:    "name",
			},
			nil,
		},
		"Validation Failed": {
			Config{},
			"mailer requires from address",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := newPostal(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.input, got.cfg)
			t.NotNil(got.client)
			t.NotNil(got.marshaller)
			t.NotNil(got.bodyReader)
		})
	}
}

func (t *MailTestSuite) TestPostal_Send() {
	tt := map[string]struct {
		input   *Transmission
		handler http.HandlerFunc
		url string
		marshaller func(v interface{}) ([]byte, error)
		bodyReader func(r io.Reader) ([]byte, error)
		want    interface{}
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
			Response{
				StatusCode: http.StatusOK,
				Body: `{"status":"success","time":0,"flags":null,"data":null}`,
				Message: "Successfully sent Postal email",
			},
		},
		"With ID": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "success",
					Data: map[string]interface{}{"message_id": "1"},
				}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			Response{
				StatusCode: http.StatusOK,
				Body: `{"status":"success","time":0,"flags":null,"data":{"message_id":"1"}}`,
				Message: "Successfully sent Postal email",
				ID: "1",
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
			Response{
				StatusCode: http.StatusOK,
				Body: `{"status":"success","time":0,"flags":null,"data":null}`,
				Message: "Successfully sent Postal email",
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
		"Empty Data Error": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{Status: "error"}
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
		"With Error Code": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "error",
					Data: map[string]interface{}{"code": "ValidationFailed"},
				}
				buf, err := json.Marshal(&res)
				t.NoError(err)
				_, err = w.Write(buf)
				t.NoError(err)
			},
			"",
			json.Marshal,
			io.ReadAll,
			fmt.Sprintf("%s - code: ValidationFailed", postalErrorMessage),
		},
		"With Error Message": {
			Trans,
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				res := postalResponse{
					Status: "error",
					Data: map[string]interface{}{"code": "ValidationFailed", "message": "Postal Message"},
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
				cfg: Config{
					URL:         url,
					FromAddress: "from",
				},
				client: server.Client(),
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
