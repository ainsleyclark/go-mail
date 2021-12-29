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
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// DataPath defines where the test data resides.
	DataPath = "testdata"
	// PNGName defines the PNG name for testing.
	PNGName = "gopher.png"
)

// Load the Env variables for testing.
func LoadEnv(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	assert.NoError(t, err)

	path := filepath.Join(filepath.Dir(wd), "/.env")
	err = godotenv.Load(path)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error loading .env file with path: %s, using system defaults.\n", path)
	}
}

// Returns a dummy transition for testing with an
// attachment.
func GetTransmission(t *testing.T) *mail.Transmission {
	t.Helper()

	wd, err := os.Getwd()
	assert.NoError(t, err)

	path := filepath.Join(filepath.Dir(wd), DataPath, PNGName)
	file, err := ioutil.ReadFile(path)
	if err != nil {

		t.Fatal("Error getting attachment with the path: "+path, err)
	}

	return &mail.Transmission{
		Recipients: strings.Split(os.Getenv("EMAIL_TO"), ","),
		//CC:         strings.Split(os.Getenv("EMAIL_CC"), ","),
		//BCC:        strings.Split(os.Getenv("EMAIL_BCC"), ","),
		Subject:   "Test - Go Mail",
		HTML:      "<h1>Hello from Go Mail!</h1>",
		PlainText: "Hello from Go Mail",
		Attachments: []mail.Attachment{
			{
				Filename: "gopher.png",
				Bytes:    file,
			},
		},
	}
}

// UtilTestSend is a helper function for performing live mailing
// tests for the drivers.
func UtilTestSend(t *testing.T, fn func(cfg mail.Config) (mail.Mailer, error), cfg mail.Config, driver string) {
	t.Helper()

	tx := GetTransmission(t)

	mailer, err := fn(cfg)
	if err != nil {
		t.Fatal("Error creating client: " + err.Error())
	}

	result, err := mailer.Send(tx)
	if err != nil {
		t.Fatalf("Error sending %s email: %s", strings.Title(driver), err.Error())
	}

	// Print for sanity
	fmt.Println(string(result.Body))

	assert.InDelta(t, result.StatusCode, http.StatusOK, 299)
	assert.NotEmpty(t, result.Message)
}
