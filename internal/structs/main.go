package structs

type (
	Gauge   float64
	Counter int64
)

type MemStorage struct {
	GaugeData   map[string]Gauge
	CounterData map[string]Counter
}

var Storage = MemStorage{
	GaugeData:   make(map[string]Gauge),
	CounterData: make(map[string]Counter),
}
