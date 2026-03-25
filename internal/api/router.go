package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/uptrace/bun"
	"time"
)

func NewRouter(db *bun.DB, graph GraphQuerier, nlq NLQueryHandler) http.Handler {
	r := chi.NewRouter()
	h := newHandler(db, graph, nlq)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/sales-orders", func(r chi.Router) { salesOrderRoutes(r, h) })
		r.Route("/billing-documents", func(r chi.Router) { billingDocumentRoutes(r, h) })
		r.Route("/deliveries", func(r chi.Router) { deliveryRoutes(r, h) })
		r.Route("/customers", func(r chi.Router) { customerRoutes(r, h) })
		r.Route("/products", func(r chi.Router) { productRoutes(r, h) })
		r.Route("/plants", func(r chi.Router) { plantRoutes(r, h) })
		r.Route("/payments", func(r chi.Router) { paymentRoutes(r, h) })
		r.Route("/journal-entries", func(r chi.Router) { journalEntryRoutes(r, h) })
		r.Route("/graph", func(r chi.Router) { graphRoutes(r, h) })
		r.Route("/search", func(r chi.Router) { searchRoutes(r, h) })
		r.Post("/nlquery", h.nlQueryHandler)
	})

	return r
}
