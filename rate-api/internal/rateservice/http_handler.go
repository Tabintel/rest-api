package rateservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type RateService interface {
	RateService(ctx context.Context, currencyPair string) (*RateResponse, error)
}

// NewHTTPHandler returns a new HTTP handler for the service.
func NewHTTPHandler(rateService RateService) http.Handler {
	h := &httpHandler{
		rateService: rateService,
	}
	router := chi.NewRouter()

	router.Post("/", h.GetRate)
	return router
}

type httpHandler struct {
	rateService RateService
}

func (h *httpHandler) GetRate(w http.ResponseWriter, r *http.Request) {
	data := &RateRequest{}
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, err := h.rateService.RateService(r.Context(), data.CurrencyPair)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the rate response to JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
