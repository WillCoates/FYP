package framework

// Graph represents a plottable line graph for graphutil
type Graph struct {
	ID    string
	Title string
	Data  []GraphPoint
}

type GraphPoint struct {
	Time  int64
	Value float64
}
