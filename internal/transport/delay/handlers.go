package delay

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"minibank-backend/internal/entity"
	"minibank-backend/pkg/auth"
	"net/http"
	"strconv"
)

type Handler struct {
	mw           auth.MiddleWare
	dealService  IDealService
	delayService IDelayService
}

// NewHandler constructor Handler struct
func NewHandler(
	mw auth.MiddleWare,
	dealService IDealService,
	delayService IDelayService,

) *Handler {
	return &Handler{
		mw:           mw,
		dealService:  dealService,
		delayService: delayService,
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

type IDealService interface {
	GetDeals(ctx context.Context, UserID int) ([]entity.EDeal, error)
	GetDealByID(ctx context.Context, UserID int, dealID int) (*entity.EDeal, error)
	InsertPay(ctx context.Context, UserID int, dealID int, paymentData map[string]interface{}) error
}

type IDelayService interface {
	CreateDelay(ctx context.Context, UserID int, dealID int, delayData map[string]interface{}) (*entity.EDelay, error)
	GetDelays(ctx context.Context, UserID int, dealID int) ([]entity.EDelay, error)
}

func (h *Handler) Register(rtr *httprouter.Router) {
	// Delay

	rtr.POST("/api/deals/:deal_id/delays", h.mw.IsAuthed(h.CreateDelay))
	rtr.GET("/api/deals/:deal_id/delays", h.mw.IsAuthed(h.GetDelays))
}

// CreateDelay [DONE | NOT TESTED]
func (h *Handler) CreateDelay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, role := GetData(r)

	if role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	dealID := p.ByName("deal_id")
	if dealID == "" {
		http.Error(w, "missing deal ID", http.StatusBadRequest)
		return
	}
	DealID, err := strconv.Atoi(dealID)
	if err != nil {
		return
	}

	var input DelayInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	input.EmployeeID = userID

	delay, err := h.delayService.CreateDelay(r.Context(), userID, DealID, input.ToMap())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, delay)
}

// GetDelays [DONE | NOT TESTED]
func (h *Handler) GetDelays(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	UserID, _ := GetData(r)
	dealID := p.ByName("deal_id")
	if dealID == "" {
		http.Error(w, "Missing deal ID", http.StatusBadRequest)
		return
	}
	DealID, err := strconv.Atoi(dealID)
	if err != nil {
		return
	}

	// Получаем задержки
	delays, err := h.delayService.GetDelays(r.Context(), UserID, DealID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, delays)
}
