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

func ExampleNewPostmark() {
	cfg := mail.Config{
		URL:         "https://postal.example.com",
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
	}

	_, err := NewPostal(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func (t *DriversTestSuite) TestNewPostmark() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
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
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := NewPostmark(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestPostmarkResponse_Unmarshal() {
	t.UtilTestUnmarshal(&postmarkResponse{}, []byte(`{"message": "Hello"}`))
}

func (t *DriversTestSuite) TestPostmarkResponse_CheckError() {
	tt := map[string]struct {
		input    postmarkResponse
		response *http.Response
		buf      []byte
		want     error
	}{
		"Success": {
			postmarkResponse{ErrorCode: 0},
			&http.Response{StatusCode: http.StatusOK},
			[]byte("test"),
			nil,
		},
		"Empty Body": {
			postmarkResponse{ErrorCode: 10},
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
			mail.ErrEmptyBody,
		},
		"Error": {
			postmarkResponse{ErrorCode: 10, Message: "message"},
			&http.Response{StatusCode: http.StatusInternalServerError},
			[]byte("test"),
			fmt.Errorf("%s - code: 10, message: message", postmarkErrorMessage),
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

func (t *DriversTestSuite) TestPostmarkResponse_Meta() {
	d := &postmarkResponse{Message: "Success", ID: "id"}
	t.UtilTestMeta(d, d.Message, d.ID)
}

func (t *DriversTestSuite) TestPostmark_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &postmark{cfg: Comfig, client: m}
	}, true)
}
