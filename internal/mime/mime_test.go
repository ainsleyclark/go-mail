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

package mime

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectBuffer(t *testing.T) {
	tt := map[string]struct {
		input string
		want  string
	}{
		"PNG": {
			"gopher.png",
			"image/png",
		},
		"JPG": {
			"gopher.jpg",
			"image/jpeg",
		},
		"SVG": {
			"gopher.svg",
			"image/svg+xml",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			wd, err := os.Getwd()
			assert.NoError(t, err)

			path := filepath.Join(filepath.Join(wd, "../../testdata"), test.input)
			file, err := os.ReadFile(path)
			assert.NoError(t, err)

			got := DetectBuffer(file)
			assert.Equal(t, test.want, got)
		})
	}
}
