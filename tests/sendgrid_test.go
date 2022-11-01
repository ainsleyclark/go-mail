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
	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
	"os"
	"testing"
)

func Test_SendGrid(t *testing.T) {
	LoadEnv(t)
	cfg := mail.Config{
		APIKey:      os.Getenv("SENDGRID_API_KEY"),
		FromAddress: os.Getenv("SENDGRID_FROM_ADDRESS"),
		FromName:    os.Getenv("SENDGRID_FROM_NAME"),
	}
	UtilTestSend(t, drivers.NewSendGrid, cfg, "SendGrid")
}
