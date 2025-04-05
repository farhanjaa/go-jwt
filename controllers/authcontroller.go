package controllers

import (
	"encoding/json"
	"go-jwt/configs"
	"go-jwt/helpers"
	"go-jwt/models"
	"html/template"
	"log"
	"net/http"
)

func AuthPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/home/auth.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var register models.Register

	// Cek apakah request dari JSON atau Form HTML
	if r.Header.Get("Content-Type") == "application/json" {
		// Jika JSON, baca dari body
		if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
			helpers.Response(w, 500, err.Error(), nil)
			return
		}
		log.Println("Parsed JSON:", register) // Debugging log
	} else {
		// Jika dari Form HTML, ambil dari FormValue
		register.Name = r.FormValue("name")
		register.Email = r.FormValue("email")
		register.Password = r.FormValue("password")
		register.PasswordConfirm = r.FormValue("password_confirm")
		log.Println("Parsed Form Data:", register) // Debugging log
	}
	defer r.Body.Close()

	// log.Println("PasswordConfirm received:", register.PasswordConfirm) // Debugging log

	// Cek password confirmation
	if register.Password != register.PasswordConfirm {
		helpers.Response(w, 400, "Password not match: "+register.Password+" password confirm: "+register.PasswordConfirm, nil)
		return
	}

	// Hash password
	passwordHash, err := helpers.HashPassword(register.Password)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	// Simpan user ke database
	user := models.User{
		Name:     register.Name,
		Email:    register.Email,
		Password: passwordHash,
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	helpers.Response(w, 201, "Register Successfully", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login models.Login

	// Parse the JSON body
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	var user models.User
	if err := configs.DB.First(&user, "email = ?", login.Email).Error; err != nil {
		helpers.Response(w, 404, "Wrong email or password", nil)
		return
	}

	if err := helpers.VerifyPassword(user.Password, login.Password); err != nil {
		helpers.Response(w, 404, "Wrong email or password", nil)
		return
	}

	// Create JWT token
	token, err := helpers.CreateToken(&user)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	// Create a success response
	response := helpers.ResponseWithData{
		Status:  "success",
		Message: "Successfully Login",
		Data:    token,
	}

	// Convert response to JSON
	res, err := json.Marshal(response)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	// Log the response to the terminal
	log.Println("Login Response: ", string(res)) // Log to VSCode terminal

	// Send the response to client (browser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
