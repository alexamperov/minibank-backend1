package admin

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
	adminService IAdminService
}
type IAdminService interface {
	GetEmployees(ctx context.Context) ([]entity.EUser, error)
	PaySalary(ctx context.Context, UserID int, Sum int) error
	GetEmployee(ctx context.Context, EmployeeID int) (entity.EUser, error)
}

// respondJSON - func for send status-code and response data
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func GetData(r *http.Request) (int, string) {
	userRole := r.Context().Value("user_role").(string)
	userID := r.Context().Value("user_id").(int)

	return userID, userRole
}

// Register TODO Grant Role endpoint
func (h *Handler) Register(rtr *httprouter.Router) {
	rtr.GET("/api/admin/employees", h.mw.IsAuthed(h.GetEmployees))
	rtr.POST("/api/admin/employees/:employee_id/salary", h.mw.IsAuthed(h.PaySalary))
	rtr.GET("/api/admin/employees/:employee_id", h.mw.IsAuthed(h.GetEmployeeByID))
}

func (h *Handler) GetEmployees(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Проверка прав доступа
	_, userRole := GetData(r)
	if userRole != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Получение списка сотрудников
	employees, err := h.adminService.GetEmployees(r.Context())
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	respondJSON(w, http.StatusOK, employees)
}

func (h *Handler) PaySalary(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	_, userRole := GetData(r)
	if userRole != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Парсинг ID сотрудника
	employeeID, err := strconv.Atoi(p.ByName("employee_id"))
	if err != nil || employeeID <= 0 {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Парсинг тела запроса
	var request struct {
		Sum int `json:"sum"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{})
		return
	}

	// Выплата зарплаты
	if err := h.adminService.PaySalary(r.Context(), employeeID, request.Sum); err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetEmployeeByID NOT TESTED
func (h *Handler) GetEmployeeByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Проверка прав доступа
	_, userRole := GetData(r)
	if userRole != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Парсинг ID сотрудника
	employeeID, err := strconv.Atoi(p.ByName("employee_id"))
	if err != nil || employeeID <= 0 {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	// Получение информации о сотруднике
	employee, err := h.adminService.GetEmployee(r.Context(), employeeID)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, employee)
}
