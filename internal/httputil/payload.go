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

// Payload TODO
type Payload interface {
	Buffer() (*bytes.Buffer, error)
	ContentType() string
	Values() map[string]string
}

type jsonData struct {
	values map[string]interface{}
}

func NewJSONData(obj interface{}) (*jsonData, error) {
	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	return &jsonData{
		values: m,
	}, nil
}

func (j *jsonData) Buffer() (*bytes.Buffer, error) {
	buf, err := json.Marshal(j.values)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(buf), nil
}

func (j *jsonData) ContentType() string {
	return "application/json"
}

func (j *jsonData) Values() map[string]string {
	m := make(map[string]string)
	for key, value := range j.values {
		m[key] = fmt.Sprintf("%v", value)
	}
	return m
}

type formData struct {
	contentType string
	values      map[string]string
	buffers     []keyNameBuff
}

type keyNameBuff struct {
	key   string
	name  string
	value []byte
}

func NewFormData() *formData {
	return &formData{}
}

func (f *formData) AddValue(key, value string) {
	if len(f.values) == 0 {
		f.values = make(map[string]string)
	}
	f.values[key] = value
}

func (f *formData) AddBuffer(key, file string, buff []byte) {
	f.buffers = append(f.buffers, keyNameBuff{key: key, name: file, value: buff})
}

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

func (f *formData) ContentType() string {
	if f.contentType == "" {
		f.Buffer()
	}
	return f.contentType
}

func (f *formData) Values() map[string]string {
	return f.values
}
