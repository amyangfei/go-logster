package logster

type Parser interface {
	Init(options string) error
	ParseLine(line string) error
	GetState(duration int) ([]*Metric, error)
}

type Metric struct {
	Name       string
	Value      float64
	Units      string
	Timestamp  int64
	MetricType string
}
