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
	"errors"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/sendgrid/rest"
	mailsg "github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

func (t *DriversTestSuite) TestSendGrid_Send() {
	tt := map[string]struct {
		input *mail.Transmission
		send  sendGridSendFunc
		want  interface{}
	}{
		"Success": {
			Trans,
			func(email *mailsg.SGMailV3) (*rest.Response, error) {
				return &rest.Response{
					StatusCode: http.StatusOK,
					Body:       "body",
					Headers:    map[string][]string{"msg": {"test"}},
				}, nil
			},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       "body",
				Headers:    map[string][]string{"msg": {"test"}},
				ID:         "",
				Message:    "",
			},
		},
		"With Attachment": {
			TransWithAttachment,
			func(email *mailsg.SGMailV3) (*rest.Response, error) {
				return &rest.Response{
					StatusCode: http.StatusOK,
					Body:       "body",
					Headers:    map[string][]string{"msg": {"test"}},
				}, nil
			},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       "body",
				Headers:    map[string][]string{"msg": {"test"}},
				ID:         "",
				Message:    "",
			},
		},
		"Validation Failed": {
			nil,
			func(email *mailsg.SGMailV3) (*rest.Response, error) {
				return nil, nil
			},
			"can't validate a nil transmission",
		},
		"Send Error": {
			Trans,
			func(email *mailsg.SGMailV3) (*rest.Response, error) {
				return nil, errors.New("send error")
			},
			"send error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			spark := sendGrid{
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
