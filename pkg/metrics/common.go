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

package metrics

import (
	"bytes"
	"io"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// ParseText read text and returns MetricFamily
func ParseText(input io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(input)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

// ParseMetricFamily read MetricFamily and returns text
func ParseMetricFamily(mfs map[string]*dto.MetricFamily) (string, error) {
	var buf bytes.Buffer
	for _, mf := range mfs {
		_, err := expfmt.MetricFamilyToText(&buf, mf)
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}
