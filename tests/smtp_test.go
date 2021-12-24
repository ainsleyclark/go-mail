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
	"github.com/ainsleyclark/go-mail"
	"os"
	"strconv"
)

func (t *MailTestSuite) Test_SMTP() {
	tx := t.GetTransmission()

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		t.FailNow("Error parsing SMTP port")
	}

	driver, err := mail.NewClient(mail.SMTP, mail.Config{
		URL:         os.Getenv("SMTP_URL"),
		FromAddress: os.Getenv("SMTP_FROM_ADDRESS"),
		FromName:    os.Getenv("SMTP_FROM_NAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Port:        port,
	})
	if err != nil {
		t.FailNow("Error creating client: " + err.Error())
	}

	result, err := driver.Send(tx)
	if err != nil {
		t.FailNow("Error sending SMTP email: " + err.Error())
	}

	t.Equal(200, result.StatusCode)
	t.NotEmpty(result.Message)
}
