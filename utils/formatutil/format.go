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

package formatutil

import (
	"fmt"
	"strings"
)

// ProgressBar generates a formatted progress bar string.
func ProgressBar(name string, progress, total int) string {
	var (
		percentage     float64
		barLength      = 50
		progressLength int
	)

	if total == 0 {
		percentage = 0
		progressLength = 0
	} else {
		percentage = float64(progress) / float64(total) * 100
		barLength = 50
		progressLength = int(percentage / 100 * float64(barLength))
	}
	progressLength = min(progressLength, barLength)
	bar := strings.Repeat("█", progressLength) + strings.Repeat(" ", barLength-progressLength)
	return fmt.Sprintf("\r%s: [%s] %3.0f%% (%d/%d)", name, bar, percentage, progress, total)
}
