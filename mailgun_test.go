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
	"context"
	"errors"
	"github.com/mailgun/mailgun-go/v4"
)

func (t *MailTestSuite) TestMailGun_Send() {
	tt := map[string]struct {
		input *Transmission
		send  mailGunSendFunc
		want  interface{}
	}{
		"Success": {
			Trans,
			func(ctx context.Context, message *mailgun.Message) (mes string, id string, err error) {
				return "success", "1", nil
			},
			Response{
				StatusCode: 200,
				Body:       "",
				Headers:    nil,
				ID:         "1",
				Message:    "success",
			},
		},
		"With Attachment": {
			TransWithAttachment,
			func(ctx context.Context, message *mailgun.Message) (mes string, id string, err error) {
				return "success", "1", nil
			},
			Response{
				StatusCode: 200,
				Body:       "",
				Headers:    nil,
				ID:         "1",
				Message:    "success",
			},
		},
		"Validation Failed": {
			nil,
			func(ctx context.Context, message *mailgun.Message) (mes string, id string, err error) {
				return "", "", nil
			},
			"can't validate a nil transmission",
		},
		"Send Error": {
			Trans,
			func(ctx context.Context, message *mailgun.Message) (mes string, id string, err error) {
				return "", "", errors.New("send error")
			},
			"send error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			spark := mailGun{
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
