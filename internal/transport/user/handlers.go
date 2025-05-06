package user

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"minibank-backend/internal/entity"
	"minibank-backend/pkg"
	"minibank-backend/pkg/auth"
	"net/http"
	"strconv"
)

type IUserService interface {
	GetUserByID(ctx context.Context, UserID int) (entity.EUser, error)
	GetUsers(ctx context.Context) ([]entity.EUser, error)

	UpdateUserByID(ctx context.Context, UserID int, data map[string]interface{}) error

	SignUp(ctx context.Context, data map[string]interface{}) (int, string, error)
	SignIn(ctx context.Context, data map[string]interface{}) (int, string, error)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

type Handler struct {
	mw          auth.MiddleWare
	userService IUserService
}

func NewUserHandler(mw auth.MiddleWare, userService IUserService) *Handler {
	return &Handler{mw: mw, userService: userService}
}

func (h *Handler) Register(rtr *httprouter.Router) {
	rtr.GET("/api/me", h.mw.IsAuthed(h.GetMe))
	rtr.PUT("/api/me", h.mw.IsAuthed(h.UpdateProfile))

	rtr.GET("/api/users/", h.mw.IsAuthed(h.GetUsers))       // For Employee and Admin
	rtr.GET("/api/users/:id", h.mw.IsAuthed(h.GetUserByID)) // For Employee and Admin

	rtr.POST("/api/auth/sign-up", h.SignUp) // For All
	rtr.POST("/api/auth/sign-in", h.SignIn) // For All

}

func GetData(r *http.Request) (int, string) {
	userRole := r.Context().Value("user_role").(string)
	userID := r.Context().Value("user_id").(int)

	return userID, userRole
}

// GetMe TESTED Service level
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, _ := GetData(r)

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, user)

}

// SignUp TESTED
// But TODO Email accepting
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input CreateUserInput
	err := pkg.GetFromBody(r.Body, &input)
	if err != nil {
		return
	}

	data, err := input.ToMap()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, token, err := h.userService.SignUp(r.Context(), data)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}

	// TODO Accepting by Email
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"user_id": userID,
		"token":   token,
	})
}

// SignIn TESTED
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var input LoginUserInput
	err := pkg.GetFromBody(r.Body, &input)
	if err != nil {
		log.Println(err)
		return
	}

	data, err := input.ToMap()
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, token, err := h.userService.SignIn(r.Context(), data)
	if err != nil {
		log.Println(err)
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"token":   token,
	})
}

// UpdateProfile NOT TESTED
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	UserID, _ := GetData(r)

	var input UpdateUserInput

	err := pkg.GetFromBody(r.Body, &input)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, err)
	}

	data, err := input.ToMap()
	if err != nil {
		log.Println(err)
	}

	err = h.userService.UpdateUserByID(r.Context(), UserID, data)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "success update",
	})

}

// GetUserByID Complete and Work
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	_, role := GetData(r)

	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	users, err := h.userService.GetUsers(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, http.StatusOK, users)
}
