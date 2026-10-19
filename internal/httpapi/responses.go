package httpapi

import (
	"encoding/json"
	"net/http"
)

type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type ProductSummary struct {
	ID       int    `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    Money  `json:"price"`
	InStock  bool   `json:"in_stock"`
	URL      string `json:"url"`
}

type IndexResponse struct {
	Count    int              `json:"count"`
	Products []ProductSummary `json:"products"`
}

type ProductResponse struct {
	ID              int     `json:"id"`
	Slug            string  `json:"slug"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	UnitPrice       Money   `json:"unit_price"`
	ListPrice       *Money  `json:"list_price,omitempty"`
	DiscountPercent int     `json:"discount_percent,omitempty"`
	Quantity        int     `json:"quantity"`
	Total           Money   `json:"total"`
	Stock           int     `json:"stock"`
	Available       bool    `json:"available"`
	Rating          float64 `json:"rating"`
	Specs           any     `json:"specs"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// The status line is already sent; an encoding error cannot be reported to the client.
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}
