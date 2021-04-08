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
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// sendGrid represents the data for sending mail via the
// SendGrid API. Configuration, the client and the
// main send function are parsed for sending
// data.
type sendGrid struct {
	cfg    Config
	client *sendgrid.Client
	send   sendGridSendFunc
}

// sendGridSendFunc defines the function for ending
// SendGrid transmissions.
type sendGridSendFunc func(email *mail.SGMailV3) (*rest.Response, error)

// Creates a new sendGrid client. Configuration is
// validated before initialisation.
func newSendGrid(cfg Config) *sendGrid {
	client := sendgrid.NewSendClient(cfg.APIKey)
	return &sendGrid{
		cfg:    cfg,
		client: client,
		send:   client.Send,
	}
}

// Send posts the go mail Transmission to the SendGrid
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (m *sendGrid) Send(t *Transmission) (Response, error) {
	err := t.Validate()
	if err != nil {
		return Response{}, err
	}

	sender := mail.NewV3Mail()

	// Add from
	from := mail.NewEmail(m.cfg.FromName, m.cfg.FromAddress)
	sender.SetFrom(from)

	// Add subject
	sender.Subject = t.Subject

	// Add to
	p := mail.NewPersonalization()
	var to []*mail.Email
	for _, recipient := range t.Recipients {
		to = append(to, mail.NewEmail("", recipient))
	}
	p.AddTos(to...)

	// Add Plain Text
	if t.PlainText != "" {
		content := mail.NewContent("text/plain", t.PlainText)
		sender.AddContent(content)
	}

	// Add HTML
	html := mail.NewContent("text/html", t.HTML)
	sender.AddContent(html)

	// Add attachments if they exist.
	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			a := mail.NewAttachment()
			a.SetContent(v.B64())
			a.SetType(v.Mime())
			a.SetFilename(v.Filename)
			a.SetDisposition("attachment")
			sender.AddAttachment(a)
		}
	}

	sender.AddPersonalizations(p)

	response, err := m.send(sender)
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Headers:    response.Headers,
		ID:         "",
		Message:    "",
	}, nil
}
