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
	"context"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/mailgun/mailgun-go/v4"
	"net/http"
	"time"
)

// mailGun represents the data for sending mail via the
// MailGun API. Configuration, the client and the
// main send function are parsed for sending
// data.
type mailGun struct {
	cfg    mail.Config
	client *mailgun.MailgunImpl
	send   mailGunSendFunc
}

// mailGunSendFunc defines the function for ending MailGun
// transmissions.
type mailGunSendFunc func(ctx context.Context, message *mailgun.Message) (mes string, id string, err error)

// NewMailGun creates a new MailGun client. Configuration
// is validated before initialisation.
func NewMailGun(cfg mail.Config) (mail.Mailer, error) {
	client := mailgun.NewMailgun(cfg.Domain, cfg.APIKey)
	return &mailGun{
		cfg:    cfg,
		client: client,
		send:   client.Send,
	}, nil
}

// Send posts the go mail Transmission to the MailGun
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (m *mailGun) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	message := m.client.NewMessage(m.cfg.FromAddress, t.Subject, t.PlainText, t.Recipients...)
	message.SetHtml(t.HTML)

	if t.HasCC() {
		for _, v := range t.CC {
			message.AddCC(v)
		}
	}

	if t.HasBCC() {
		for _, v := range t.BCC {
			message.AddBCC(v)
		}
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			message.AddBufferAttachment(v.Filename, v.Bytes)
		}
	}

	const Timeout = 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*Timeout)
	defer cancel()

	msg, id, err := m.send(ctx, message)
	if err != nil {
		return mail.Response{}, err
	}

	return mail.Response{
		StatusCode: http.StatusOK,
		Body:       "",
		Headers:    nil,
		ID:         id,
		Message:    msg,
	}, nil
}