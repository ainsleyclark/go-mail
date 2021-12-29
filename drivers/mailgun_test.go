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
	"github.com/ainsleyclark/go-mail/mail"
	mocks "github.com/ainsleyclark/go-mail/mocks/client"
)

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
			got, err := NewMailGun(test.input)
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
	d := &mailgunResponse{Message: "Error"}
	t.UtilTestCheckError(d, d.Message, true)
}

func (t *DriversTestSuite) TestMailgunResponse_Meta() {
	d := &mailgunResponse{Message: "Success", ID: "id"}
	t.UtilTestMeta(d, d.Message, d.ID)
}

func (t *DriversTestSuite) TestMailGun_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &mailGun{cfg: Comfig, client: m}
	})
}
