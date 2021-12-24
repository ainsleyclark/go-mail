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

func (t *MailTestSuite) TestAttachments_Exists() {
	tt := map[string]struct {
		input Attachments
		want  bool
	}{
		"Exists": {
			Attachments{
				Attachment{
					Filename: PNGName,
				},
			},
			true,
		},
		"Nil": {
			nil,
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.Exists()
			t.Equal(test.want, got)
		})
	}
}

func (t *MailTestSuite) TestAttachment_Mime() {
	tt := map[string]struct {
		input string
		want  string
	}{
		"PNG": {
			PNGName,
			"image/png",
		},
		"JPG": {
			JPGName,
			"image/jpeg",
		},
		"SVG": {
			SVGName,
			"image/svg+xml",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			a := t.Attachment(test.input)
			got := a.Mime()
			t.Equal(test.want, got)
		})
	}
}

func (t *MailTestSuite) TestAttachment_B64() {
	a := Attachment{
		Bytes: []byte("hello"),
	}
	got := a.B64()
	t.Equal("aGVsbG8=", got)
}
