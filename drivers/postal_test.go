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
	mocks "github.com/ainsleyclark/go-mail/internal/mocks/client"
	"github.com/ainsleyclark/go-mail/mail"
	"log"
	"net/http"
)

func ExampleNewPostal() {
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
			"driver requires from address",
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

func (t *DriversTestSuite) TestPostalResponse_Unmarshal() {
	t.UtilTestUnmarshal(&postalResponse{}, []byte(`{"status": "success"}`))
}

func (t *DriversTestSuite) TestPostalResponse_CheckError() {
	tt := map[string]struct {
		input    postalResponse
		response *http.Response
		buf      []byte
		want     error
	}{
		"Success": {
			postalResponse{Status: "success"},
			&http.Response{StatusCode: http.StatusOK},
			[]byte("test"),
			nil,
		},
		"Empty Body": {
			postalResponse{},
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
			mail.ErrEmptyBody,
		},
		"Error": {
			postalResponse{Data: map[string]interface{}{"code": "code", "message": "message"}},
			&http.Response{StatusCode: http.StatusInternalServerError},
			[]byte("test"),
			fmt.Errorf("%s - code: code, message: message", postalErrorMessage),
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

func (t *DriversTestSuite) TestPostalResponse_Meta() {
	d := &postalResponse{
		Data: map[string]interface{}{"message_id": 10},
	}
	t.UtilTestMeta(d, "Successfully sent Postal email", "10")
}

func (t *DriversTestSuite) TestPostal_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &postal{cfg: Comfig, client: m}
	}, true)
}
