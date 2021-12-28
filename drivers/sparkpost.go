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
	"encoding/json"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/client"
	"github.com/ainsleyclark/go-mail/mail"
	"net/http"
	"strings"
	"time"
)

// sparkPost represents the data for sending mail via the
// SparkPost API. Configuration, the client and the
// main send function are parsed for sending
// data.
type sparkPost struct {
	cfg    mail.Config
	client client.Requester
}

const (
	// sparkpostEndpoint defines the endpoint to POST to.
	// See: https://www.sparkpost.com/api#/reference/transmissions
	sparkpostEndpoint = "/api/v1/transmissions"
	// sparkpostErrorMessage defines the message when an error occurred
	// when sending mail via the Sparkpost API.
	sparkpostErrorMessage = "error sending transmission to Sparkpost API"
)

// NewSparkPost creates a new SparkPost client. Configuration
// is validated before initialisation.
func NewSparkPost(cfg mail.Config) (mail.Mailer, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &sparkPost{
		cfg:    cfg,
		client: client.New(cfg.URL),
	}, nil
}

type (
	// spTransmission is the JSON structure accepted by and returned
	// from the SparkPost Transmissions API.
	spTransmission struct {
		ID                   string                 `json:"id,omitempty"`
		State                string                 `json:"state,omitempty"`
		Options              *spTransmissionOptions `json:"options,omitempty"`
		Recipients           []spRecipient          `json:"recipients"`
		CampaignID           string                 `json:"campaign_id,omitempty"`
		Description          string                 `json:"description,omitempty"`
		Metadata             interface{}            `json:"metadata,omitempty"`
		SubstitutionData     interface{}            `json:"substitution_data,omitempty"`
		ReturnPath           string                 `json:"return_path,omitempty"`
		Content              spContent              `json:"content"`
		TotalRecipients      *int                   `json:"total_recipients,omitempty"`
		NumGenerated         *int                   `json:"num_generated,omitempty"`
		NumFailedGeneration  *int                   `json:"num_failed_generation,omitempty"`
		NumInvalidRecipients *int                   `json:"num_invalid_recipients,omitempty"`
	}
	// spTransmissionOptions specifies settings to apply to this Transmission.
	// If not specified, and present in TmplOptions, those values will be used.
	spTransmissionOptions struct {
		OpenTracking         *bool      `json:"open_tracking,omitempty"`
		ClickTracking        *bool      `json:"click_tracking,omitempty"`
		Transactional        *bool      `json:"transactional,omitempty"`
		StartTime            *time.Time `json:"start_time,omitempty"`
		Sandbox              *bool      `json:"sandbox,omitempty"`
		SkipSuppression      *bool      `json:"skip_suppression,omitempty"`
		IPPool               string     `json:"ip_pool,omitempty"`
		InlineCSS            *bool      `json:"inline_css,omitempty"`
		PerformSubstitutions *bool      `json:"perform_substitutions,omitempty"`
	}
	// spContent is what will be sent to recipients.
	// Knowledge of SparkPost's substitution/templating capabilities will come in handy here.
	// https://www.sparkpost.com/api#/introduction/substitutions-reference
	spContent struct {
		HTML         string            `json:"html,omitempty"`
		Text         string            `json:"text,omitempty"`
		Subject      string            `json:"subject,omitempty"`
		From         spFrom            `json:"from,omitempty"`
		ReplyTo      string            `json:"reply_to,omitempty"`
		Headers      map[string]string `json:"headers,omitempty"`
		EmailRFC822  string            `json:"email_rfc822,omitempty"`
		Attachments  []spAttachment    `json:"attachments,omitempty"`
		InlineImages []interface{}     `json:"inline_images,omitempty"`
	}
	// spFrom describes the nested object way of specifying the `From` header.
	// Content.From can be specified this way, or as a plain string.
	spFrom struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	// spResponse contains information about the last HTTP response from
	// the SparkPost API.
	spResponse struct {
		Results map[string]interface{} `json:"results,omitempty"`
		Errors  []spError              `json:"errors,omitempty"`
	}
	// spError mirrors the error format returned by SparkPost APIs.
	spError struct {
		Message     string `json:"message"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Part        string `json:"part,omitempty"`
		Line        int    `json:"line,omitempty"`
	}
	// spRecipient represents one email (you guessed it) recipient.
	spRecipient struct {
		Address          spAddress   `json:"address"`
		ReturnPath       string      `json:"return_path,omitempty"`
		Tags             []string    `json:"tags,omitempty"`
		Metadata         interface{} `json:"metadata,omitempty"`
		SubstitutionData interface{} `json:"substitution_data,omitempty"`
	}
	// spAddress describes the nested object way of specifying the
	// Recipient's email address. Recipient.Address can also be
	// a plain string.
	spAddress struct {
		Email    string `json:"email"`
		Name     string `json:"name,omitempty"`
		HeaderTo string `json:"header_to,omitempty"`
	}
	// spAttachment contains metadata and the contents of the
	// file to attach.
	spAttachment struct {
		MIMEType string `json:"type"`
		Filename string `json:"name"`
		B64Data  string `json:"data"`
	}
)

// HasError determines if the Sparkpost call was successful
// by evaluating the error slice length within the response.
func (p *spResponse) HasError() bool {
	return len(p.Errors) != 0
}

// Error returns a formatted response error for a Sparkpost
// response.
func (p *spResponse) Error() error {
	if len(p.Errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s - code: %s, message: %s", sparkpostErrorMessage, p.Errors[0].Code, p.Errors[0].Message)
}

// ToResponse transforms a spResponse into a Go Mail response.
// Checks if the id is attached and sets accordingly.
func (p *spResponse) ToResponse(resp *http.Response, buf []byte) mail.Response {
	response := mail.Response{
		StatusCode: resp.StatusCode,
		Body:       string(buf),
		Headers:    resp.Header,
		Message:    "Successfully sent Sparkpost email",
	}
	if val, ok := p.Results["id"]; ok {
		response.ID = fmt.Sprintf("%v", val)
	}
	return response
}

// Send posts the go mail Transmission to the SparkPost
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (s *sparkPost) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	headerTo := strings.Join(t.Recipients, ",")

	m := spTransmission{
		Content: spContent{
			HTML:    t.HTML,
			Text:    t.PlainText,
			Subject: t.Subject,
			From: spFrom{
				Email: s.cfg.FromAddress,
				Name:  s.cfg.FromName,
			},
			ReplyTo: "",
			Headers: make(map[string]string),
		},
	}

	for _, r := range t.Recipients {
		m.Recipients = append(m.Recipients, spRecipient{
			Address: spAddress{Email: r, HeaderTo: headerTo},
		})
	}

	if t.HasCC() {
		for _, c := range t.CC {
			m.Recipients = append(m.Recipients, spRecipient{
				Address: spAddress{Email: c, HeaderTo: headerTo},
			})
			m.Content.Headers["cc"] = strings.Join(t.CC, ",")
		}
	}

	if t.HasBCC() {
		for _, b := range t.BCC {
			m.Recipients = append(m.Recipients, spRecipient{
				Address: spAddress{Email: b, HeaderTo: headerTo},
			})
		}
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			m.Content.Attachments = append(m.Content.Attachments, spAttachment{
				MIMEType: v.Mime(),
				Filename: v.Filename,
				B64Data:  v.B64(),
			})
		}
	}

	// Ensure the API Key is set for authorisation
	// and add the JSON content type.
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Authorization", s.cfg.APIKey)

	buf, resp, err := s.client.Do(m, sparkpostEndpoint, headers)
	if err != nil {
		return mail.Response{}, err
	}

	// Unmarshal the buffer into a postalResponse.
	response := spResponse{}
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return mail.Response{}, err
	}

	if response.HasError() {
		return mail.Response{}, response.Error()
	}

	return response.ToResponse(resp, buf), nil
}
