package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"capstone1/config"
	"capstone1/internal/auth"
	"capstone1/internal/db"
	"capstone1/internal/session"
	"capstone1/internal/vault"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	dbConn := db.ConnectDB(cfg)
	db.InitTables(dbConn)

	r := mux.NewRouter()

	// Registration
	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		user, err := auth.Register(dbConn, req.Username, req.Password, req.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"username": user.Username,
			"totp":     user.TOTPSecret,
			"role":     user.Role,
		})
	}).Methods("POST")

	// Login
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			TOTP     string `json:"totp"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		token, err := auth.Login(dbConn, req.Username, req.Password, req.TOTP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		session.LogSession(dbConn, 0, req.Username+" logged in")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}).Methods("POST")

	// Vault: Add secret
	r.HandleFunc("/vault/add", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name   string `json:"name"`
			Secret string `json:"secret"`
			UserID int    `json:"user_id"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		if err := vault.StoreSecret(dbConn, req.UserID, req.Name, req.Secret); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		session.LogSession(dbConn, req.UserID, "Added vault secret "+req.Name)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("POST")

	// Vault: Get secret
	r.HandleFunc("/vault/get/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]
		secret, err := vault.RetrieveSecret(dbConn, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"secret": secret})
	}).Methods("GET")

	// Sessions: User sessions
	r.HandleFunc("/sessions/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		uid, _ := strconv.Atoi(mux.Vars(r)["user_id"])
		sessions, err := session.GetUserSessions(dbConn, uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sessions)
	}).Methods("GET")

	// Sessions: All sessions
	r.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		sessions, err := session.GetAllSessions(dbConn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sessions)
	}).Methods("GET")

	log.Println("Server running at http://localhost:" + cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
