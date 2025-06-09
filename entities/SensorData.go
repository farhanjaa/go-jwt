package entities

// Struct untuk data dari ESP32
type SensorData struct {
	Device      string  `json:"device"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   int64   `json:"timestamp"`
}
