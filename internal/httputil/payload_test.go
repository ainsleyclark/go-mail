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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONData_AddStruct(t *testing.T) {
	tt := map[string]struct {
		input interface{}
		want  interface{}
	}{
		"Success": {
			map[string]interface{}{"test": 1},
			nil,
		},
		"Marshal Error": {
			map[string]interface{}{"test": make(chan struct{})},
			"unsupported type",
		},
		"Unmarshal Error": {
			1,
			"cannot unmarshal",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			pl := NewJSONData()
			err := pl.AddStruct(test.input)
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
			assert.NotNil(t, pl.original)
			assert.NotNil(t, pl.values)
		})
	}
}

func TestJSONData_Buffer(t *testing.T) {
	tt := map[string]struct {
		input jsonData
		want  interface{}
	}{
		"Success": {
			jsonData{values: map[string]interface{}{"test": 1}},
			`{"test":1}`,
		},
		"Marshal Error": {
			jsonData{values: map[string]interface{}{"test": make(chan struct{})}},
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
	pl := jsonData{}
	got := pl.ContentType()
	assert.Equal(t, JSONContentType, got)
}

func TestJSONData_Values(t *testing.T) {
	pl := jsonData{values: map[string]interface{}{"test": 1}}
	got := pl.Values()
	want := map[string]string{"test": "1"}
	assert.Equal(t, want, got)
}
