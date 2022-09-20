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
)

// postal represents the entity for sending mail via the
// Postal API.
//
// See:
// https://docs.postalserver.io/developer/api
// https://apiv1.postalserver.io/controllers/send/message.html
type postal struct {
	cfg    mail.Config
	client client.Requester
}

const (
	// postalEndpoint defines the endpoint to POST to.
	postalEndpoint = "%s/api/v1/send/message"
	// postalErrorMessage defines the message when an error occurred
	// when sending mail via the Postal API.
	postalErrorMessage = "error sending transmission to Postal API"
)

// NewPostal creates a new Postal client. Configuration
// is validated before initialisation.
func NewPostal(cfg mail.Config) (mail.Mailer, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &postal{
		cfg:    cfg,
		client: client.New(cfg.Client),
	}, nil
}

type (
	// postalTransmission defines the data to be sent to the Postal API.
	postalTransmission struct {
		To          []string           `json:"to"`
		CC          []string           `json:"cc"`
		BCC         []string           `json:"bcc"`
		From        string             `json:"from"`
		Sender      string             `json:"sender"`
		Subject     string             `json:"subject"`
		HTML        string             `json:"html_body"`
		PlainText   string             `json:"plain_body"`
		Attachments []postalAttachment `json:"attachments"`
		Headers     map[string]string  `json:"headers"`
	}
	// postalAttachment defines a singular Postal mail attachment.
	postalAttachment struct {
		Name        string `json:"name"`
		ContentType string `json:"content_type"`
		Data        string `json:"data"`
	}
	// postalResponse defines the data sent back from the Postal API.
	// Status can either be "success" or "error" and data is
	// dynamic dependent on if an error occurred during processing.
	//
	// Example JSON Responses:
	// {"status":"success","time":0.08,"flags":{},"data":{"message_id":"080c21de-52f9-4be1-9cbe-19d63450949c@rp.postal.example.com","messages":{"info@ainsleyclark.com":{"id":28,"token":"WEjrFfpnynRm"}}}}
	// {"status":"error","time":0.0,"flags":{},"data":{"code":"NoRecipients","message":"There are no recipients defined to receive this message"}}
	postalResponse struct {
		Status string                 `json:"status"`
		Time   float32                `json:"time"`
		Flags  map[string]interface{} `json:"flags"`
		Data   map[string]interface{} `json:"data"`
	}
)

func (r *postalResponse) Unmarshal(buf []byte) error {
	resp := &postalResponse{}
	err := json.Unmarshal(buf, resp)
	if err != nil {
		return err
	}
	*r = *resp
	return nil
}

func (r *postalResponse) CheckError(response *http.Response, buf []byte) error {
	if r.Status == "success" {
		return nil
	}
	if len(buf) == 0 {
		return mail.ErrEmptyBody
	}
	msg := postalErrorMessage
	if code, ok := r.Data["code"]; ok {
		msg = fmt.Sprintf("%s - code: %s", msg, code)
	}
	if message, ok := r.Data["message"]; ok {
		msg = fmt.Sprintf("%s, message: %s", msg, message)
	}
	return errors.New(msg)
}

func (r *postalResponse) Meta() httputil.Meta {
	m := httputil.Meta{
		Message: "Successfully sent Postal email",
	}
	if val, ok := r.Data["message_id"]; ok {
		m.ID = fmt.Sprintf("%v", val)
	}
	return m
}

func (d *postal) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	tx := postalTransmission{
		To:        t.Recipients,
		CC:        t.CC,
		BCC:       t.BCC,
		From:      d.cfg.FromAddress,
		Sender:    d.cfg.FromName,
		Subject:   t.Subject,
		HTML:      t.HTML,
		PlainText: t.PlainText,
	}

	if t.HasAttachments() {
		for _, v := range t.Attachments {
			tx.Attachments = append(tx.Attachments, postalAttachment{
				Name:        v.Filename,
				ContentType: v.Mime(),
				Data:        v.B64(),
			})
		}
	}

	tx.Headers = t.Headers

	pl, err := newJSONData(tx)
	if err != nil {
		return mail.Response{}, err
	}

	req := httputil.NewHTTPRequest(http.MethodPost, fmt.Sprintf(postalEndpoint, d.cfg.URL))
	req.AddHeader("X-Server-API-Key", d.cfg.APIKey)

	return d.client.Do(context.Background(), req, pl, &postalResponse{})
}
