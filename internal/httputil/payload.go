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

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	JSONContentType = "application/json"
)

// jsonData defines the payload for JSON types.
type jsonData struct {
	original interface{}
	values   map[string]interface{}
}

// NewJSONData creates a new JSON Data Payload type.
func NewJSONData() *jsonData {
	return &jsonData{}
}

// AddStruct adds a struct type to the JSON Data payload.
// Returns an error if the struct could not be marshalled or unmarshalled.
func (j *jsonData) AddStruct(obj interface{}) error {
	buf, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return err
	}

	j.values = m
	j.original = obj
	return nil
}

// Buffer returns the byte buffer for making the request.
func (j *jsonData) Buffer() (*bytes.Buffer, error) {
	buf, err := json.Marshal(j.values)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

// ContentType returns the `Content-Type` header.
func (j *jsonData) ContentType() string {
	return JSONContentType
}

// Values returns a map of key - value pairs used for testing
// and debugging.
func (j *jsonData) Values() map[string]string {
	m := make(map[string]string)
	for key, value := range j.values {
		m[key] = fmt.Sprintf("%v", value)
	}
	return m
}

// formData defines the payload for URL encoded types.
type formData struct {
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
func NewFormData() *formData {
	return &formData{}
}

// AddValue adds a key - value string pair to the Payload.
func (f *formData) AddValue(key, value string) {
	if len(f.values) == 0 {
		f.values = make(map[string]string)
	}
	f.values[key] = value
}

// AddBuffer adds a file buffer to the Payload with a filename.
func (f *formData) AddBuffer(key, fileName string, buff []byte) {
	f.buffers = append(f.buffers, keyNameBuff{
		key:   key,
		name:  fileName,
		value: buff,
	})
}

// Buffer returns the byte buffer for making the request.
func (f *formData) Buffer() (*bytes.Buffer, error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	defer writer.Close()

	for key, val := range f.values {
		if tmp, err := writer.CreateFormField(key); err == nil {
			tmp.Write([]byte(val))
		} else {
			return nil, err
		}
	}

	for _, buff := range f.buffers {
		if tmp, err := writer.CreateFormFile(buff.key, buff.name); err == nil {
			r := bytes.NewReader(buff.value)
			io.Copy(tmp, r)
		} else {
			return nil, err
		}
	}

	f.contentType = writer.FormDataContentType()

	return data, nil
}

// ContentType returns the `Content-Type` header.
func (f *formData) ContentType() string {
	if f.contentType == "" {
		f.Buffer()
	}
	return f.contentType
}

// Values returns a map of key - value pairs used for testing
// and debugging.
func (f *formData) Values() map[string]string {
	return f.values
}
