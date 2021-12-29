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
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	mocks "github.com/ainsleyclark/go-mail/mocks/client"
	"net/http"
)

func (t *DriversTestSuite) TestNewSendgrid() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
				URL:         "https://Sendgrid.example.com",
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
			got, err := NewSendGrid(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestSendgridResponse_Unmarshal() {
	t.UtilTestUnmarshal(&sgResponse{}, []byte(`{"errors": []}`))
	res := sgResponse{}
	err := res.Unmarshal(nil)
	t.NoError(err)
}

func (t *DriversTestSuite) TestSendgridResponse_CheckError() {
	tt := map[string]struct {
		input    sgResponse
		response *http.Response
		buf      []byte
		want     error
	}{
		"Success": {
			sgResponse{Errors: nil},
			&http.Response{StatusCode: http.StatusOK},
			[]byte("test"),
			nil,
		},
		"No Errors": {
			sgResponse{},
			&http.Response{StatusCode: http.StatusInternalServerError},
			nil,
			nil,
		},
		"Error": {
			sgResponse{Errors: []sgError{{Message: "message", Field: "field", Help: "help"}}},
			&http.Response{StatusCode: http.StatusInternalServerError},
			[]byte("test"),
			fmt.Errorf("%s - message: message, field: field, help: help", sendgridErrorMessage),
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

func (t *DriversTestSuite) TestSendgridResponse_Meta() {
	d := &sgResponse{}
	t.UtilTestMeta(d, "Successfully sent Sendgrid email", "")
}

func (t *DriversTestSuite) TestSendgrid_Send() {
	t.UtilTestSend(func(m *mocks.Requester) mail.Mailer {
		return &sendGrid{cfg: Comfig, client: m}
	}, true)
}
