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
	"fmt"
	"log"

	"github.com/flightaware/go-mail/drivers"
	"github.com/flightaware/go-mail/mail"
)

// SendGrid example for Go Mail
func SendGrid() {
	cfg := mail.Config{
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
	}

	mailer, err := drivers.NewSendGrid(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	tx := &mail.Transmission{
		Recipients: []string{"hello@gophers.com"},
		Subject:    "My email",
		HTML:       "<h1>Hello from Go Mail!</h1>",
		PlainText:  "plain text",
	}

	result, err := mailer.Send(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", result)
}
