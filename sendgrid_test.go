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
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

func (t *MailTestSuite) TestSendGrid_Send() {
	trans := Transmission{
		Recipients: []string{"recipient@test.com"},
		Subject:    "Subject",
		HTML:       "<h1>HTML</h1>",
		PlainText:  "PlainText",
	}

	tt := map[string]struct {
		input *Transmission
		send  sendGridSendFunc
		want  interface{}
	}{
		"Success": {
			&trans,
			func(email *mail.SGMailV3) (*rest.Response, error) {
				return &rest.Response{
					StatusCode: http.StatusOK,
					Body:       "body",
					Headers:    map[string][]string{"msg": {"test"}},
				}, nil
			},
			Response{
				StatusCode: http.StatusOK,
				Body:       "body",
				Headers:    map[string][]string{"msg": {"test"}},
				ID:         "",
				Message:    "",
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
			func(email *mail.SGMailV3) (*rest.Response, error) {
				return &rest.Response{
					StatusCode: http.StatusOK,
					Body:       "body",
					Headers:    map[string][]string{"msg": {"test"}},
				}, nil
			},
			Response{
				StatusCode: http.StatusOK,
				Body:       "body",
				Headers:    map[string][]string{"msg": {"test"}},
				ID:         "",
				Message:    "",
			},
		},
		"Validation Failed": {
			nil,
			func(email *mail.SGMailV3) (*rest.Response, error) {
				return nil, nil
			},
			"can't validate a nil transmission",
		},
		"Send Error": {
			&trans,
			func(email *mail.SGMailV3) (*rest.Response, error) {
				return nil, errors.New("send error")
			},
			"send error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			spark := sendGrid{
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
