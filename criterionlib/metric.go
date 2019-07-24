package criterionlib

// Metric - A map of metrics to their values.
type Metric struct {
	name  string
	value float64
}

// NewMetric - Constructs a new metric with name and value.
func NewMetric(name string, value float64) Metric {
	return Metric{name: name, value: value}
}

// Name - Metric name.
func (m Metric) Name() string {
	return m.name
}

// Value - Metric value.
func (m Metric) Value() float64 {
	return m.value
}
