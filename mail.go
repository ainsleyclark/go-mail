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
)

// Mailer defines the sender for go-mail returning a
// Response or error when an email is sent.
type Mailer interface {
	Send(t *Transmission) (Response, error)
}

// Response represents the data passed back from a
// successful transmission. Where possible, a
// status code, body, headers will be
// returned within the response.
type Response struct {
	StatusCode int                 // e.g. 200
	Body       string              // e.g. {"result: success"}
	Headers    map[string][]string // e.g. map[X-Ratelimit-Limit:[600]]
	ID         string              // e.g "100"
	Message    interface{}         // e.g "Email sent successfully"
}

const (
	// SparkPost driver type.
	SparkPost = "sparkpost"
	// MailGun driver type.
	MailGun = "mailgun"
	// SendGrid driver type.
	SendGrid = "sendgrid"
	// SMTP driver type.
	SMTP = "smtp"
)

// NewClient
//
// Creates a new Mailer based on the input driver.
// Sparkpost, MailGun or SendGrid can be passed.
// Returns an error if a driver did not match,
// Or there was an error creating the client.
func NewClient(driver string, cfg Config) (Mailer, error) {
	switch driver {
	case SparkPost:
		return newSparkPost(cfg)
	case MailGun:
		return newMailGun(cfg), nil
	case SendGrid:
		return newSendGrid(cfg), nil
	case SMTP:
		return newSMTP(cfg), nil
	}
	return nil, errors.New(driver + " not supported")
}
