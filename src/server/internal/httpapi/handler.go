package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"app/server/internal/auth"
	"app/server/internal/user"
)

type Handler struct {
	users         user.Store
	auth          *auth.Service
	allowedOrigin string
}

func NewHandler(users user.Store, authService *auth.Service, allowedOrigin string) *Handler {
	return &Handler{users: users, auth: authService, allowedOrigin: allowedOrigin}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/dummy-data", h.method(http.MethodGet, h.health))
	mux.HandleFunc("/login", h.method(http.MethodPost, h.login))
	mux.HandleFunc("/create-user", h.method(http.MethodPost, h.createUser))
	mux.HandleFunc("/show-users", h.method(http.MethodGet, h.authorized(h.showUsers)))
	mux.HandleFunc("/delete", h.method(http.MethodPost, h.authorized(h.deleteUser)))
	mux.HandleFunc("/logout", h.method(http.MethodGet, h.authorized(h.logout)))
	return h.cors(mux)
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string          `json:"access-token"`
	User        user.PublicUser `json:"user"`
	Error       string          `json:"error"`
}

type usersResponse struct {
	Users   []user.PublicUser `json:"user_array"`
	Success string            `json:"success"`
	Error   string            `json:"error"`
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"title": "Golang + Angular Starter Kit"})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request credentials
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	authenticated, err := h.users.Authenticate(request.Email, request.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, err := h.auth.Create(authenticated.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	writeJSON(w, http.StatusOK, loginResponse{AccessToken: token, User: authenticated})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var request registration
	if err := decodeJSON(r, &request); err != nil || request.Name == "" || request.Email == "" || request.Password == "" {
		writeError(w, http.StatusBadRequest, "name, email, and password are required")
		return
	}
	_, err := h.users.Create(request.Name, request.Email, request.Password)
	if errors.Is(err, user.ErrEmailAlreadyExists) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"success": "User successfully created", "error": ""})
}

func (h *Handler) showUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.users.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list users")
		return
	}
	writeJSON(w, http.StatusOK, usersResponse{Users: users, Success: "success"})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID int `json:"id"`
	}
	if err := decodeJSON(r, &request); err != nil || request.ID <= 0 {
		writeError(w, http.StatusBadRequest, "a valid user id is required")
		return
	}
	if err := h.users.Delete(request.ID); errors.Is(err, user.ErrNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "could not delete user")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"success": "User successfully deleted", "error": ""})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	claims, err := h.auth.Validate(r.Header.Get("Authorization"))
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	h.auth.Revoke(claims.SessionID)
	writeJSON(w, http.StatusOK, map[string]string{"success": "success", "error": ""})
}

func (h *Handler) authorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := h.auth.Validate(r.Header.Get("Authorization")); err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		next(w, r)
	}
}

func (h *Handler) method(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != method {
			w.Header().Set("Allow", method)
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		next(w, r)
	}
}

func (h *Handler) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", h.allowedOrigin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		next.ServeHTTP(w, r)
	})
}

func decodeJSON(r *http.Request, destination interface{}) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message, "status": strconv.Itoa(status)})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
