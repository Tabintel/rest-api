package rateservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type RateService interface{
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
		return 
	}

	rate , err := h.rateService.RateService(r.Context(), data.CurrencyPair)
	if err != nil {
		return
	}

	// return rate
}