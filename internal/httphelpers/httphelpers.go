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

package httphelpers

import (
	"bytes"
	"io"
	"mime/multipart"
)

type FormData struct {
	ContentType string
	values      []keyValuePair
	buffers     []keyNameBuff
}

type keyValuePair struct {
	key   string
	value string
}

type keyNameBuff struct {
	key   string
	name  string
	value []byte
}

func (f *FormData) AddValue(key, value string) {
	f.values = append(f.values, keyValuePair{key: key, value: value})
}

func (f *FormData) AddBuffer(key, file string, buff []byte) {
	f.buffers = append(f.buffers, keyNameBuff{key: key, name: file, value: buff})
}

func (f *FormData) GetPayloadBuffer() (*bytes.Buffer, error) {
	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	defer writer.Close()

	for _, keyVal := range f.values {
		if tmp, err := writer.CreateFormField(keyVal.key); err == nil {
			tmp.Write([]byte(keyVal.value))
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

	f.ContentType = writer.FormDataContentType()

	return data, nil
}
