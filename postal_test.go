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
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
)

func (t *MailTestSuite) TestPostal_Send() {
	trans := Transmission{
		Recipients: []string{"recipient@test.com"},
		Subject:    "Subject",
		HTML:       "<h1>HTML</h1>",
		PlainText:  "PlainText",
	}

	tt := map[string]struct {
		input *Transmission
		send  mailGunSendFunc
		want  interface{}
	}{
		"Success": {
			&trans,
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
	}

	for name, test := range tt {
		t.Run(name, func() {
			fmt.Println(test)
			//httptest.NewServer(handler)
			//t.Equal(test.want, resp)
		})
	}
}
