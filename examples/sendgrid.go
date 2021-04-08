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
	"github.com/ainsleyclark/go-mail"
)

func SendGrid() {
	cfg := mail.Config{
		APIKey:      "my-key",
		FromAddress: "hello@gophers.com",
		FromName:    "Gopher",
	}

	driver, err := mail.NewClient(mail.SendGrid, cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	tx := &mail.Transmission{
		Recipients: []string{"hello@gophers.com"},
		Subject:    "My email",
		HTML:       "<h1>Hello from go mail!</h1>",
		PlainText:  "plain text",
	}

	result, err := driver.Send(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
