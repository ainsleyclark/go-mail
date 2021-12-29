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

package res

import (
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/client"
	"github.com/ainsleyclark/go-mail/mail"
	"net/http"
)

// mailchimp represents the entity for sending mail via the
// MailChimp API.
//
// See: https://mailchimp.com/developer/transactional/api/messages/send-new-message/
type mailchimp struct {
	cfg    mail.Config
	client clientold.Requester
}

// NewMailChimp creates a new Postal client. Configuration
// is validated before initialisation.
func NewMailChimp(cfg mail.Config) (mail.Mailer, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &mailchimp{
		cfg:    cfg,
		client: clientold.New("https://mandrillapp.com/api/1.0"),
	}, nil
}

const (
	// mailchimpEndpoint defines the endpoint to POST to.
	mailchimpEndpoint = "/messages/send"
	// mailchimpErrorMessage defines the message when an error occurred
	// when sending mail via the Postal API.
	mailchimpErrorMessage = "error sending transmission to MailChimp API"
)

// {"status":"error","code":-1,"name":"Invalid_Key","message":"Invalid API key"}

type (
	mailchimpTransmission struct {
		APIKey  string           `json:"key"`
		Message mailchimpMessage `json:"message"`
		Async   bool             `json:"async"`
		IPPool  string           `json:"ip_pool"`
		SendAt  string           `json:"send_at"`
	}
	// mailchimpMessage defines the data to be sent to the Postal API.
	mailchimpMessage struct {
		HTML      string        `json:"html"`
		Text      string        `json:"text"`
		Subject   string        `json:"subject"`
		FromEmail string        `json:"from_email"`
		FromName  string        `json:"from_name"`
		To        []mailchimpTo `json:"to"`
		Headers   struct {
		} `json:"headers"`
		Important               bool          `json:"important"`
		TrackOpens              bool          `json:"track_opens"`
		TrackClicks             bool          `json:"track_clicks"`
		AutoText                bool          `json:"auto_text"`
		AutoHTML                bool          `json:"auto_html"`
		InlineCSS               bool          `json:"inline_css"`
		URLStripQs              bool          `json:"url_strip_qs"`
		PreserveRecipients      bool          `json:"preserve_recipients"`
		ViewContentLink         bool          `json:"view_content_link"`
		BccAddress              string        `json:"bcc_address"`
		TrackingDomain          string        `json:"tracking_domain"`
		SigningDomain           string        `json:"signing_domain"`
		ReturnPathDomain        string        `json:"return_path_domain"`
		Merge                   bool          `json:"merge"`
		MergeLanguage           string        `json:"merge_language"`
		GlobalMergeVars         []interface{} `json:"global_merge_vars"`
		MergeVars               []interface{} `json:"merge_vars"`
		Tags                    []interface{} `json:"tags"`
		Subaccount              string        `json:"subaccount"`
		GoogleAnalyticsDomains  []interface{} `json:"google_analytics_domains"`
		GoogleAnalyticsCampaign string        `json:"google_analytics_campaign"`
		Metadata                struct {
			Website string `json:"website"`
		} `json:"metadata"`
		RecipientMetadata []interface{}         `json:"recipient_metadata"`
		Attachments       []mailchimpAttachment `json:"attachments"`
		Images            []interface{}         `json:"images"`
	}
	mailchimpTo struct {
		// the email address of the recipient
		Email string `json:"email"`
		// the optional display name to use for the recipient
		Name string `json:"name"`
		// the header type to use for the recipient, defaults to "to" if not provided Possible values: "to", "cc", or "bcc".
		Type string `json:"type"`
	}
	// mailchimpAttachment defines a singular Postal mail attachment.
	mailchimpAttachment struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Data string `json:"content"`
	}
)

// Send posts the Go Mail Transmission to the MailChimp
// API. Transmissions are validated before sending
// and attachments are added. Returns an error
// upon failure.
func (p *mailchimp) Send(t *mail.Transmission) (mail.Response, error) {
	err := t.Validate()
	if err != nil {
		return mail.Response{}, err
	}

	m := mailchimpTransmission{
		APIKey: p.cfg.APIKey,
		Message: mailchimpMessage{
			HTML:      t.HTML,
			Text:      t.PlainText,
			Subject:   t.Subject,
			FromEmail: p.cfg.FromAddress,
			FromName:  p.cfg.FromName,
		},
	}

	for _, recipient := range t.Recipients {
		m.Message.To = append(m.Message.To, mailchimpTo{
			Email: recipient,
			Type:  "to",
		})
	}

	if t.HasCC() {
		for _, c := range t.CC {
			m.Message.To = append(m.Message.To, mailchimpTo{
				Email: c,
				Type:  "cc",
			})
		}
	}

	if t.HasBCC() {
		for _, b := range t.BCC {
			m.Message.To = append(m.Message.To, mailchimpTo{
				Email: b,
				Type:  "bcc",
			})
		}
	}

	if t.Attachments.Exists() {
		for _, v := range t.Attachments {
			m.Message.Attachments = append(m.Message.Attachments, mailchimpAttachment{
				Name: v.Filename,
				Type: v.Mime(),
				Data: v.B64(),
			})
		}
	}

	// Ensure the API Key is set for authorisation
	// and add the JSON content type.
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	buf, resp, err := p.client.Do(m, mailchimpEndpoint, headers)
	fmt.Println(string(buf), resp)
	if err != nil {
		return mail.Response{
			StatusCode: resp.StatusCode,
			Body:       string(buf),
			Headers:    resp.Header,
			ID:         "",
			// TODO - Message
			Message: nil,
		}, err
	}

	return mail.Response{}, nil

	//// Unmarshal the buffer into a postalResponse.
	//response := postalResponse{}
	//err = json.Unmarshal(buf, &response)
	//if err != nil {
	//	return mail.Response{}, err
	//}
	//
	//// Bail if the status is not `success` and return formatted
	//// error code.
	//if response.HasError() {
	//	return mail.Response{}, response.Error()
	//}
	//
	//return response.ToResponse(resp, buf), nil
}
