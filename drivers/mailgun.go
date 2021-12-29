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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/client"
	"github.com/ainsleyclark/go-mail/internal/httputil"
	"github.com/ainsleyclark/go-mail/mail"
	"net/http"
	"strings"
)

// mailGun represents the data for sending mail via the
// MailGun API. Configuration, the client and the
// main send function are parsed for sending
// data.
type mailGun struct {
	cfg    mail.Config
	client client.Requester
}

const (
	// mailgunEndpoint defines the endpoint to POST to.
	mailgunEndpoint = "/v3/%s/messages"
)

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
		client: client.NewClient(),
	}, nil
}

type (
	// mailGunResponse defines the data sent back from the MailGun API.
	// ID is included on successful transmission.
	//
	// Example JSON Response:
	// {"message":"Need at least one of 'text' or 'html' parameters specified"}
	// {"message":"from parameter is missing"}
	// {"id":"<20211229082318.a988bed7abe472bd@sandboxa6807a568a404524b2b216817d7ed775.mailgun.org>","message":"Queued. Thank you."}
	mailgunResponse struct {
		Message string `json:"message"`
		ID      string `json:"id,omitempty"`
	}
)

func (m *mailgunResponse) HasError(response *http.Response) bool {
	return !client.Is2XX(response.StatusCode)
}

func (m *mailgunResponse) Error() error {
	return errors.New(m.Message)
}

func (m *mailgunResponse) Unmarshal(buf []byte) error {
	resp := &mailgunResponse{}
	err := json.Unmarshal(buf, resp)
	if err != nil {
		return err
	}
	*m = *resp
	return nil
}

func (m *mailgunResponse) Meta() httputil.Meta {
	return httputil.Meta{
		Message: m.Message,
		ID:      m.ID,
	}
}

// Send posts the Go Mail Transmission to the MailGun
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (m *mailGun) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	f := httputil.NewFormData()
	f.AddValue("from", fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.FromAddress))
	f.AddValue("subject", t.Subject)
	f.AddValue("html", t.HTML)
	f.AddValue("text", t.PlainText)

	for _, to := range t.Recipients {
		f.AddValue("to", to)
	}

	if t.HasCC() {
		for _, c := range t.CC {
			f.AddValue("cc", c)
		}
	}

	if t.HasBCC() {
		for _, b := range t.BCC {
			f.AddValue("bcc", b)
		}
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			f.AddBuffer("attachment", v.Filename, v.Bytes)
		}
	}

	url := fmt.Sprintf("%s/%s", m.cfg.URL, strings.TrimPrefix(fmt.Sprintf(mailgunEndpoint, m.cfg.Domain), "/"))
	req := httputil.NewHTTPRequest(http.MethodPost, url)
	req.SetBasicAuth("api", m.cfg.APIKey)

	return m.client.Do(context.Background(), req, f, &mailgunResponse{})
}
