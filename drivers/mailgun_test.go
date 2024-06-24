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
	"errors"
	"log"
	"net/http"

	mocks "github.com/flightaware/go-mail/internal/mocks/client"
	"github.com/flightaware/go-mail/mail"
)

func ExampleNewMailgun() {
	cfg := mail.Config{
		URL:         "https://api.eu.mailgun.net", // Or https://api.mailgun.net
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
		Domain:      "my-domain.com",
	}

	_, err := NewMailgun(cfg)
	if err != nil {
		log.Fatalln(err)
	}
}

func (t *DriversTestSuite) TestNewMailGun() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
				URL:         "https://mailgun.example.com",
				FromAddress: "addr",
				FromName:    "name",
				APIKey:      "key",
				Domain:      "domain",
			},
			nil,
		},
		"Validation Failed": {
			mail.Config{},
			"driver requires from address",
		},
		"No Domain": {
			mail.Config{
				FromName:    "name",
				FromAddress: "hello@gophers.com",
				APIKey:      "key",
			},
			"driver requires a domain",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := NewMailgun(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestMailgunResponse_Unmarshal() {
	t.UtilTestUnmarshal(&mailgunResponse{}, []byte(`{"message": "Hello"}`))
}

func (t *DriversTestSuite) TestMailgunResponse_CheckError() {
	tt := map[string]struct {
		response *http.Response
		buf      []byte
		want     error
	}{
		"2xx": {
			&http.Response{StatusCode: http.StatusOK},
			[]byte("test"),
			nil,
		},
		"Empty Body": {
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
			mail.ErrEmptyBody,
		},
		"Error": {
			&http.Response{StatusCode: http.StatusInternalServerError},
			[]byte("test"),
			errors.New("error"),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			resp := mailgunResponse{Message: "error"}
			err := resp.CheckError(test.response, test.buf)
			if err != nil {
				t.Contains(err.Error(), test.want.Error())
				return
			}
			t.Equal(test.want, err)
		})
	}
}

func (t *DriversTestSuite) TestMailgunResponse_Meta() {
	d := &mailgunResponse{Message: "Success", ID: "id"}
	t.UtilTestMeta(d, d.Message, d.ID)
}

func (t *DriversTestSuite) TestMailGun_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &mailGun{cfg: Comfig, client: m}
	}, false)
}
