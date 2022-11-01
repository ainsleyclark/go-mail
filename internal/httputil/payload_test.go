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
	"github.com/ainsleyclark/go-mail/internal/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"testing"
)

func TestNewJSONData(t *testing.T) {
	tt := map[string]struct {
		input interface{}
		want  interface{}
	}{
		"Success": {
			map[string]interface{}{"test": 1},
			map[string]interface{}{"test": float64(1)},
		},
		"Marshal Error": {
			map[string]interface{}{"test": make(chan struct{})},
			"Error marshalling payload",
		},
		"Unmarshal Error": {
			1,
			"Error unmarshalling payload",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			pl, err := NewJSONData(test.input)
			if err != nil {
				assert.Contains(t, errors.Message(err), test.want)
				return
			}
			assert.Equal(t, test.want, pl.values)
			assert.NotNil(t, pl.original)
		})
	}
}

func TestJSONData_Buffer(t *testing.T) {
	tt := map[string]struct {
		input JSONData
		want  interface{}
	}{
		"Success": {
			JSONData{values: map[string]interface{}{"test": 1}},
			`{"test":1}`,
		},
		"Marshal Error": {
			JSONData{values: map[string]interface{}{"test": make(chan struct{})}},
			"unsupported type",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := test.input.Buffer()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.Equal(t, test.want, got.String())
		})
	}
}

func TestJsonData_ContentType(t *testing.T) {
	pl := JSONData{}
	got := pl.ContentType()
	assert.Equal(t, JSONContentType, got)
}

func TestJSONData_Values(t *testing.T) {
	pl := JSONData{values: map[string]interface{}{"test": 1}}
	got := pl.Values()
	want := map[string]string{"test": "1"}
	assert.Equal(t, want, got)
}

func TestFormData_AddValue(t *testing.T) {
	pl := NewFormData()
	pl.AddValue("key", "value")
	want := map[string]string{"key": "value"}
	assert.Equal(t, want, pl.values)
}

func TestFormData_AddBuffer(t *testing.T) {
	pl := NewFormData()
	pl.AddBuffer("key", "file", []byte("value"))
	want := []keyNameBuff{
		{key: "key", name: "file", value: []byte("value")},
	}
	assert.Equal(t, want, pl.buffers)
}

type mockWriterError struct{}

func (m *mockWriterError) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestFormData_Buffer(t *testing.T) {
	tt := map[string]struct {
		input  FormData
		writer func(w io.Writer) *multipart.Writer
		want   interface{}
	}{
		"Success": {
			FormData{
				values: map[string]string{
					"key": "value",
				},
				buffers: []keyNameBuff{
					{key: "key", name: "file", value: []byte("value")},
				},
			},
			multipart.NewWriter,
			"Content-Disposition",
		},
		"Value Error": {
			FormData{
				values: map[string]string{"key": "value"},
			},
			func(w io.Writer) *multipart.Writer {
				return multipart.NewWriter(&mockWriterError{})
			},
			"write error",
		},
		"Buffer Error": {
			FormData{
				buffers: []keyNameBuff{
					{key: "key", name: "file", value: []byte("value")},
				},
			},
			func(w io.Writer) *multipart.Writer {
				return multipart.NewWriter(&mockWriterError{})
			},
			"write error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			orig := newWriter
			defer func() { newWriter = orig }()
			newWriter = test.writer

			got, err := test.input.Buffer()
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}

			assert.Contains(t, got.String(), test.want)
		})
	}
}

func TestFormData_ContentType(t *testing.T) {
	pl := FormData{values: map[string]string{"test": "1"}}
	got := pl.ContentType()
	want := "multipart/form-data"
	assert.Contains(t, got, want)
}

func TestFormData_Values(t *testing.T) {
	pl := FormData{values: map[string]string{"test": "1"}}
	got := pl.Values()
	want := map[string]string{"test": "1"}
	assert.Equal(t, want, got)
}
