// Copyright 2022 Ainsley Clark. All rights reserved.
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
	"errors"
	"fmt"
)

func ExampleTransmission_Validate() {
	t := Transmission{}
	fmt.Println(t.Validate())
	// Output: transmission requires recipients
}

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

func ExampleTransmission_HasCC() {
	t := Transmission{
		CC: []string{"cc@gophers.com"},
	}
	fmt.Println(t.HasCC())
	// Output: true
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

func ExampleTransmission_HasBCC() {
	t := Transmission{
		BCC: []string{"bcc@gophers.com"},
	}
	fmt.Println(t.HasBCC())
	// Output: true
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

func ExampleTransmission_HasAttachments() {
	t := Transmission{
		Attachments: []Attachment{
			{
				Filename: "gopher.svg",
				Bytes:    []byte("svg"),
			},
		},
	}
	fmt.Println(t.HasAttachments())
	// Output: true
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
