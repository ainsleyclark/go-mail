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
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/client"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/mailgun/mailgun-go/v4"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// mailGun represents the data for sending mail via the
// MailGun API. Configuration, the client and the
// main send function are parsed for sending
// data.
type mailGun struct {
	cfg    mail.Config
	client client.Requester
}

// NewMailGun creates a new MailGun client. Configuration
// is validated before initialisation.
func NewMailGun(cfg mail.Config) (mail.Mailer, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	if cfg.Domain == "" {
		return nil, errors.New("driver requires a domain")
	}
	return &mailGun{
		cfg:    cfg,
		client: client.New(cfg.URL),
	}, nil
}

type (
	mailgunTransmission struct {
		// Email address for From header
		From string `json:"from"`
		// Email address of the recipient(s). Example: "Bob <bob@host.com>".
		// You can use commas to separate multiple recipients.
		To string `json:"to"`
		// Same as To but for Cc
		CC string `json:"cc"`
		// Same as To but for Bcc
		BCC string `json:"bcc"`
		// Message subject
		Subject string `json:"subject"`
		//	Body of the message. (text version)
		Text string `json:"text"`
		// Body of the message. (HTML version)
		HTML string `json:"html"`
		// AMP part of the message. Please follow google guidelines to compose and send AMP emails.
		AMPHtml string `json:"amp-html"`
		// attachment	File attachment. You can post multiple attachment values. Important: You must use multipart/form-data encoding when sending attachments.
		// Name of a template stored via template API. See Templates for more information
		Template string `json:"template"`
		// Use this parameter to send a message to specific version of a template
		TemplateVersion string `json:"t:version"`
	}
	// StoredAttachment structures contain information on an attachment associated with a stored message.
	StoredAttachment struct {
		Size        int    `json:"size"`
		Url         string `json:"url"`
		Name        string `json:"name"`
		ContentType string `json:"content-type"`
	}
type BufferAttachment struct {
Filename string
Buffer   []byte
}
)

func (m *mailGun) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	tx := mailgunTransmission{
		From:    fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.FromAddress),
		To:      strings.Join(t.Recipients, ","),
		CC:      strings.Join(t.CC, ","),
		BCC:     strings.Join(t.BCC, ","),
		Subject: t.Subject,
		Text:    t.PlainText,
		HTML:    t.HTML,
	}

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	defer writer.Close()

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			file, err := writer.CreateFormFile("attachment", v.Filename)
			if err != nil {
				return mail.Response{}, err
			}
			io.Copy(file, bytes.NewReader(v.Bytes))
		}
	}

	if tmp, err := writer.CreateFormField("from"); err == nil {
		tmp.Write([]byte(fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.FromAddress)))
	} else {
		return nil, err
	}

	file, err := writer.CreateFormFile("from")
	if err != nil {
		return mail.Response{}, err
	}
	err := writer.CreateFormField("from", fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.FromAddress))
	if err != nil {
		return mail.Response{}, err
	}

	return mail.Response{}, nil
}

// Send posts the go mail Transmission to the MailGun
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (m *mailGun) OLD(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	mailg := mailgun.NewMailgun("fff", "Fff")

	message := mailg.NewMessage(m.cfg.FromAddress, t.Subject, t.PlainText, t.Recipients...)
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

	msg, id, err := mailg.Send(ctx, message)
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
