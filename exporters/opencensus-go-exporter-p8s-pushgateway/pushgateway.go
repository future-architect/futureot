// Copyright 2019, Future Corporation
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

// Package prometheus contains a Prometheus exporter that supports exporting
// OpenCensus views as Prometheus metrics.package pushgateway
package pushgateway

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	"go.opencensus.io/stats/view"
)

type Exporter struct {
	opts   Options
	pusher *push.Pusher
	c      *collector
}

// ExportView is a dummy method. Exporter is a compatible with view.Exporter
func (e *Exporter) ExportView(vd *view.Data) {
}

type Options struct {
	// GatewayEndpoint is the full url to the Prometheus gateway.
	// For example, http://prometheus-gateway:9091
	// Default value is http://localhost:9091.
	GatewayEndpoint string

	// JobName is the job name assigned to scraped metrics by default.
	// Default value is process name.
	JobName string

	Namespace string

	// UserName is for authentaction
	UserName string

	// Password is for authentaction
	Password string

	OnError func(err error)

	// ConstLabels will be set as labels on all views.
	ConstLabels prometheus.Labels

	ReportInterval time.Duration
}

func (o *Options) onError(err error) {
	if o.OnError != nil {
		o.OnError(err)
	} else {
		log.Printf("Failed to export to Prometheus: %v", err)
	}
}

// NewExporter returns an exporter that exports stats to Prometheus pushgateway.
func NewExporter(o Options) (*Exporter, error) {
	if o.GatewayEndpoint == "" {
		o.GatewayEndpoint = "http://localhost:9091"
	}
	if o.JobName == "" {
		exe, e := os.Executable()
		if e != nil {
			return nil, errors.Wrap(e, "Error during getting process name for default job name")
		}
		o.JobName = filepath.Base(exe)
	}
	if o.ReportInterval == 0 {
		o.ReportInterval = time.Second
	}

	pusher := push.New(o.GatewayEndpoint, o.JobName)
	if o.UserName != "" || o.Password != "" {
		pusher.BasicAuth(o.UserName, o.Password)
	}
	collector := newCollector(o, pusher)
	e := &Exporter{
		opts:   o,
		pusher: pusher,
		c:      collector,
	}
	collector.ensureRegisteredOnce()
	return e, nil
}

// collector implements prometheus.Collector
type collector struct {
	opts Options
	mu   sync.Mutex // mu guards all the fields.

	registerOnce sync.Once

	pusher *push.Pusher

	// reader reads metrics from all registered producers.
	reader *metricexport.Reader
}

func newCollector(opts Options, pusher *push.Pusher) *collector {
	return &collector{
		pusher: pusher,
		opts:   opts,
		reader: metricexport.NewReader()}
}

func (c *collector) ensureRegisteredOnce() {
	c.registerOnce.Do(func() {
		c.pusher.Collector(c).Push()
		go func() {
			for {
				<-time.After(c.opts.ReportInterval)
				c.pusher.Push()
			}
		}()
	})
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	de := &descExporter{c: c, descCh: ch}
	c.reader.ReadAndExport(de)
}

// Collect fetches the statistics from OpenCensus
// and delivers them as Prometheus Metrics.
// Collect is invoked every time a prometheus.Gatherer is run
// for example when the HTTP endpoint is invoked by Prometheus.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	me := &metricExporter{c: c, metricCh: ch}
	c.reader.ReadAndExport(me)
}

func (c *collector) toDesc(metric *metricdata.Metric) *prometheus.Desc {
	return prometheus.NewDesc(
		metricName(c.opts.Namespace, metric),
		metric.Descriptor.Description,
		toPromLabels(metric.Descriptor.LabelKeys),
		c.opts.ConstLabels)
}

type metricExporter struct {
	c        *collector
	metricCh chan<- prometheus.Metric
}

