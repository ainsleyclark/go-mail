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

func (t *MailTestSuite) TestTransmission_Validate() {
	tt := map[string]struct {
		input *Transmission
		want  error
	}{
		"Success": {
			&Transmission{
				Recipients: []string{"hello@test.com"},
				Subject:    "subject",
				HTML:       "<h1>Hello</h1>",
			},
			nil,
		},
		"Nil": {
			nil,
			errors.New("can't validate a nil transmission"),
		},
		"No Recipients": {
			&Transmission{
				HTML:    "html",
				Subject: "subject",
			},
			errors.New("transmission requires recipients"),
		},
		"No Subject": {
			&Transmission{
				Recipients: []string{"hello@test.com"},
				HTML:       "html",
			},
			errors.New("transmission requires a subject"),
		},
		"No HTML": {
			&Transmission{
				Recipients: []string{"hello@test.com"},
				Subject:    "subject",
			},
			errors.New("transmission requires html content"),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.Validate()
			t.Equal(test.want, got)
		})
	}
}

func (t *MailTestSuite) TestTransmission_HasAttachments() {
	tt := map[string]struct {
		input Transmission
		want  bool
	}{
		"Exists": {
			Transmission{
				Attachments: []Attachment{{Filename: PNGName}},
			},
			true,
		},
		"Nil": {
			Transmission{},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.HasAttachments()
			t.Equal(test.want, got)
		})
	}
}
