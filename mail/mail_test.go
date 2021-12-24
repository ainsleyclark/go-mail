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
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// MailTestSuite defines the helper used for mail
// testing.
type MailTestSuite struct {
	suite.Suite
	base string
}

// Assert testing has begun.
func TestMail(t *testing.T) {
	suite.Run(t, new(MailTestSuite))
}

// Assigns test base.
func (t *MailTestSuite) SetupSuite() {
	wd, err := os.Getwd()
	t.NoError(err)
	t.base = wd
}

const (
	// DataPath defines where the test data resides.
	DataPath = "testdata"
	// PNGName defines the PNG name for testing.
	PNGName = "gopher.png"
	// JPGName defines the JPG name for testing.
	JPGName = "gopher.jpg"
	// SVGName defines the SVG name testing.
	SVGName = "gopher.svg"
)

// Returns a PNG attachment for testing.
func (t *MailTestSuite) Attachment(name string) Attachment {
	path := filepath.Join(filepath.Dir(t.base), DataPath, name)
	file, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fail("error getting attachment with the path: "+path, err)
	}

	return Attachment{
		Filename: name,
		Bytes:    file,
	}
}
