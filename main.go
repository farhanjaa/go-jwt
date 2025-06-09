package main

import (
	"go-jwt/configs"
	"go-jwt/controllers"
	"go-jwt/middleware"
	routes "go-jwt/routers"
	"go-jwt/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Jalankan WebSocket Server di goroutine terpisah
	go server.StartMQTTWebSocketServer()

	// Koneksi ke database
	configs.ConnectDB()

	// Buat router utama
	r := mux.NewRouter()

	// Setup static file serving
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes halaman
	r.HandleFunc("/auth", controllers.AuthPage).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/register", controllers.Register).Methods("GET")
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/logout", controllers.Logout).Methods("GET")

	// Route untuk menerima data IoT dari ESP32
	r.HandleFunc("/iot-data", server.HandleIoTData).Methods("POST")

	// Monitoring hanya untuk admin
	r.Handle("/monitoring", middleware.Auth(middleware.AdminOnly(http.HandlerFunc(controllers.MonitoringPage)))).Methods("GET")

	// Endpoint kontrol relay (tanpa edit server.go)
	r.HandleFunc("/relay/on", func(w http.ResponseWriter, r *http.Request) {
		if server.MqttClient != nil && server.MqttClient.IsConnected() {
			server.MqttClient.Publish("emqx/IoTcontrol", 0, false, "Turn on")
			w.Write([]byte("Relay ON sent"))
		} else {
			http.Error(w, "MQTT client not connected", http.StatusServiceUnavailable)
		}
	}).Methods("GET", "POST")
	// Endpoint kontrol relay (tanpa edit server.go)

	r.HandleFunc("/relay/off", func(w http.ResponseWriter, r *http.Request) {
		if server.MqttClient != nil && server.MqttClient.IsConnected() {
			server.MqttClient.Publish("emqx/IoTcontrol", 0, false, "Turn off")
			w.Write([]byte("Relay OFF sent"))
		} else {
			http.Error(w, "MQTT client not connected", http.StatusServiceUnavailable)
		}
	}).Methods("GET", "POST")

	// Routes API
	apiRouter := r.PathPrefix("/api").Subrouter()
	routes.AuthRoutes(apiRouter)
	routes.UserRoutes(apiRouter)

	log.Println("🌐 HTTP server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
