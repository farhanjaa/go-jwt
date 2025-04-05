package main

import (
	"go-jwt/configs"
	"go-jwt/controllers"
	routes "go-jwt/routers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Koneksi ke database
	configs.ConnectDB()

	// Buat router utama
	r := mux.NewRouter()

	// Setup static file serving (untuk CSS, JS, dll.)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Tambahkan route untuk halaman autentikasi (GET request)
	r.HandleFunc("/auth", controllers.AuthPage).Methods("GET")

	// Tambahkan route untuk Login & Register dari Form HTML
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/register", controllers.Register).Methods("POST")

	// Buat subrouter untuk API (JSON request)
	apiRouter := r.PathPrefix("/api").Subrouter()
	routes.AuthRoutes(apiRouter)
	routes.UserRoutes(apiRouter)

	// Menjalankan server
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
