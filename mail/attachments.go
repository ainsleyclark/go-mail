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

package mail

import (
	"encoding/base64"

	"github.com/flightaware/go-mail/internal/mime"
)

// Attachment defines an email attachment for Go Mail.
// It contains useful information for sending files via
// the mail driver.
type Attachment struct {
	Filename string
	Bytes    []byte
}

// Mime returns the mime type of the byte data.
func (a Attachment) Mime() string {
	return mime.DetectBuffer(a.Bytes)
}

// B64 returns the base 64 encoding of the attachment.
func (a Attachment) B64() string {
	return base64.StdEncoding.EncodeToString(a.Bytes)
}
