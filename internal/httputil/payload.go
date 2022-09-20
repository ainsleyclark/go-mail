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

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ainsleyclark/go-mail/internal/errors"
	"io"
	"mime/multipart"
)

// Payload defines the methods used for creating  HTTP payload
// helper.
type Payload interface {
	// Buffer returns the byte buffer used for making the
	// HTTP Request
	Buffer() (*bytes.Buffer, error)
	// ContentType returns the `Content-Type` header used for
	// making the HTTP Request
	ContentType() string
	// Values returns a map of key - value pairs used for testing
	// and debugging.
	Values() map[string]string
}

const (
	// JSONContentType is the Content-Type header for
	// JSON payloads.
	JSONContentType = "application/json"
)

// JSONData defines the payload for JSON types.
type JSONData struct {
	original interface{}
	values   map[string]interface{}
}

// NewJSONData creates a new JSON Data Payload type.
// It adds a struct type to the JSON Data payload.
// Returns an error if the struct could not be marshalled or unmarshalled.
func NewJSONData(obj interface{}) (*JSONData, error) {
	const op = "HTTPUtil.NewJSONData"

	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, &errors.Error{Code: errors.INTERNAL, Message: "Error marshalling payload", Operation: op, Err: err}
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, &errors.Error{Code: errors.INTERNAL, Message: "Error unmarshalling payload", Operation: op, Err: err}
	}

	return &JSONData{
		original: obj,
		values:   m,
	}, nil
}

// Buffer returns the byte buffer for making the request.
func (j *JSONData) Buffer() (*bytes.Buffer, error) {
	const op = "JSONData.Buffer"
	buf, err := json.Marshal(j.values)
	if err != nil {
		return nil, &errors.Error{Code: errors.INTERNAL, Message: "Error marshalling values", Operation: op, Err: err}
	}
	return bytes.NewBuffer(buf), nil
}

// ContentType returns the `Content-Type` header.
func (j *JSONData) ContentType() string {
	return JSONContentType
}

// Values returns a map of key - value pairs used for testing
// and debugging.
func (j *JSONData) Values() map[string]string {
	m := make(map[string]string)
	for key, value := range j.values {
		m[key] = fmt.Sprintf("%v", value)
	}
	return m
}

// FormData defines the payload for URL encoded types.
type FormData struct {
	contentType string
	values      map[string]string
	buffers     []keyNameBuff
}

// keyNameBuff defines the buffer for multipart attachments.
type keyNameBuff struct {
	key   string
	name  string
	value []byte
}

// NewFormData creates a new Form Data Payload type.
func NewFormData() *FormData {
	return &FormData{}
}

// AddValue adds a key - value string pair to the Payload.
func (f *FormData) AddValue(key, value string) {
	if len(f.values) == 0 {
		f.values = make(map[string]string)
	}
	f.values[key] = value
}

// AddBuffer adds a file buffer to the Payload with a filename.
func (f *FormData) AddBuffer(key, fileName string, buff []byte) {
	f.buffers = append(f.buffers, keyNameBuff{
		key:   key,
		name:  fileName,
		value: buff,
	})
}

var newWriter = multipart.NewWriter

// Buffer returns the byte buffer for making the request.
func (f *FormData) Buffer() (*bytes.Buffer, error) {
	const op = "FormData.Buffer"

	data := &bytes.Buffer{}
	writer := newWriter(data)
	defer writer.Close()

	for key, val := range f.values {
		if tmp, err := writer.CreateFormField(key); err == nil {
			tmp.Write([]byte(val)) // nolint
		} else {
			return nil, &errors.Error{Code: errors.INTERNAL, Message: "Error creating form field", Operation: op, Err: err}
		}
	}

	for _, buff := range f.buffers {
		if tmp, err := writer.CreateFormFile(buff.key, buff.name); err == nil {
			r := bytes.NewReader(buff.value)
			io.Copy(tmp, r) // nolint
		} else {
			return nil, &errors.Error{Code: errors.INTERNAL, Message: "Error creating form file", Operation: op, Err: err}
		}
	}

	f.contentType = writer.FormDataContentType()

	return data, nil
}

// ContentType returns the `Content-Type` header.
func (f *FormData) ContentType() string {
	if f.contentType == "" {
		f.Buffer() // nolint
	}
	return f.contentType
}

// Values returns a map of key - value pairs used for testing
// and debugging.
func (f *FormData) Values() map[string]string {
	return f.values
}
