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

	// Render the register.html page for GET request
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("views/home/register.html")
		if err != nil {
			http.Error(w, "Error loading register page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

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
		Role:     "user", // default atau input dari form
	}

	if err := configs.DB.Create(&user).Error; err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	helpers.Response(w, 201, "Register Successfully", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Hapus token dari cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "token", // sesuaikan jika kamu menyimpan token di cookie dengan nama lain
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Hapus cookie
	})
	// Redirect ke halaman login
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var login models.Login

	// Handle GET for rendering login page
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("views/home/auth.html")
		if err != nil {
			http.Error(w, "Error loading login page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	// Handle POST from form or JSON
	if r.Header.Get("Content-Type") == "application/json" {
		// JSON login (API)
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			helpers.Response(w, 500, err.Error(), nil)
			return
		}
	} else {
		// Form login
		login.Email = r.FormValue("email")
		login.Password = r.FormValue("password")
	}

	if login.Email == "" || login.Password == "" {
		if r.Header.Get("Content-Type") == "application/json" {
			helpers.Response(w, 400, "Email and password must not be empty", nil)
		} else {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
		}
		return
	}
	// Cari user di DB
	var user models.User
	if err := configs.DB.First(&user, "email = ?", login.Email).Error; err != nil {
		// Handle untuk Form atau JSON
		if r.Header.Get("Content-Type") == "application/json" {
			helpers.Response(w, 404, "Wrong email or password", nil)
		} else {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
		}
		return
	}

	// Verifikasi password
	if err := helpers.VerifyPassword(user.Password, login.Password); err != nil {
		if r.Header.Get("Content-Type") == "application/json" {
			helpers.Response(w, 404, "Wrong email or password", nil)
		} else {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
		}
		return
	}

	// Berhasil login → token disiapkan
	token, err := helpers.CreateToken(&user)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	// Always set the cookie, for both JSON and form logins
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // for localhost
		SameSite: http.SameSiteLaxMode,
	})

	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "Login berhasil",
			"data":    token,
			"role":    user.Role,
		})
		return
	}

	// If login is from HTML form, redirect after setting cookie
	http.Redirect(w, r, "/static/iot-dashboard/landing.html", http.StatusSeeOther)
}

func MonitoringPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/iot-dashboard/index.html") // pastikan file ini ada
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
