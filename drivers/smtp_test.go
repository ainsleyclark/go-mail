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
	"fmt"
	"net/smtp"

	"github.com/flightaware/go-mail/mail"
)

func (t *DriversTestSuite) TestNewSMTP() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
				URL:         "https://smtp.example.com",
				FromAddress: "addr",
				FromName:    "name",
				Password:    "password",
			},
			nil,
		},
		"No url": {
			mail.Config{},
			"driver requires a url",
		},
		"No From Address": {
			mail.Config{
				URL: "https://smtp.example.com",
			},
			"driver requires from address",
		},
		"No From Name": {
			mail.Config{
				URL:         "https://smtp.example.com",
				FromAddress: "hello@gophers.com",
			},
			"driver requires from name",
		},
		"No Password": {
			mail.Config{
				URL:         "https://smtp.example.com",
				FromAddress: "hello@gophers.com",
				FromName:    "name",
			},
			"driver requires a password",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := NewSMTP(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestSMTP_Send() {
	tt := map[string]struct {
		input *mail.Transmission
		send  smtpSendFunc
		want  interface{}
	}{
		"Success": {
			Trans,
			func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
				return nil
			},
			mail.Response{
				StatusCode: 200,
				Message:    "Email sent successfully",
			},
		},
		"With Attachment": {
			TransWithAttachment,
			func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
				return nil
			},
			mail.Response{
				StatusCode: 200,
				Message:    "Email sent successfully",
			},
		},
		"Validation Failed": {
			nil,
			func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
				return nil
			},
			"can't validate a nil transmission",
		},
		"Send Error": {
			Trans,
			func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
				return errors.New("send error")
			},
			"send error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			spark := smtpClient{
				cfg: mail.Config{
					FromAddress: "from",
				},
				send: test.send,
			}
			resp, err := spark.Send(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, resp)
		})
	}
}

func (t *DriversTestSuite) TestSMTP_Bytes() {
	t.T().Skip()

	m := smtpClient{}
	got := m.bytes(&mail.Transmission{
		Recipients: []string{"hello@gmail.com"},
		Subject:    "Subject",
		HTML:       "<h1>Hey!</h1>",
		PlainText:  "Hey!",
		//Attachments: []mail.Attachment{
		//	{
		//		Filename: "test.jpg",
		//		Bytes:    []byte("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg=="),
		//	},
		//},
	})
	fmt.Println(string(got))
}
