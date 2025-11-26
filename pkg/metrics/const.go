package metrics

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	All   Status = "all"
	Ok    Status = "ok"
	Error Status = "error"
)

//nolint:gochecknoglobals
var buckets = []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 40, 50, 60, 90, 120}
