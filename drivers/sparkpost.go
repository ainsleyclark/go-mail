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
	"strings"
	"time"
)

// postal represents the entity for sending mail via the
// Postal API.
//
// See:
// https://developers.sparkpost.com/api/
// https://developers.sparkpost.com/api/transmissions/#transmissions-create-a-transmission
type sparkPost struct {
	cfg    mail.Config
	client client.Requester
}

const (
	// sparkpostEndpoint defines the endpoint to POST to.
	// See: https://www.sparkpost.com/api#/reference/transmissions
	sparkpostEndpoint = "%s/api/v1/transmissions"
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
		client: client.New(),
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
	//
	// Example JSON Response:
	// {"results":{"total_rejected_recipients":0,"total_accepted_recipients":1,"id":"7029753512321354395"}}
	// {"errors":[{"message":"content.subject is a required field","code":"1400"}]}
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

func (r *spResponse) Unmarshal(buf []byte) error {
	resp := &spResponse{}
	err := json.Unmarshal(buf, resp)
	if err != nil {
		return err
	}
	*r = *resp
	return nil
}

func (r *spResponse) CheckError(response *http.Response, buf []byte) error {
	if len(r.Errors) == 0 {
		return nil
	}
	if len(buf) == 0 {
		return mail.ErrEmptyBody
	}
	return fmt.Errorf("%s - code: %s, message: %s", sparkpostErrorMessage, r.Errors[0].Code, r.Errors[0].Message)
}

func (r *spResponse) Meta() httputil.Meta {
	m := httputil.Meta{
		Message: "Successfully sent Sparkpost email",
	}
	if val, ok := r.Results["id"]; ok {
		m.ID = fmt.Sprintf("%v", val)
	}
	return m
}

func (d *sparkPost) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	headerTo := strings.Join(t.Recipients, ",")

	tx := spTransmission{
		Content: spContent{
			HTML:    t.HTML,
			Text:    t.PlainText,
			Subject: t.Subject,
			From: spFrom{
				Email: d.cfg.FromAddress,
				Name:  d.cfg.FromName,
			},
			ReplyTo: "",
			Headers: make(map[string]string),
		},
	}

	for _, r := range t.Recipients {
		tx.Recipients = append(tx.Recipients, spRecipient{
			Address: spAddress{Email: r, HeaderTo: headerTo},
		})
	}

	if t.HasCC() {
		for _, c := range t.CC {
			tx.Recipients = append(tx.Recipients, spRecipient{
				Address: spAddress{Email: c, HeaderTo: headerTo},
			})
			tx.Content.Headers["cc"] = strings.Join(t.CC, ",")
		}
	}

	if t.HasBCC() {
		for _, b := range t.BCC {
			tx.Recipients = append(tx.Recipients, spRecipient{
				Address: spAddress{Email: b, HeaderTo: headerTo},
			})
		}
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			tx.Content.Attachments = append(tx.Content.Attachments, spAttachment{
				MIMEType: v.Mime(),
				Filename: v.Filename,
				B64Data:  v.B64(),
			})
		}
	}

	pl := httputil.NewJSONData()
	err = pl.AddStruct(tx)
	if err != nil {
		return mail.Response{}, err
	}

	req := httputil.NewHTTPRequest(http.MethodPost, fmt.Sprintf(sparkpostEndpoint, d.cfg.URL))
	req.AddHeader("Authorization", d.cfg.APIKey)

	return d.client.Do(context.Background(), req, pl, &spResponse{})
}
