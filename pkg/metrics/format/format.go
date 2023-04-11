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

package fmt

import (
	"fmt"
	"strings"

	"github.com/fgouteroux/prom/pkg/metrics"
	"github.com/fgouteroux/prom/pkg/utils"
)

// Format prometheus metrics text as following convention
// https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format
func Format(input string) (string, error) {
	var metricsFmt string

	// remove duplicates lines
	strSlice := utils.UniqueStringSlice(strings.Split(input, "\n"))

	// end with a newline
	content := fmt.Sprintf("%s\n", strings.Join(strSlice, "\n"))

	// remove return carriage
	content = strings.ReplaceAll(content, "\r", "")

	// text to metric family
	mfs, err := metrics.ParseText(strings.NewReader(content))
	if err != nil {
		return metricsFmt, err
	}

	// metric family to text
	metricsFmt, err = metrics.ParseMetricFamily(mfs)
	if err != nil {
		return metricsFmt, err
	}

	return metricsFmt, nil
}
