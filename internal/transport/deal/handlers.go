package deal

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
	// Deals

	// Employee can see deals of client by ClientID | Admin can see deals of Employee which he signed
	rtr.GET("/api/all-deals/:user_id", h.mw.IsAuthed(h.GetDealsOfUser)) // For Employee and Admin

	rtr.GET("/api/deals", h.mw.IsAuthed(h.GetDeals)) // For User and Employee

	// Client can see his deals | Employee can see deals which he accepted before
	rtr.GET("/api/deals/:deal_id", h.mw.IsAuthed(h.GetDealByID)) // For User and Employee

	// Pay
	rtr.POST("/api/deals/:deal_id/pay", h.mw.IsAuthed(h.InsertPay)) //For User Role

}

// DEAL HANDLERS

// GetDeals [DONE | NOT TESTED] [Need Filter]
func (h *Handler) GetDeals(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, role := GetData(r)
	deals, err := h.dealService.GetDeals(r.Context(), userID) // Get Deals of User or of Employee
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if role == "user" {
	} else if role == "employee" {
	} // TODO Filter for role

	respondJSON(w, http.StatusOK, deals)
}

// GetDealByID [DONE | NOT TESTED] [Need Filter]
func (h *Handler) GetDealByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userID, role := GetData(r)

	dealID := p.ByName("deal_id")
	if dealID == "" {
		http.Error(w, "missing deal ID", http.StatusBadRequest)
		return
	}

	DealID, err := strconv.Atoi(dealID)
	if err != nil {
		return
	}

	deal, err := h.dealService.GetDealByID(r.Context(), userID, DealID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	// TODO Add Pays into struct Deal

	if role == "user" {
	} else if role == "employee" {
	} // TODO Filter for role

	respondJSON(w, http.StatusOK, deal)
}

func (h *Handler) GetDealsOfUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

// PAY HANDLERS

// InsertPay [DONE | NOT TESTED]
func (h *Handler) InsertPay(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	var input PayInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	input.DealID = DealID

	if err := h.dealService.InsertPay(r.Context(), userID, DealID, input.ToMap()); err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
