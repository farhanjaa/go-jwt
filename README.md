# 🏠 IoT Smart Home Backend with Golang

This project is a **Smart Home IoT Backend** built with **Golang**, featuring secure device control and real-time environmental monitoring (temperature and humidity). The backend system uses **JWT authentication**, **RBAC (Role-Based Access Control)**, and **InfluxDB 3.0** as a time-series database for storing sensor data.

---

## 🔐 Authentication & Authorization

- **JWT (JSON Web Token)** for stateless login & session management
- **RBAC (Role-Based Access Control)** to define user permissions:
  - `admin`: full control
  - `user`: limited access to devices and monitoring
  - `device`: only allowed to push data

---

## ⚙️ Features

- 🌡️ Real-time temperature & humidity tracking
- 💡 Control smart appliances (relays, lamps, etc.)
- 🔐 JWT login and token verification
- 📡 WebSocket live update integration
- 📦 REST API with secure route grouping
- 🧑‍💼 Role management using RBAC
- 🧠 Data logging with **InfluxDB 3.0** for time-series sensor data
- 🛠️ Clean and modular Go project structure

---

## 🧰 Tech Stack

| Layer         | Technology                    |
|---------------|-------------------------------|
| Backend Lang  | Go (Golang)                   |
| WebSocket     | Gorilla WebSocket             |
| REST API      | net/http + Gorilla Mux        |
| Auth          | Golang JWT v5                 |
| RBAC          | Custom middleware             |
| DB (Users)    | PostgreSQL / SQLite (GORM)    |
| DB (IoT Data) | **InfluxDB 3.0** (time-series)|
| Security      | bcrypt, token expiration      |

---

## 📦 InfluxDB 3.0 Integration

Sensor data (temperature, humidity) is pushed to **InfluxDB 3.0**, which stores values in a high-performance time-series format.

### Example measurement schema

- **Measurement**: `sensor_data`
- **Tags**: `device_id`, `location`
- **Fields**:
  - `temperature` (float)
  - `humidity` (float)
- **Timestamp**: `received_at`

### Sample Go insert code

```go
writeAPI := influxClient.WriteAPIBlocking(org, bucket)

point := influxdb3.NewPointWithMeasurement("sensor_data").
    AddTag("device_id", "esp32_01").
    AddField("temperature", 28.4).
    AddField("humidity", 60.2).
    SetTime(time.Now())

err := writeAPI.WritePoint(context.Background(), point)
```
---

## 🚀 API Overview

### 🔐 Auth Routes

| Method | Endpoint         | Description            |
|--------|------------------|------------------------|
| POST   | `/login`         | User login, returns JWT |
| POST   | `/register`      | Create new user (admin only) |

### 🌡️ Sensor & Device Control

| Method | Endpoint             | Role       | Description             |
|--------|----------------------|------------|-------------------------|
| GET    | `/sensor/data`       | user/admin | Get latest sensor data  |
| POST   | `/sensor/update`     | device     | Submit new sensor data  |
| POST   | `/relay/:id/on`      | user/admin | Turn relay ON           |
| POST   | `/relay/:id/off`     | user/admin | Turn relay OFF          |

---



## 🛡️ RBAC Middleware Example

```go
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    role := GetUserRoleFromJWT(r)
    if role != "admin" {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
    }
    next(w, r)
  }
}

📡 MQTT Topics
Ensure your ESP32 firmware publishes/subscribes to the following topics:

-emqx/IoTdata → Publishes temperature & humidity

-emqx/IoTcontrol/relay1 → Controls Lampu Ruang Tamu

-emqx/IoTcontrol/relay2 → Controls Lampu Dapur

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
🧪 Running the Project
1. Clone the repository
```plaintext
git clone https://github.com/yourusername/iot-smart-home-backend.git
cd iot-smart-home-backend
```
2. Configure .env
```plaintext
JWT_SECRET=your_jwt_secret
DB_URL=your_database_url
PORT=8080
```
3. Run the backend server
```plaintext
go run main.go
```
🧾 Example JWT Payload
```plaintext
{
  "user_id": 1,
  "role": "admin",
  "exp": 1723769123
}

```

📌 Notes
-The backend supports WebSocket communication for streaming real-time data to the front end.

-Designed to integrate with frontend dashboards (e.g., HTML + Chart.js).

-IoT devices send POST requests or open WebSocket channels to push updates.
