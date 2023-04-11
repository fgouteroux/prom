/*
Copyright © 2023 François Gouteroux <francois.gouteroux@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package check

import (
	"bytes"
	"fmt"
	"io"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
)

const (
	successExitCode = 0
	failureExitCode = 1
	// Exit code 3 is used for "one or more lint issues detected".
	lintErrExitCode = 3
)

// Check performs a linting pass on input metrics.
// https://github.com/prometheus/prometheus/blob/6ddadd98b44cca7d55b27c20477123afac2201d7/cmd/promtool/main.go#L755
func Check(input io.Reader) (int, []string) {
	var errors []string
	var buf bytes.Buffer
	tee := io.TeeReader(input, &buf)

	l := promlint.New(tee)
	problems, err := l.Lint()
	if err != nil {
		errors = append(errors, fmt.Sprintf("error while linting: %v", err))
		return failureExitCode, errors
	}

	for _, p := range problems {
		errors = append(errors, fmt.Sprintln(p.Metric, p.Text))
	}

	if len(problems) > 0 {
		return lintErrExitCode, errors
	}

	return successExitCode, errors
}
