package sensor

type Payload struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func (p Payload) EqualTo(p2 Payload) bool {
	return p.Temperature == p2.Temperature && p.Humidity == p2.Humidity
}
