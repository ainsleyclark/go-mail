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

package drivers

import (
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	mailsg "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// sendGrid represents the data for sending mail via the
// SendGrid API. mail.Configuration, the client and the
// main send function are parsed for sending
// data.
type sendGrid struct {
	cfg    mail.Config
	client *sendgrid.Client
	send   sendGridSendFunc
}

// sendGridSendFunc defines the function for ending
// SendGrid mail.Transmissions.
type sendGridSendFunc func(email *mailsg.SGMailV3) (*rest.Response, error)

// NewSendGrid creates a new sendGrid client. Configuration
// is validated before initialisation.
func NewSendGrid(cfg mail.Config) (mail.Mailer, error) {
	client := sendgrid.NewSendClient(cfg.APIKey)
	return &sendGrid{
		cfg:    cfg,
		client: client,
		send:   client.Send,
	}, nil
}

// Send posts the go mail mail.Transmission to the SendGrid
// API. mail.Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (m *sendGrid) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	sender := mailsg.NewV3Mail()

	// Add from
	from := mailsg.NewEmail(m.cfg.FromName, m.cfg.FromAddress)
	sender.SetFrom(from)

	// Add subject
	sender.Subject = t.Subject

	// Add to
	p := mailsg.NewPersonalization()
	var to []*mailsg.Email
	for _, recipient := range t.Recipients {
		to = append(to, mailsg.NewEmail("", recipient))
	}
	p.AddTos(to...)

	// Add CC
	if t.HasCC() {
		var cc []*mailsg.Email
		for _, v := range t.CC {
			cc = append(cc, mailsg.NewEmail("", v))
		}
		p.AddCCs(cc...)
	}

	// Add BCC
	if t.HasBCC() {
		var bcc []*mailsg.Email
		for _, v := range t.BCC {
			bcc = append(bcc, mailsg.NewEmail("", v))
		}
		p.AddBCCs(bcc...)
	}

	// Add Plain Text
	if t.PlainText != "" {
		content := mailsg.NewContent("text/plain", t.PlainText)
		sender.AddContent(content)
	}

	// Add HTML
	html := mailsg.NewContent("text/html", t.HTML)
	sender.AddContent(html)

	// Add attachments if they exist.
	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			a := mailsg.NewAttachment()
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
		return mail.Response{}, err
	}

	return mail.Response{
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Headers:    response.Headers,
		ID:         "",
		Message:    "",
	}, nil
}
