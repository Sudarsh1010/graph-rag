package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func salesOrderRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "sales_order_headers")
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.getEntity(w, r, "sales_order_headers", "sales_order")
	})
}

func billingDocumentRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "billing_document_headers")
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.getEntity(w, r, "billing_document_headers", "billing_document")
	})
}

func deliveryRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "outbound_delivery_headers")
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.getEntity(w, r, "outbound_delivery_headers", "delivery_document")
	})
}

func customerRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "business_partners")
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.getEntity(w, r, "business_partners", "business_partner")
	})
}

func productRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "products")
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.getEntity(w, r, "products", "product")
	})
}

func plantRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "plants")
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.getEntity(w, r, "plants", "plant")
	})
}

func paymentRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "payments_accounts_receivable")
	})
}

func journalEntryRoutes(r chi.Router, h *handler) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h.listEntities(w, r, "journal_entry_items_accounts_receivable")
	})
}

func graphRoutes(r chi.Router, h *handler) {
	r.Post("/query", h.graphQuery)
}

func searchRoutes(r chi.Router, h *handler) {
	r.Get("/", h.searchEntities)
}
