// Copyright © 2024 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fileutil_test

import (
	"fmt"
	"github.com/zsjinwei/openim-tools/log/file-rotatelogs/internal/fileutil"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/lestrrat-go/strftime"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFn(t *testing.T) {
	// Mock time
	ts := []time.Time{
		{},
		(time.Time{}).Add(24 * time.Hour),
	}

	for _, xt := range ts {
		pattern, err := strftime.New("/path/to/%Y/%m/%d")
		if !assert.NoError(t, err, `strftime.New should succeed`) {
			return
		}
		clock := clockwork.NewFakeClockAt(xt)
		fn := fileutil.GenerateFn(pattern, clock, 24*time.Hour)
		expected := fmt.Sprintf("/path/to/%04d/%02d/%02d",
			xt.Year(),
			xt.Month(),
			xt.Day(),
		)

		if !assert.Equal(t, expected, fn) {
			return
		}
	}
}
