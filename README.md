# 🏠 IoT Smart Home Backend with Golang

A lightweight, scalable backend built in Go for managing smart home appliances and streaming real-time sensor data using **WebSocket**, **RESTful API**, and **MQTT**. This backend is designed to communicate with ESP32 microcontrollers to monitor temperature and humidity and control home appliances such as lights or relays.

---

## 📌 Key Features

- ⚙️ Real-time temperature & humidity monitoring via **WebSocket**
- 💡 Relay (device) control using simple **REST API** endpoints
- 🔗 MQTT client integration (publish/subscribe)
- ⚡ Fast and concurrent processing using **Goroutines**
- 🧱 Modular code structure for easy extension

---


## 🖼️ System Architecture

```plaintext
+-------------+        MQTT         +-------------------+       HTTP / WS       +-------------+
|   ESP32     |  <-------------->   |    Golang Backend |  <----------------->  |  Web Client |
|  Sensors +  |                    |                   |                        | (Chart.js)  |
|  Relays     |                    |                   |                        |             |
+-------------+                    +-------------------+                        +-------------+

```
## 🏃 Getting Started

### 1. Requirements

- Go 1.18+
- MQTT Broker (e.g., [EMQX](https://www.emqx.io/) or Mosquitto)
- ESP32 devices with sensor firmware
- Frontend client (HTML+Chart.js or React)

### 2. Clone and Run

```bash
git clone https://github.com/yourusername/iot-smart-home-backend.git
cd iot-smart-home-backend
go mod tidy
go run main.go
```
🔌 API Endpoints

| Method | Endpoint      | Function                            |
| ------ | ------------- | ----------------------------------- |
| GET    | `/ws`         | WebSocket connection for live data  |
| POST   | `/relay/on`   | Turn ON Relay 1 (Living Room Light) |
| POST   | `/relay/off`  | Turn OFF Relay 1                    |
| POST   | `/relay2/on`  | Turn ON Relay 2 (Kitchen Light)     |
| POST   | `/relay2/off` | Turn OFF Relay 2                    |

📡 MQTT Topics

| Direction    | Topic                | Payload Format                                   |
| ------------ | -------------------- | ------------------------------------------------ |
| ➕ Publish    | `emqx/IoTcontrol`    | `{"relay1": "ON"}` or `{"relay2": "OFF"}`        |
| 📥 Subscribe | `sensor/temperature` | `{"temperature": 28.5, "timestamp": 1719322431}` |
| 📥 Subscribe | `sensor/humidity`    | `{"humidity": 65.1, "timestamp": 1719322431}`    |

#🧠 How it Works
-ESP32 devices publish temperature & humidity data to MQTT broker.

-The Golang backend subscribes to these topics and forwards the data to connected WebSocket clients.

-Users can toggle relays (like lights) from the frontend, which sends POST requests to the backend.

-Backend publishes MQTT control messages (e.g., {"relay1": "ON"}) to the broker to control devices.s

🧪 Sample Payloads
Sensor Data (from ESP32)(json):
```plaintext
{
  "temperature": 29.2,
  "humidity": 70.1,
  "timestamp": 1719322431
}
```
Relay Control from Backend to MQTT Broker:
```plaintext
{
  "relay2": "OFF"
}
```
