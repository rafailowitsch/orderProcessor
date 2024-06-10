package httphandler

import (
	"context"
	"encoding/json"
	_ "github.com/swaggo/http-swagger" // http-swagger middleware
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

// @Summary Get an order
// @Description Get order by ID
// @ID get-order-by-id
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /order/{id} [get]
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

// @Summary Get all orders
// @Description Get all orders
// @ID get-all-orders
// @Accept json
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 500 {object} map[string]string
// @Router /order [get]
func (h *Handler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.ReadAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
