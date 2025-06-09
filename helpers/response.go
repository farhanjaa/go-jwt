package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseWithData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Role    string `json:"role"`
}

type ResponseWithoutData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(w http.ResponseWriter, code int, message string, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var response interface{}
	status := "success"

	if code >= 400 {
		status = "failed"
	}

	if payload != nil {
		response = &ResponseWithData{
			Status:  status,
			Message: message,
			Data:    payload, // Ensure Data field is populated
			Role:    "",      // Tambahkan role jika diperlukan, bisa diisi dari payload atau context
		}
	} else {
		response = &ResponseWithoutData{
			Status:  status,
			Message: message,
		}
	}

	// Encode response ke JSON
	res, err := json.MarshalIndent(response, "", "  ") // Indentasi agar lebih mudah dibaca di terminal
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Cetak JSON ke terminal (VSCode)
	log.Printf("Response Sent: \n%s\n", res)

	// Kirim JSON response ke client
	w.Write(res)
}
