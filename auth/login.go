package auth

import (
	"log"
	"net/http"

	"github.com/Maiar0/api-sqlite-base-go/server"
)

func Register(mux *http.ServeMux) {
	log.Printf("[login.go] Register login routes")
	mux.HandleFunc("/auth/login", login)
	mux.HandleFunc("/auth/register", register)
}
func login(w http.ResponseWriter, r *http.Request) {
	//is it post?
	log.Printf("[login.go] login handler called")
	if r.Method != http.MethodPost {
		server.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}

}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func register(w http.ResponseWriter, r *http.Request) {
	//is it post?
	log.Printf("[login.go] register handler called")
	if r.Method != http.MethodPost {
		server.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	//read body into struct
	var req registerRequest
	if err := server.ReadRequestBody(w, r, &req); err != nil {
		server.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[login.go] register request: %+v", req) //TODO:: remove dont read password in logs
	//logic
	db, err := GetUserStore()
	if err != nil {
		server.WriteJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}
	result, err := db.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		server.WriteJSONError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	log.Printf("[login.go] User created with result: %+v", result)
	//response
	server.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "User created successfully"})

}
