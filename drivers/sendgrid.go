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
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/client"
	"github.com/ainsleyclark/go-mail/internal/httputil"
	"github.com/ainsleyclark/go-mail/mail"
	"net/http"
)

// sendGrid represents the entity for sending mail via the
// SendGrid API.
//
// See:
// https://docs.sendgrid.com/api-reference/how-to-use-the-sendgrid-v3-api
// https://docs.sendgrid.com/api-reference/mail-send/mail-send
type sendGrid struct {
	cfg    mail.Config
	client client.Requester
}

const (
	// sendGridEndpoint defines the endpoint to POST to.
	// The host for Web API v3 requests is always https://sendgrid.com/v3/
	sendGridEndpoint = "https://api.sendgrid.com/v3/mail/send"
	// sendgridErrorMessage defines the message when an error occurred
	// when sending mail via the SendGrid API.
	sendgridErrorMessage = "error sending transmission to SendGrid API"
)

// NewSendGrid creates a new sendGrid client. Configuration
// is validated before initialisation.
func NewSendGrid(cfg mail.Config) (mail.Mailer, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &sendGrid{
		cfg:    cfg,
		client: client.New(),
	}, nil
}

type (
	sgTransmission struct {
		From             *sgEmail             `json:"from,omitempty"`
		Subject          string               `json:"subject,omitempty"`
		Personalizations []*sgPersonalization `json:"personalizations,omitempty"`
		Content          []*sgContent         `json:"content,omitempty"`
		Attachments      []*sgAttachment      `json:"attachments,omitempty"`
		TemplateID       string               `json:"template_id,omitempty"`
		Sections         map[string]string    `json:"sections,omitempty"`
		Headers          map[string]string    `json:"headers,omitempty"`
		Categories       []string             `json:"categories,omitempty"`
		CustomArgs       map[string]string    `json:"custom_args,omitempty"`
		SendAt           int                  `json:"send_at,omitempty"`
		BatchID          string               `json:"batch_id,omitempty"`
		IPPoolID         string               `json:"ip_pool_name,omitempty"`
		ReplyTo          *sgEmail             `json:"reply_to,omitempty"`
	}
	sgPersonalization struct {
		To                  []*sgEmail             `json:"to,omitempty"`
		From                *sgEmail               `json:"from,omitempty"`
		CC                  []*sgEmail             `json:"cc,omitempty"`
		BCC                 []*sgEmail             `json:"bcc,omitempty"`
		Subject             string                 `json:"subject,omitempty"`
		Headers             map[string]string      `json:"headers,omitempty"`
		Substitutions       map[string]string      `json:"substitutions,omitempty"`
		CustomArgs          map[string]string      `json:"custom_args,omitempty"`
		DynamicTemplateData map[string]interface{} `json:"dynamic_template_data,omitempty"`
		Categories          []string               `json:"categories,omitempty"`
		SendAt              int                    `json:"send_at,omitempty"`
	}
	// sgEmail holds email name and address info
	sgEmail struct {
		Name    string `json:"name,omitempty"`
		Address string `json:"email,omitempty"`
	}
	// sgContent defines content of the mail body
	sgContent struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}
	// sgAttachment holds attachement information
	sgAttachment struct {
		Content     string `json:"content,omitempty"`
		Type        string `json:"type,omitempty"`
		Name        string `json:"name,omitempty"`
		Filename    string `json:"filename,omitempty"`
		Disposition string `json:"disposition,omitempty"`
		ContentID   string `json:"content_id,omitempty"`
	}
	// sgResponse contains the response data from the SendGrid
	// API.
	// Note: No response data is passed if the response code is 2xx
	//
	// Example JSON Response:
	// {"errors":[{"message":"The from object must be provided for every email send. It is an object that requires the email parameter, but may also contain a name parameter.  e.g. {\"email\" : \"example@example.com\"}  or {\"email\" : \"example@example.com\", \"name\" : \"Example Recipient\"}.","field":"from.email","help":"http://sendgrid.com/docs/API_Reference/Web_API_v3/Mail/errors.html#message.from"}]}
	sgResponse struct {
		Errors []sgError `json:"errors"`
	}
	// sgError defines a singular validation error from the API.
	sgError struct {
		Message string `json:"message"`
		Field   string `json:"field"`
		Help    string `json:"help"`
	}
)

func (r *sgResponse) Unmarshal(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	resp := &sgResponse{}
	err := json.Unmarshal(buf, resp)
	if err != nil {
		return err
	}
	*r = *resp
	return nil
}

func (r *sgResponse) CheckError(response *http.Response, buf []byte) error {
	if client.Is2XX(response.StatusCode) {
		return nil
	}
	if len(r.Errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s - message: %s, field: %s, help: %s", sendgridErrorMessage, r.Errors[0].Message, r.Errors[0].Field, r.Errors[0].Help)
}

func (r *sgResponse) Meta() httputil.Meta {
	return httputil.Meta{
		Message: "Successfully sent Sendgrid Email",
		// No response data from Sendgrid
		ID: "",
	}
}

func (d *sendGrid) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	tx := sgTransmission{
		From: &sgEmail{
			Name:    d.cfg.FromName,
			Address: d.cfg.FromAddress,
		},
		Subject: t.Subject,
		Personalizations: []*sgPersonalization{
			{Subject: t.Subject},
		},
		Content: []*sgContent{
			{Type: "text/plain", Value: t.PlainText},
			{Type: "text/html", Value: t.HTML},
		},
		Attachments: nil,
	}

	for _, r := range t.Recipients {
		tx.Personalizations[0].To = append(tx.Personalizations[0].To, &sgEmail{
			Address: r,
		})
	}

	if t.HasCC() {
		for _, c := range t.CC {
			tx.Personalizations[0].CC = append(tx.Personalizations[0].CC, &sgEmail{
				Address: c,
			})
		}
	}

	if t.HasBCC() {
		for _, b := range t.BCC {
			tx.Personalizations[0].BCC = append(tx.Personalizations[0].BCC, &sgEmail{
				Address: b,
			})
		}
	}

	if t.HasAttachments() {
		for _, v := range t.Attachments {
			tx.Attachments = append(tx.Attachments, &sgAttachment{
				Content:     v.B64(),
				Type:        v.Mime(),
				Name:        "",
				Filename:    v.Filename,
				Disposition: "attachment",
			})
		}
	}

	pl := httputil.NewJSONData()
	err = pl.AddStruct(tx)
	if err != nil {
		return mail.Response{}, err
	}

	req := httputil.NewHTTPRequest(http.MethodPost, sendGridEndpoint)
	req.AddHeader("Authorization", "Bearer "+d.cfg.APIKey)

	return d.client.Do(context.Background(), req, pl, &sgResponse{})
}
