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
	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
	"log"
)

// SMTP example for Go Mail
func SMTP() {
	cfg := mail.Config{
		URL:         "smtp.gmail.com",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
		Password:    "my-password",
		Port:        587,
	}

	mailer, err := drivers.NewSMTP(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	tx := &mail.Transmission{
		Recipients: []string{"hello@gophers.com"},
		CC:         []string{"cc@gophers.com"},
		BCC:        []string{"bcc@gophers.com"},
		Subject:    "My email",
		HTML:       "<h1>Hello from Go Mail!</h1>",
		PlainText:  "Hello from Go Mail!",
	}

	result, err := mailer.Send(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", result)
}
