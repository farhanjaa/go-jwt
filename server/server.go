package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go-jwt/entities"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var MqttClient mqtt.Client

var broadcast = make(chan entities.SensorData)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	influxToken  = "3Ip4wU-Cq7HV-x3ZawZoSttwqFfi8giSjaic53G4UxxN1B0OvmmaugRJg8WXpvbHc0tc7LrSAFgNgpjdmqe7yw=="
	influxURL    = "https://us-east-1-1.aws.cloud2.influxdata.com"
	influxDBName = "sensor_data"
)

// Handler untuk menerima data dari ESP32
func HandleIoTData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	jsonData, _ := json.Marshal(data)

	log.Printf("🔧 Sending to WebSocket: %s", jsonData) // ✅ Tambahkan ini

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Error sending to client:", err)
			client.Close()
			delete(clients, client)
		}
	}

	w.Write([]byte("Data received and broadcasted"))
}

func StartMQTTWebSocketServer() {
	// Init MQTT
	const statusTopic = "emqx/IoTstatus"
	broker := "tcp://broker.emqx.io:1883"
	topic := "emqx/IoTdata"
	clientID := "golang-backend"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// Init InfluxDB client
	influxClient, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     influxURL,
		Token:    influxToken,
		Database: influxDBName,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create InfluxDB client: %v", err)
	}

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Connected to MQTT broker.")

		// Subscriber untuk data sensor
		if token := c.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			handleMQTTMessage(msg, influxClient)
		}); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}

		// Subscriber untuk status relay
		if token := c.Subscribe(statusTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
			handleRelayStatus(msg)
		}); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}

	clientMQTT := mqtt.NewClient(opts)
	MqttClient = clientMQTT

	if token := clientMQTT.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// WebSocket routes
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	fmt.Println("🚀 WebSocket server started on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("❌ Failed to start WebSocket server:", err)
	}

}

func handleRelayStatus(msg mqtt.Message) {
	status := string(msg.Payload())
	log.Printf("⚙️ Relay Status: %s", status)

	// Kirim ke semua WebSocket client
	data := map[string]string{
		"type":   "status",
		"status": status,
	}
	jsonData, _ := json.Marshal(data)

	for client := range clients {
		client.WriteMessage(websocket.TextMessage, jsonData)
	}
}

func handleMQTTMessage(msg mqtt.Message, influxClient *influxdb3.Client) {
	fmt.Printf("📩 Received MQTT message on topic %s\n", msg.Topic())

	var data entities.SensorData
	if err := json.Unmarshal(msg.Payload(), &data); err != nil {
		log.Println("❌ Error decoding JSON:", err)
		return
	}

	// Convert timestamp to nanoseconds
	timestampNano := data.Timestamp * 1_000_000_000

	// Write to InfluxDB
	line := fmt.Sprintf("iot_sensor,device=%s temperature=%.2f,humidity=%.2f %d",
		data.Device, data.Temperature, data.Humidity, timestampNano)

	err := influxClient.Write(context.Background(), []byte(line))
	if err != nil {
		log.Println("❌ Failed to write to InfluxDB:", err)
		return
	}

	// Send to WebSocket clients
	broadcast <- data
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade error: %v", err)
		return
	}
	clients[ws] = true
	fmt.Println("🔗 New WebSocket client connected")

	// Loop untuk menahan koneksi tetap hidup
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("❌ WebSocket error: %v", err)
			delete(clients, ws)
			ws.Close()
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("❌ WebSocket error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
