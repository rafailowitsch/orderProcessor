package httphandler

import (
	"context"
	"encoding/json"
	"net/http"
	"orderProcessor/internal/domain"
	"strings"
)

type Repository interface {
	Read(ctx context.Context, orderUID string) (*domain.Order, error)
	ReadAll(ctx context.Context) ([]*domain.Order, error)
	//Create(ctx context.Context, order *domain.Order) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

//func (h *Handler) SaveOrder(w http.ResponseWriter, r *http.Request) {
//	var order domain.Order
//
//	err := json.NewDecoder(r.Body).Decode(&order)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	err = h.repo.Create(r.Context(), &order)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := strings.TrimPrefix(r.URL.Path, "/order/")
	if orderUID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	order, err := h.repo.Read(r.Context(), orderUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.ReadAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
