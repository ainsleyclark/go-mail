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
	"errors"
	"fmt"
	"github.com/ainsleyclark/go-mail/mail"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
)

// smtpClient represents the data for sending mail via
// plain ol SMTP. Configuration, the client and the
// main send function are parsed for sending
// data.
type smtpClient struct {
	cfg  mail.Config
	send smtpSendFunc
}

// smtpSendFunc defines the function for ending
// SMTP mail.Transmissions.
type smtpSendFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// NewSMTP creates a new smtp client. Configuration
// is validated before initialisation.
func NewSMTP(cfg mail.Config) (mail.Mailer, error) {
	if cfg.URL == "" {
		return nil, errors.New("driver requires a url")
	}
	if cfg.FromAddress == "" {
		return nil, errors.New("driver requires from address")
	}
	if cfg.FromName == "" {
		return nil, errors.New("driver requires from name")
	}
	if cfg.Password == "" {
		return nil, errors.New("driver requires a password")
	}
	return &smtpClient{
		cfg:  cfg,
		send: smtp.SendMail,
	}, nil
}

// Send mail via plain SMTP. mail.Transmissions are validated
// before sending and attachments are added. Returns
// an error upon failure.
func (m *smtpClient) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	auth := smtp.PlainAuth("", m.cfg.FromAddress, m.cfg.Password, m.cfg.URL)
	err = m.send(m.cfg.URL+":"+strconv.Itoa(m.cfg.Port), auth, m.cfg.FromAddress, m.getTo(t), m.bytes(t))
	if err != nil {
		return mail.Response{}, err
	}

	return mail.Response{
		StatusCode: http.StatusOK,
		Message:    "Email sent successfully",
	}, nil
}

// getTo returns the merged mail.Transmission recipients, CC and
// BCC email addresses.
func (m *smtpClient) getTo(t *mail.Transmission) []string {
	var to []string
	to = append(t.Recipients, t.CC...)
	to = append(to, t.BCC...)
	return to
}

// Processes the mail.Transmission and returns the bytes for
// sending. Mime types are set dependent on the
// content passed.
func (m *smtpClient) bytes(t *mail.Transmission) []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("Subject: %s\n", t.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(t.Recipients, ",")))

	if t.HasCC() {
		buf.WriteString(fmt.Sprintf("CC: %s\n", strings.Join(t.CC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if t.HasAttachments() {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/html; charset=\"ascii\"\n")
	}

	if t.PlainText != "" {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
		buf.WriteString(t.PlainText)
	}

	if t.HTML != "" {
		buf.WriteString("Content-Type: text/html; charset=\"ascii\"\n")
		buf.WriteString(t.HTML)
	}

	if t.HasAttachments() {
		for _, v := range t.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", v.Mime()))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", v.Filename))
			buf.WriteString(fmt.Sprintf("\n--%s", v.B64()))
		}
		buf.WriteString("--")
	}

	return buf.Bytes()
}
