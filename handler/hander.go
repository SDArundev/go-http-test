package handler

import (
	"fmt"
	"net/http"
	"test-app/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	DevMode bool
	DB *db.Queries
}

func New(devMode bool, db *db.Queries) *Handler {
	return &Handler{
		DevMode: devMode,
		DB: db,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if h.DevMode {
		user, err := h.DB.GetUserByUsername(r.Context(), "admin")
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "Hello, World in dev mode! User: %s", user.Username)
	} else {			
		fmt.Fprintf(w, "Hello, World!")
	}
}

// search users get params from url
func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	users, err := h.DB.SearchUsers(r.Context(), pgtype.Text{String: username, Valid: true})
	if err != nil {
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Users: %v", users)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.DB.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Error getting all users", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Users: %v", users)
}
// login user and retuner jwt token
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	user, err := h.DB.Login(r.Context(), db.LoginParams{
		Username: username,
		Password: password,
	})
	if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
		return
	}

	if user.Username == "" {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}	
	
	fmt.Fprintf(w, "User: %v", token)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if err := h.DB.Register(r.Context(), db.RegisterParams{
		Username: username,
		Password: password,

	}); err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User registered: %s", username)
}

// admin dashboard, check if isAdmin to login		
func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	user, err := h.DB.VerifyIfUserIsAdmin(r.Context(), username)
	if err != nil {
		http.Error(w, "Error checking if user is admin", http.StatusInternalServerError)
		return
	}

	if !user {
		http.Error(w, "User is not an admin", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Admin dashboard: %v", user)
}