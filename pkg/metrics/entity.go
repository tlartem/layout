package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Entity struct {
	total    *prometheus.CounterVec
	duration *prometheus.HistogramVec
	current  *prometheus.GaugeVec
}

func NewEntity() *Entity {
	m := &Entity{}

	m.total = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "entity_processed_total",
		Help: "Total number of processed entities",
	}, []string{"name", "status"})
	prometheus.MustRegister(m.total)

	m.duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "entity_processed_duration",
		Help:    "Duration of processing entities (in seconds)",
		Buckets: buckets,
	}, []string{"name"})
	prometheus.MustRegister(m.duration)

	m.current = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "entity_processing_current",
		Help: "Current number of processing entities",
	}, []string{"name"})
	prometheus.MustRegister(m.current)

	return m
}

func (e *Entity) Total(name string, status Status) {
	e.total.WithLabelValues(name, status.String()).Inc()
}

func (e *Entity) TotalAdd(name string, status Status, counter int) {
	e.total.WithLabelValues(name, status.String()).Add(float64(counter))
}

func (e *Entity) Duration(name string, startTime time.Time) {
	e.duration.WithLabelValues(name).Observe(time.Since(startTime).Seconds())
}

func (e *Entity) Current(name string, value float64) {
	e.current.WithLabelValues(name).Set(value)
}
