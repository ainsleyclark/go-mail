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
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// MailTestSuite defines the helper used for mail
// testing.
type MailTestSuite struct {
	suite.Suite
}

// Assert testing has begun.
func TestMail(t *testing.T) {
	suite.Run(t, new(MailTestSuite))
}

const (
	// DataPath defines where the test data resides.
	DataPath = "testdata"
	// PNGName defines the PNG name for testing.
	PNGName = "gopher.png"
)

// Returns a dummy transition for testing with an
// attachment.
func (t *MailTestSuite) GetTransmission() *mail.Transmission {
	wd, err := os.Getwd()
	t.NoError(err)

	err = godotenv.Load(filepath.Join(filepath.Dir(wd), "/.env"))
	if err != nil {
		t.FailNow("Error loading .env file")
	}

	//path := filepath.Join(filepath.Dir(wd), DataPath, PNGName)
	//file, err := ioutil.ReadFile(path)
	//if err != nil {
	//	t.FailNow("Error getting attachment with the path: "+path, err)
	//}

	return &mail.Transmission{
		Recipients: strings.Split(os.Getenv("EMAIL_TO"), ","),
		//CC:         strings.Split(os.Getenv("EMAIL_CC"), ","),
		//BCC:        strings.Split(os.Getenv("EMAIL_BCC"), ","),
		Subject:   "Test - Go Mail",
		HTML:      "<h1>Hello from Go Mail!</h1>",
		PlainText: "Hello from Go Mail!",
		//Attachments: mail.Attachments{
		//	mail.Attachment{
		//		Filename: "gopher.png",
		//		Bytes:    file,
		//	},
		//},
	}
}
