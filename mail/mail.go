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

import "errors"

var (
	// Set true to write the HTTP requests in curl for to stdout
	Debug = false

	ErrEmptyBody = errors.New("error, empty body")
)

// Mailer defines the sender for go-mail returning a
// Response or error when an email is sent.
type Mailer interface {
	Send(t *Transmission) (Response, error)
}

const (
	// SparkPost driver type.
	SparkPost = "sparkpost"
	// MailGun driver type.
	MailGun = "mailgun"
	// SendGrid driver type.
	SendGrid = "sendgrid"
	// Postmark driver type.
	Postmark = "postmark"
	// Postal driver type.
	Postal = "postal"
	// SMTP driver type.
	SMTP = "smtp"
)
