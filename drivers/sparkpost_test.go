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
	SparkpostHeaders = http.Header{"Authorization": []string{""}}
)

func (t *DriversTestSuite) TestNewSparkPost() {
	tt := map[string]struct {
		input mail.Config
		want  interface{}
	}{
		"Success": {
			mail.Config{
				URL:         "https://api.eu.sparkpost.com",
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
		"Error": {
			mail.Config{
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
			got, err := NewSparkPost(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.NotNil(got)
		})
	}
}

func (t *DriversTestSuite) TestSparkpost_Send() {
	tt := map[string]struct {
		input *mail.Transmission
		mock  func(m *mocks.Requester)
		want  interface{}
	}{
		"Success": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, sparkpostEndpoint, SparkpostHeaders).
					Return([]byte(`{"results":{"total_rejected_recipients":0,"total_accepted_recipients":1,"id":"1"}}`), &http.Response{StatusCode: http.StatusOK}, nil)
			},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"results":{"total_rejected_recipients":0,"total_accepted_recipients":1,"id":"1"}}`,
				Message:    "Successfully sent Sparkpost email",
				ID:         "1",
			},
		},
		"With Attachment": {
			TransWithAttachment,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, sparkpostEndpoint, SparkpostHeaders).
					Return([]byte(`{"results":{"total_rejected_recipients":0,"total_accepted_recipients":1,"id":"1"}}`), &http.Response{StatusCode: http.StatusOK}, nil)
			},
			mail.Response{
				StatusCode: http.StatusOK,
				Body:       `{"results":{"total_rejected_recipients":0,"total_accepted_recipients":1,"id":"1"}}`,
				Message:    "Successfully sent Sparkpost email",
				ID:         "1",
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
				m.On("Do", mock.Anything, sparkpostEndpoint, SparkpostHeaders).
					Return([]byte("output"), nil, errors.New("do error"))
			},
			"do error",
		},
		"Unmarshal Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, sparkpostEndpoint, SparkpostHeaders).
					Return([]byte(`wrong`), nil, nil)
			},
			"invalid character",
		},
		"Response Error": {
			Trans,
			func(m *mocks.Requester) {
				m.On("Do", mock.Anything, sparkpostEndpoint, SparkpostHeaders).
					Return([]byte(`{"errors": [{"message": "Error", "code": "10"}]}`), nil, nil)
			},
			fmt.Sprintf("%s - code: 10, message: Error", sparkpostErrorMessage),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			m := &mocks.Requester{}
			if test.mock != nil {
				test.mock(m)
			}

			sp := sparkPost{
				cfg:    mail.Config{FromAddress: "from"},
				client: m,
			}

			resp, err := sp.Send(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}

			t.Equal(test.want, resp)
		})
	}
}
