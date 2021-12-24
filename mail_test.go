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

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

// MailTestSuite defines the helper used for mail
// testing.
type MailTestSuite struct {
	suite.Suite
	base string
}

// Assert testing has begun.
func TestMail(t *testing.T) {
	suite.Run(t, new(MailTestSuite))
}

// Assigns test base.
func (t *MailTestSuite) SetupSuite() {
	wd, err := os.Getwd()
	t.NoError(err)
	t.base = wd
}

const (
	// DataPath defines where the test data resides.
	DataPath = "testdata"
	// PNGName defines the PNG name for testing.
	PNGName = "gopher.png"
	// JPGName defines the JPG name for testing.
	JPGName = "gopher.jpg"
	// SVGName defines the SVG name testing.
	SVGName = "gopher.svg"
)

var (
	// Trans is the transmission used for testing.
	Trans = &Transmission{
		Recipients: []string{"recipient@test.com"},
		CC:         []string{"recipient@test.com"},
		BCC:        []string{"recipient@test.com"},
		Subject:    "Subject",
		HTML:       "<h1>HTML</h1>",
		PlainText:  "PlainText",
	}
	// Trans is the transmission with an
	// attachment used for testing.
	TransWithAttachment = &Transmission{
		Recipients: []string{"recipient@test.com"},
		Subject:    "Subject",
		HTML:       "<h1>HTML</h1>",
		PlainText:  "PlainText",
		Attachments: Attachments{
			Attachment{
				Filename: "test.jpg",
			},
		},
	}
)

// Returns a PNG attachment for testing.
func (t *MailTestSuite) Attachment(name string) Attachment {
	path := t.base + string(os.PathSeparator) + DataPath + string(os.PathSeparator) + name
	file, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fail("error getting attachment with the path: "+path, err)
	}

	return Attachment{
		Filename: name,
		Bytes:    file,
	}
}

func (t *MailTestSuite) TestNewClient() {
	cfg := Config{
		APIKey:      "key",
		FromAddress: "hello@test.com",
		FromName:    "Test",
	}

	tt := map[string]struct {
		driver string
		input  Config
		want   error
	}{
		"SparkPost": {
			SparkPost,
			cfg,
			nil,
		},
		"MailGun": {
			MailGun,
			cfg,
			nil,
		},
		"SendGrid": {
			SendGrid,
			cfg,
			nil,
		},
		"Postal": {
			Postal,
			cfg,
			nil,
		},
		"SMTP": {
			SMTP,
			cfg,
			nil,
		},
		"Error": {
			"wrong",
			Config{},
			errors.New("wrong not supported"),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := NewClient(test.driver, test.input)
			if err != nil {
				t.Equal(test.want, err)
				return
			}
			t.NotNil(got)
		})
	}
}

// Docs
func ExampleNewClient() {
	mailer, err := NewClient(SparkPost, Config{
		APIKey:      "my-key",
		FromAddress: "hello@test.com",
		FromName:    "Test",
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(mailer)
}
