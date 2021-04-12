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

package mail

// sendGrid represents the data for sending mail via the
// SendGrid API. Configuration, the client and the
// main send function are parsed for sending
// data.
//type smtpClient struct {
//	cfg Config
//	//send   sendGridSendFunc
//}
//
//// sendGridSendFunc defines the function for ending
//// SendGrid transmissions.
////type sendGridSendFunc func(email *mail.SGMailV3) (*rest.Response, error)
//
//// Creates a new sendGrid client. Configuration is
//// validated before initialisation.
//func newSMTP(cfg Config) *sendGrid {
//	return &sendGrid{
//		cfg: cfg,
//	}
//}
//
//// Send posts the go mail Transmission to the SendGrid
//// API. Transmissions are validated before sending
//// and attachments are added. Returns an error
//// upon failure.
//func (m *smtpClient) Send(t *Transmission) (Response, error) {
//	err := t.Validate()
//	if err != nil {
//		return Response{}, err
//	}
//
//	//auth := smtp.PlainAuth("gopher@example.net", "user@example.com", "password", m.cfg.URL)
//
////	smtp.SendMail(auth, "", "", "", "")
//
//	return Response{}, nil
//}
