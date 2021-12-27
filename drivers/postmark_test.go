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
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/ainsleyclark/go-mail/mocks"
	"github.com/stretchr/testify/mock"
	"net/http"
)

var (
	PostmarkHeaders = http.Header{"X-Postmark-Server-Token": []string{""}}
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

//func (t *DriversTestSuite) TestPostalResponse_HasError() {
//	tt := map[string]struct {
//		input postalResponse
//		want  bool
//	}{
//		"Error": {
//			postalResponse{Status: "success"},
//			false,
//		},
//		"No Error": {
//			postalResponse{Status: "error"},
//			true,
//		},
//	}
//
//	for name, test := range tt {
//		t.Run(name, func() {
//			got := test.input.HasError()
//			t.Equal(test.want, got)
//		})
//	}
//}
//
//func (t *DriversTestSuite) TestPostalResponse_Error() {
//	tt := map[string]struct {
//		input postalResponse
//		want  string
//	}{
//		"Default": {
//			postalResponse{},
//			postalErrorMessage,
//		},
//		"Code": {
//			postalResponse{Data: map[string]interface{}{"code": "ValidationFailed"}},
//			fmt.Sprintf("%s - code: ValidationFailed", postalErrorMessage),
//		},
//		"All": {
//			postalResponse{Data: map[string]interface{}{"code": "ValidationFailed", "message": "Postal Message"}},
//			fmt.Sprintf("%s - code: ValidationFailed, message: Postal Message", postalErrorMessage),
//		},
//	}
//
//	for name, test := range tt {
//		t.Run(name, func() {
//			got := test.input.Error()
//			t.Contains(got.Error(), test.want)
//		})
//	}
//}

func (t *DriversTestSuite) TestPostmark_Send() {
	tt := map[string]struct {
		input   *mail.Transmission
		handler func(m *mocks.Requester)
		want    interface{}
	}{
		"Success": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, postmarkEndpoint, PostmarkHeaders).
					Return([]byte(`{"To":"hello@gophers.com","MessageID":"10","ErrorCode":0,"Message":"OK"}`), &http.Response{StatusCode: http.StatusOK}, nil)
			},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"To":"hello@gophers.com","MessageID":"10","ErrorCode":0,"Message":"OK"}`,
				Headers:    nil,
				ID:         "10",
				Message:    "OK",
			},
		},
		"With Attachment": {
			TransWithAttachment,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, postmarkEndpoint, PostmarkHeaders).
					Return([]byte(`{"To":"hello@gophers.com","MessageID":"10","ErrorCode":0,"Message":"OK"}`), &http.Response{StatusCode: http.StatusOK}, nil)
			},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"To":"hello@gophers.com","MessageID":"10","ErrorCode":0,"Message":"OK"}`,
				Headers:    nil,
				ID:         "10",
				Message:    "OK",
			},
		},
		"Validation Failed": {
			nil,
			nil,
			"can't validate a nil transmission",
		},
		"Do Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, postmarkEndpoint, PostmarkHeaders).
					Return([]byte("output"), nil, errors.New("do error"))
			},
			"do error",
		},
		"Unmarshal Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, postmarkEndpoint, PostmarkHeaders).
					Return([]byte(`wrong`), nil, nil)
			},
			"invalid character",
		},
		"Response Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, postmarkEndpoint, PostmarkHeaders).
					Return([]byte(`{"ErrorCode": 10, "Message": "Error"}`), nil, nil)
			},
			fmt.Sprintf("%s - code: 10, message: Error", postmarkErrorMessage),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			m := &mocks.Requester{}
			if test.handler != nil {
				test.handler(m)
			}

			ptl := postmark{
				cfg:    mail.Config{FromAddress: "from"},
				client: m,
			}

			resp, err := ptl.Send(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}

			t.Equal(test.want, resp)
		})
	}
}