// ExportMetrics exports to the Prometheus.
// Each OpenCensus Metric will be converted to
// corresponding Prometheus Metric:
// TypeCumulativeInt64 and TypeCumulativeFloat64 will be a Counter Metric,
// TypeCumulativeDistribution will be a Histogram Metric.
// TypeGaugeFloat64 and TypeGaugeInt64 will be a Gauge Metric
func (me *metricExporter) ExportMetrics(ctx context.Context, metrics []*metricdata.Metric) error {
	for _, metric := range metrics {
		desc := me.c.toDesc(metric)
		for _, ts := range metric.TimeSeries {
			tvs := toLabelValues(ts.LabelValues)
			for _, point := range ts.Points {
				metric, err := toPromMetric(desc, metric, point, tvs)
				if err != nil {
					me.c.opts.onError(err)
				} else if metric != nil {
					me.metricCh <- metric
				}
			}
		}
	}
	return nil
}

type descExporter struct {
	c      *collector
	descCh chan<- *prometheus.Desc
}

// ExportMetrics exports descriptor to the Prometheus.
// It is invoked when request to scrape descriptors is received.
func (me *descExporter) ExportMetrics(ctx context.Context, metrics []*metricdata.Metric) error {
	for _, metric := range metrics {
		desc := me.c.toDesc(metric)
		me.descCh <- desc
	}
	return nil
}

func toPromLabels(mls []metricdata.LabelKey) (labels []string) {
	for _, ml := range mls {
		labels = append(labels, sanitize(ml.Key))
	}
	return labels
}

func metricName(namespace string, m *metricdata.Metric) string {
	var name string
	if namespace != "" {
		name = namespace + "_"
	}
	return name + sanitize(m.Descriptor.Name)
}

func toPromMetric(
	desc *prometheus.Desc,
	metric *metricdata.Metric,
	point metricdata.Point,
	labelValues []string) (prometheus.Metric, error) {
	switch metric.Descriptor.Type {
	case metricdata.TypeCumulativeFloat64, metricdata.TypeCumulativeInt64:
		pv, err := toPromValue(point)
		if err != nil {
			return nil, err
		}
		return prometheus.NewConstMetric(desc, prometheus.CounterValue, pv, labelValues...)

	case metricdata.TypeGaugeFloat64, metricdata.TypeGaugeInt64:
		pv, err := toPromValue(point)
		if err != nil {
			return nil, err
		}
		return prometheus.NewConstMetric(desc, prometheus.GaugeValue, pv, labelValues...)

	case metricdata.TypeCumulativeDistribution:
		switch v := point.Value.(type) {
		case *metricdata.Distribution:
			points := make(map[float64]uint64)
			// Histograms are cumulative in Prometheus.
			// Get cumulative bucket counts.
			cumCount := uint64(0)
			for i, b := range v.BucketOptions.Bounds {
				cumCount += uint64(v.Buckets[i].Count)
				points[b] = cumCount
			}
			return prometheus.NewConstHistogram(desc, uint64(v.Count), v.Sum, points, labelValues...)
		default:
			return nil, typeMismatchError(point)
		}
	case metricdata.TypeSummary:
		// TODO: [rghetia] add support for TypeSummary.
		return nil, nil
	default:
		return nil, fmt.Errorf("aggregation %T is not yet supported", metric.Descriptor.Type)
	}
}

func toLabelValues(labelValues []metricdata.LabelValue) (values []string) {
	for _, lv := range labelValues {
		if lv.Present {
			values = append(values, lv.Value)
		} else {
			values = append(values, "")
		}
	}
	return values
}

func typeMismatchError(point metricdata.Point) error {
	return fmt.Errorf("point type %T does not match metric type", point)

}

func toPromValue(point metricdata.Point) (float64, error) {
	switch v := point.Value.(type) {
	case float64:
		return v, nil
	case int64:
		return float64(v), nil
	default:
		return 0.0, typeMismatchError(point)
	}
}
