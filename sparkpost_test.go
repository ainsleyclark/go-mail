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
	"errors"
	sp "github.com/SparkPost/gosparkpost"
	"net/http"
)

func (t *MailTestSuite) TestNewSparkPost() {
	tt := map[string]struct {
		input Config
		want  interface{}
	}{
		"Success": {
			Config{
				URL:         "https://api.eu.sparkpost.com",
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
		"Error": {
			Config{
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
			got, err := newSparkPost(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.input, got.cfg)
			t.NotNil(got.client)
			t.NotNil(got.send)
		})
	}
}

func (t *MailTestSuite) TestSparkPost_Send() {
	trans := Transmission{
		Recipients: []string{"recipient@test.com"},
		Subject:    "Subject",
		HTML:       "<h1>HTML</h1>",
		PlainText:  "PlainText",
	}

	tt := map[string]struct {
		input *Transmission
		send  sparkSendFunc
		want  interface{}
	}{
		"Success": {
			&trans,
			func(t *sp.Transmission) (id string, res *sp.Response, err error) {
				return "1", &sp.Response{
					HTTP:    &http.Response{StatusCode: 200, Header: nil},
					Verbose: map[string]string{"msg": "value"},
					Body:    []byte("body"),
				}, nil
			},
			Response{
				StatusCode: 200,
				Body:       "body",
				Headers:    nil,
				ID:         "1",
				Message:    map[string]string{"msg": "value"},
			},
		},
		"With Attachment": {
			&Transmission{
				Recipients: []string{"recipient@test.com"},
				Subject:    "Subject",
				HTML:       "<h1>HTML</h1>",
				PlainText:  "PlainText",
				Attachments: Attachments{
					Attachment{
						Filename: "test.jpg",
					},
				},
			},
			func(t *sp.Transmission) (id string, res *sp.Response, err error) {
				return "1", &sp.Response{
					HTTP:    &http.Response{StatusCode: 200, Header: nil},
					Verbose: map[string]string{"msg": "value"},
					Body:    []byte("body"),
				}, nil
			},
			Response{
				StatusCode: 200,
				Body:       "body",
				Headers:    nil,
				ID:         "1",
				Message:    map[string]string{"msg": "value"},
			},
		},
		"Validation Failed": {
			nil,
			func(t *sp.Transmission) (id string, res *sp.Response, err error) {
				return "", nil, nil
			},
			"can't validate a nil transmission",
		},
		"Send Error": {
			&trans,
			func(t *sp.Transmission) (id string, res *sp.Response, err error) {
				return "", nil, errors.New("send error")
			},
			"send error",
		},
		"Response Error": {
			&trans,
			func(t *sp.Transmission) (id string, res *sp.Response, err error) {
				return "0", &sp.Response{
					Errors: sp.SPErrors{
						sp.SPError{Message: "resp error"},
					},
				}, nil
			},
			"resp error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			spark := sparkPost{
				cfg: Config{
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
