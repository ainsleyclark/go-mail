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

package drivers

import (
	"fmt"
	"log"
	"net/http"

	mocks "github.com/flightaware/go-mail/internal/mocks/client"
	"github.com/flightaware/go-mail/mail"
)

func ExampleNewSparkPost() {
	cfg := mail.Config{
		URL:         "https://api.eu.sparkpost.com", // Or https://api.sparkpost.com/api/v1
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
	}

	_, err := NewSparkPost(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func (t *DriversTestSuite) TestNewSparkPost() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
				URL:         "https://api.eu.sparkpost.com",
				APIKey:      "key",
				FromAddress: "addr",
				FromName:    "name",
			},
			nil,
		},
		"Validation Failed": {
			mail.Config{},
			"driver requires from address",
		},
		"Error": {
			mail.Config{
				URL:         "http://",
				APIKey:      "key",
				FromAddress: "addr",
				FromName:    "name",
			},
			"API base url must be https!",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := NewSparkPost(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestSSparkPostResponse_Unmarshal() {
	t.UtilTestUnmarshal(&spResponse{}, []byte(`{"results": {}}`))
}

func (t *DriversTestSuite) TestSparkPostResponse_CheckError() {
	tt := map[string]struct {
		input    spResponse
		response *http.Response
		buf      []byte
		want     error
	}{
		"Success": {
			spResponse{Errors: nil},
			&http.Response{StatusCode: http.StatusOK},
			[]byte("test"),
			nil,
		},
		"No Errors": {
			spResponse{Errors: []spError{{Message: "message", Code: "code"}}},
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
			mail.ErrEmptyBody,
		},
		"Error": {
			spResponse{Errors: []spError{{Message: "message", Code: "code"}}},
			&http.Response{StatusCode: http.StatusInternalServerError},
			[]byte("test"),
			fmt.Errorf("%s - code: code, message: message", sparkpostErrorMessage),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			err := test.input.CheckError(test.response, test.buf)
			if err != nil {
				t.Contains(err.Error(), test.want.Error())
				return
			}
			t.Equal(test.want, err)
		})
	}
}

func (t *DriversTestSuite) TestSparkPostResponse_Meta() {
	d := &spResponse{
		Results: map[string]interface{}{"id": "10"},
		Errors:  nil,
	}
	t.UtilTestMeta(d, "Successfully sent SparkPost email", "10")
}

func (t *DriversTestSuite) TestSparkPost_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &sparkPost{cfg: Comfig, client: m}
	}, true)
}
