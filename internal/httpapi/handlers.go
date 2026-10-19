package httpapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"catalog/internal/catalog"
)

const (
	defaultQty = 1
	maxQty     = 10
)

type handlers struct {
	store *catalog.Store
}

func money(bani int64, c catalog.Currency) Money {
	return Money{Amount: catalog.Convert(bani, c), Currency: string(c)}
}

func (h *handlers) index(w http.ResponseWriter, _ *http.Request) {
	products := h.store.List()
	resp := IndexResponse{
		Count:    len(products),
		Products: make([]ProductSummary, 0, len(products)),
	}
	for _, p := range products {
		resp.Products = append(resp.Products, ProductSummary{
			ID:       p.ID,
			Slug:     p.Slug,
			Name:     p.Name,
			Category: string(p.Category),
			Price:    money(p.UnitPrice(), catalog.RON),
			InStock:  p.Stock > 0,
			URL:      fmt.Sprintf("/products/%d", p.ID),
		})
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *handlers) product(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid_id", "id must be a positive integer")
		return
	}

	query := r.URL.Query()

	currency := catalog.RON
	if raw := query.Get("currency"); raw != "" {
		if currency, err = catalog.ParseCurrency(raw); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_currency", "currency must be one of RON, EUR, HUF")
			return
		}
	}

	qty := defaultQty
	if raw := query.Get("qty"); raw != "" {
		if qty, err = strconv.Atoi(raw); err != nil || qty < 1 || qty > maxQty {
			writeError(w, http.StatusBadRequest, "invalid_qty", fmt.Sprintf("qty must be an integer between 1 and %d", maxQty))
			return
		}
	}

	p, err := h.store.Get(id)
	if errors.Is(err, catalog.ErrNotFound) {
		writeError(w, http.StatusNotFound, "product_not_found", fmt.Sprintf("product %d does not exist", id))
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "internal server error")
		return
	}

	unit := p.UnitPrice()
	resp := ProductResponse{
		ID:          p.ID,
		Slug:        p.Slug,
		Name:        p.Name,
		Description: p.Description,
		Category:    string(p.Category),
		UnitPrice:   money(unit, currency),
		Quantity:    qty,
		Total:       money(unit*int64(qty), currency),
		Stock:       p.Stock,
		Available:   p.Stock >= qty,
		Rating:      p.Rating,
		Specs:       p.Specs,
	}
	if p.DiscountPercent > 0 {
		list := money(p.PriceMinor, currency)
		resp.ListPrice = &list
		resp.DiscountPercent = p.DiscountPercent
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *handlers) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}
