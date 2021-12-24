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

import "errors"

func (t *MailTestSuite) TestConfig_Validate() {
	tt := map[string]struct {
		input Config
		want  error
	}{
		"Success": {
			Config{
				APIKey:      "key",
				FromAddress: "hello@test.com",
				FromName:    "Test",
			},
			nil,
		},
		"No From Address": {
			Config{
				APIKey:   "key",
				FromName: "Test",
			},
			errors.New("mailer requires from address"),
		},
		"No From Name": {
			Config{
				APIKey:      "key",
				FromAddress: "hello@test.com",
			},
			errors.New("mailer requires from name"),
		},
		"No Key": {
			Config{
				FromAddress: "hello@test.com",
				FromName:    "Test",
			},
			errors.New("mailer requires api key"),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.Validate()
			t.Equal(test.want, got)
		})
	}
}

func (t *MailTestSuite) TestConfig_HasCC() {
	tt := map[string]struct {
		input Transmission
		want  bool
	}{
		"With": {
			Transmission{CC: []string{"hello@test.com"}},
			true,
		},
		"Without": {
			Transmission{},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.HasCC()
			t.Equal(test.want, got)
		})
	}
}

func (t *MailTestSuite) TestConfig_HasBCC() {
	tt := map[string]struct {
		input Transmission
		want  bool
	}{
		"With": {
			Transmission{BCC: []string{"hello@test.com"}},
			true,
		},
		"Without": {
			Transmission{},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.HasBCC()
			t.Equal(test.want, got)
		})
	}
}
