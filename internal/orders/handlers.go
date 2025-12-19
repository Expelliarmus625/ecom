package orders

import (
	"log"
	"net/http"

	"github.com/expelliarmus625/ecom/internal/json"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.ListOrders(r.Context())
	if err != nil {
		http.Error(w, err.Error() ,http.StatusInternalServerError)
	}
	json.Write(w, orders, http.StatusOK)
}

func (h *Handler) ListOrderItems(w http.ResponseWriter, r *http.Request) {
	orderId := chi.URLParam(r, "id")
	orderItems, err := h.service.ListOrderItems(r.Context(), orderId)
	if err != nil {
		http.Error(w, err.Error() ,http.StatusInternalServerError)
	}
	json.Write(w, orderItems, http.StatusOK)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	//Call Orders service to place the actual order
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil{
		log.Println(err)

		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}	
	//Send http response
	json.Write(w, createdOrder, http.StatusCreated)
}
