package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/Maiar0/api-sqlite-base-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterLoginPaths(mux *http.ServeMux) {
	log.Printf("[login.go] Register login routes")
	mux.HandleFunc("/auth/login", login)
	mux.HandleFunc("/auth/register", register)
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type loginResponse struct {
	TokenStr string `json:"token"`
	Message  string `json:"message,omitempty"`
}

func login(w http.ResponseWriter, r *http.Request) {
	//is it post?
	log.Printf("[login.go] login handler called")
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	var req loginRequest
	if err := utils.ReadRequestBody(w, r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[login.go] login request: %+v", req.Username) //TODO:: remove dont read password in logs
	//logic
	db, err := GetUserStore()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}
	user, err := db.GetUserByUsername(req.Username)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Error fetching user")
		return
	}
	if user == nil || !CheckPasswordHash(req.Password, user.passwordHash) {
		utils.WriteJSONError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	//response
	token, err := GenerateJWT(user.uuid, user.username, user.email, 24*time.Hour)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Error generating token")
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, loginResponse{TokenStr: token, Message: "Login successful"})

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
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	//read body into struct
	var req registerRequest
	if err := utils.ReadRequestBody(w, r, &req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	log.Printf("[login.go] register request: %+v", req) //TODO:: remove dont read password in logs
	//logic
	db, err := GetUserStore()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}
	result, err := db.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	log.Printf("[login.go] User created with result: %+v", result)
	//response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "User created successfully"})

}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TestAuthorized(w http.ResponseWriter, r *http.Request) {
	//is it get?
	log.Printf("[login.go] testAuthorized handler called")
	if r.Method != http.MethodGet {
		utils.WriteJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"data": "true"})
}
