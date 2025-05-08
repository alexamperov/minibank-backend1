package apply

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"minibank-backend/internal/entity"
	"minibank-backend/pkg"
	"minibank-backend/pkg/auth"
	"net/http"
	"strconv"
)

type Handler struct {
	mw auth.MiddleWare

	applyService IApplyService
}

// NewHandler constructor Handler struct
func NewHandler(
	mw auth.MiddleWare,
	applyService IApplyService,
) *Handler {
	return &Handler{
		mw:           mw,
		applyService: applyService,
	}
}

// respondJSON - func for send status-code and response data
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// GetData for get id and role from JWT-token
func GetData(r *http.Request) (int, string) {
	userRole := r.Context().Value("user_role").(string)
	userID := r.Context().Value("user_id").(int)

	return userID, userRole
}

type IApplyService interface {
	GetAllApplies(ctx context.Context, UserID int, isUser bool) ([]entity.EApply, error)
	GetApplyByID(ctx context.Context, UserID int, applyID int) (*entity.EApply, error)
	CreateApply(ctx context.Context, UserID int, applyData map[string]interface{}) (int, error)
	AcceptApply(ctx context.Context, UserID int, applyID int) error
	DenyApply(ctx context.Context, UserID int, applyID int) error
}

func (h *Handler) Register(rtr *httprouter.Router) {
	//Apply
	rtr.GET("/api/applies", h.mw.IsAuthed(h.GetAllApplies))          // For Employee
	rtr.GET("/api/applies/:apply_id", h.mw.IsAuthed(h.GetApplyByID)) // For User and Employee

	rtr.POST("/api/applies", h.mw.IsAuthed(h.CreateApply)) // For User

	rtr.GET("/api/applies/:apply_id/accept", h.mw.IsAuthed(h.AcceptApply)) // For Employee
	rtr.GET("/api/applies/:apply_id/deny", h.mw.IsAuthed(h.DenyApply))     // For Employee
}

// GetAllApplies
// Get All own applies for clients or get all applies of client for employee
// [DONE | NOT TESTED]
func (h *Handler) GetAllApplies(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, role := GetData(r)

	if userID == 0 {
		respondJSON(w, http.StatusUnauthorized, "not authorized")
	}

	clientUserID := r.URL.Query().Get("user_id")
	if clientUserID == "" && role == "user" { // For Users, if param user_id is empty
		applies, err := h.applyService.GetAllApplies(r.Context(), userID, true)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondJSON(w, http.StatusOK, applies)
	} else { // For employees
		ClientID, err := strconv.Atoi(clientUserID)
		if err != nil {
			return
		}
		applies, err := h.applyService.GetAllApplies(r.Context(), ClientID, false)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error())
		}
		respondJSON(w, http.StatusOK, applies)
	}

}

// GetApplyByID [DONE | NOT TESTED]
func (h *Handler) GetApplyByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	userID, role := GetData(r)

	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}

	applyID := p.ByName("apply_id")
	if applyID == "" {
		http.Error(w, "Missing apply ID", http.StatusBadRequest)
		return
	}

	ApplyID, err := strconv.Atoi(applyID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	apply, err := h.applyService.GetApplyByID(r.Context(), userID, ApplyID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if role == "user" {
	} else if role == "employee" {
	} // TODO Filter for role

	respondJSON(w, http.StatusOK, apply)
}

// CreateApply [DONE | NOT TESTED]
func (h *Handler) CreateApply(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, role := GetData(r)

	if role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var input ApplyInput
	err := pkg.GetFromBody(r.Body, &input)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	input.UserID = userID

	status, err := h.applyService.CreateApply(r.Context(), userID, input.ToMap())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, status, map[string]interface{}{"status": "created"})
}

// AcceptApply [DONE | NOT TESTED]
func (h *Handler) AcceptApply(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, role := GetData(r)

	if role != "employee" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	applyID := p.ByName("apply_id")
	if applyID == "" {
		http.Error(w, "missing apply ID", http.StatusBadRequest)
		return
	}

	ApplyID, err := strconv.Atoi(applyID)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.applyService.AcceptApply(r.Context(), userID, ApplyID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"status": "accepted"})
}

// DenyApply [DONE | NOT TESTED]
func (h *Handler) DenyApply(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Проверка прав
	UserID, role := GetData(r)
	if role != "employee" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Получаем ID заявки
	applyID := p.ByName("apply_id")
	if applyID == "" {
		http.Error(w, "Missing apply ID", http.StatusBadRequest)
		return
	}

	ApplyID, err := strconv.Atoi(applyID)
	if err != nil {
		return
	}

	// Обрабатываем отклонение
	err = h.applyService.DenyApply(r.Context(), UserID, ApplyID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"status": "denied"})
}
