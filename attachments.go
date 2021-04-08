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

import (
	"encoding/base64"
	"github.com/gabriel-vasile/mimetype"
)

// Attachments defines the slice of mail attachments.
type Attachments []Attachment

// Attachment defines the mail file that has been
// uploaded via the forms endpoint. It contains
// useful information for sending files over
// the mail driver.
type Attachment struct {
	Filename string
	Bytes    []byte
}

// Determines if there are any attachments in the slice.
func (a Attachments) Exists() bool {
	return len(a) != 0
}

// Returns the Mime type of the byte data.
func (a Attachment) Mime() string {
	return mimetype.Detect(a.Bytes).String()
}

// Returns the base 64 encode of the byte data.
func (a Attachment) B64() string {
	return base64.StdEncoding.EncodeToString(a.Bytes)
}
