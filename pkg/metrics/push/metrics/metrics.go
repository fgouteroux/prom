// Package Metrics
package metrics

import (
	"fmt"
	"io"
	"time"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/prometheus/prompb"

	dto "github.com/prometheus/client_model/go"
)

var MetricMetadata_MetricType_value = map[string]int32{
	"UNKNOWN":        0,
	"COUNTER":        1,
	"GAUGE":          2,
	"HISTOGRAM":      3,
	"GAUGEHISTOGRAM": 4,
	"SUMMARY":        5,
	"INFO":           6,
	"STATESET":       7,
}

// FormatData convert metric family to a writerequest
func FormatData(mf map[string]*dto.MetricFamily, jobLabel string) *prompb.WriteRequest {
	wr := &prompb.WriteRequest{}

	for metricName, data := range mf {
		// Set metadata writerequest
		mtype := MetricMetadata_MetricType_value[data.Type.String()]
		metadata := prompb.MetricMetadata{
			MetricFamilyName: data.GetName(),
			Type:             prompb.MetricMetadata_MetricType(mtype),
			Help:             data.GetHelp(),
		}
		wr.Metadata = append(wr.Metadata, metadata)

		for _, metric := range data.Metric {
			timeserie := prompb.TimeSeries{
				Labels: []prompb.Label{
					{
						Name:  "__name__",
						Value: metricName,
					},
					{
						Name:  "job",
						Value: jobLabel,
					},
				},
			}

			for _, label := range metric.Label {
				labelname := label.GetName()
				if labelname == "job" {
					labelname = fmt.Sprintf("%s_exported", labelname)
				}
				timeserie.Labels = append(timeserie.Labels, prompb.Label{
					Name:  labelname,
					Value: label.GetValue(),
				})
			}

			timeserie.Samples = []prompb.Sample{
				{
					Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
					Value:     GetValue(metric),
				},
			}

			wr.Timeseries = append(wr.Timeseries, timeserie)
		}
	}
	return wr
}

// ParseTextReader consumes an io.Reader and returns the MetricFamily
func ParseTextReader(input io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(input)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

// GetValue return the value of a timeserie without the need to give value type
func GetValue(m *dto.Metric) float64 {
	switch {
	case m.Gauge != nil:
		return m.GetGauge().GetValue()
	case m.Counter != nil:
		return m.GetCounter().GetValue()
	case m.Untyped != nil:
		return m.GetUntyped().GetValue()
	default:
		return 0.
	}
}

// ParseTextAndFormat return the data in the expected prometheus metrics write request format
func ParseTextAndFormat(input io.Reader, jobLabel string) (*prompb.WriteRequest, error) {
	mf, err := ParseTextReader(input)
	if err != nil {
		return nil, err
	}
	return FormatData(mf, jobLabel), nil
}
