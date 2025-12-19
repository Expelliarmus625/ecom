package products

import (
	"log"
	"net/http"

	"github.com/expelliarmus625/ecom/internal/json"
	"github.com/go-chi/chi/v5"
)

type Handler struct{
	service Service
}

func NewHandler(service Service) *Handler{
	return &Handler{
		service: service, 
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {	

  products, err := h.service.ListProducts(r.Context()) 
	if err != nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, products, http.StatusOK)
}

func (h *Handler) FindProductByID(w http.ResponseWriter, r *http.Request) {	
	id := chi.URLParam(r, "id")

	log.Default().Printf("%s", id)
  product, err := h.service.FindProductByID(r.Context(), id)
	if err != nil{
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, product, http.StatusOK)
}
