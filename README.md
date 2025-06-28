# 🏠 IoT Smart Home Backend with Golang

This project is a **Smart Home IoT Backend** built with **Golang**, featuring secure device control and environmental monitoring (temperature and humidity). The backend system integrates **JWT authentication** and **RBAC (Role-Based Access Control)** to ensure secure and role-aware access for users and devices.

---

## 🔐 Authentication & Authorization

- **JWT (JSON Web Token)** for stateless, secure login and session handling
- **RBAC (Role-Based Access Control)** to differentiate access levels between:
  - `admin`: full access to all endpoints
  - `user`: restricted access (view & limited control)
  - `device`: allowed to push sensor data or control signals

---

## ⚙️ Features

- 📡 Real-time data collection from IoT sensors (e.g., ESP32)
- 🌡️ Monitor temperature and humidity via WebSocket
- 💡 Control relays (e.g., lamps, plugs) via REST API
- 🔐 Secure login system using JWT tokens
- 👥 Role-based authorization (admin/user/device)
- 🌍 WebSocket integration for live updates
- 📦 RESTful API for device and user management
- 📄 Clean folder structure and maintainable Go modules

---

## 🧰 Tech Stack

- **Language**: Go (Golang)
- **Routing**: `net/http`, `gorilla/mux`
- **WebSocket**: `gorilla/websocket`
- **Authentication**: `github.com/golang-jwt/jwt/v5`
- **Authorization**: Custom RBAC middleware
- **Database**: PostgreSQL or SQLite via GORM
- **Security**: JWT, password hashing (bcrypt), middleware validation

---


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

## 🛡️ RBAC Example

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
