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
	t.UtilTestCheckError_Error(&postmarkResponse{ErrorCode: 10}, postmarkErrorMessage, true)
	t.UtilTestCheckError_Success(&postmarkResponse{})
}

func (t *DriversTestSuite) TestPostmarkResponse_Meta() {
	d := &postmarkResponse{Message: "Success", ID: "id"}
	t.UtilTestMeta(d, d.Message, d.ID)
}

func (t *DriversTestSuite) TestPostmark_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &postmark{cfg: Comfig, client: m}
	})
}
