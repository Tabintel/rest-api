package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tabintel/rest-api/rate-api/internal/config"
	"github.com/tabintel/rest-api/rate-api/internal/rateservice"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		// log error with fatal
	}

	rateService := rateservice.NewRate(cfg)
	rateHttpHandler := rateservice.NewHTTPHandler(rateService)

	router := chi.NewRouter()
	router.Mount("/v1/rate", rateHttpHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
