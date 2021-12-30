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

import "fmt"

func ExampleAttachment_Mime() {
	svg := `
<svg width="100" height="100">
  <circle cx="50" cy="50" r="40" stroke="green" stroke-width="4" fill="yellow" />
</svg>`

	a := Attachment{
		Filename: "circle.svg",
		Bytes:    []byte(svg),
	}

	fmt.Println(a.Mime())
	// Output: image/svg+xml
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

func ExampleAttachment_B64() {
	svg := `
<svg width="100" height="100">
  <circle cx="50" cy="50" r="40" stroke="green" stroke-width="4" fill="yellow" />
</svg>`

	a := Attachment{
		Filename: "circle.svg",
		Bytes:    []byte(svg),
	}

	fmt.Println(a.B64())
	// Output: Cjxzdmcgd2lkdGg9IjEwMCIgaGVpZ2h0PSIxMDAiPgogIDxjaXJjbGUgY3g9IjUwIiBjeT0iNTAiIHI9IjQwIiBzdHJva2U9ImdyZWVuIiBzdHJva2Utd2lkdGg9IjQiIGZpbGw9InllbGxvdyIgLz4KPC9zdmc+
}

func (t *MailTestSuite) TestAttachment_B64() {
	a := Attachment{
		Bytes: []byte("hello"),
	}
	got := a.B64()
	t.Equal("aGVsbG8=", got)
}
