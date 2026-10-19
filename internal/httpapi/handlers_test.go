package httpapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"catalog/internal/catalog"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := httptest.NewServer(NewRouter(catalog.NewStore(catalog.SeedProducts()), logger))
	t.Cleanup(srv.Close)
	return srv
}

func getJSON(t *testing.T, srv *httptest.Server, path string, wantStatus int, into any) {
	t.Helper()
	res, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	defer res.Body.Close()
	if res.StatusCode != wantStatus {
		t.Fatalf("GET %s: status = %d, want %d", path, res.StatusCode, wantStatus)
	}
	if into != nil {
		if err := json.NewDecoder(res.Body).Decode(into); err != nil {
			t.Fatalf("GET %s: decode body: %v", path, err)
		}
	}
}

func TestIndex(t *testing.T) {
	srv := newTestServer(t)

	var got IndexResponse
	getJSON(t, srv, "/", http.StatusOK, &got)

	if got.Count != 3 || len(got.Products) != 3 {
		t.Fatalf("count = %d, len(products) = %d, want 3 and 3", got.Count, len(got.Products))
	}
	if book := got.Products[1]; book.Price != (Money{Amount: 71.92, Currency: "RON"}) {
		t.Errorf("book price = %+v, want discounted 71.92 RON", book.Price)
	}
	if tee := got.Products[2]; tee.InStock {
		t.Errorf("tee in_stock = true, want false")
	}
}

func TestProduct(t *testing.T) {
	srv := newTestServer(t)

	t.Run("discounted book in EUR", func(t *testing.T) {
		var got ProductResponse
		getJSON(t, srv, "/products/2?currency=eur&qty=2", http.StatusOK, &got)

		if got.UnitPrice != (Money{Amount: 14.47, Currency: "EUR"}) {
			t.Errorf("unit_price = %+v", got.UnitPrice)
		}
		if got.Total != (Money{Amount: 28.94, Currency: "EUR"}) {
			t.Errorf("total = %+v", got.Total)
		}
		if got.ListPrice == nil || *got.ListPrice != (Money{Amount: 18.09, Currency: "EUR"}) {
			t.Errorf("list_price = %+v", got.ListPrice)
		}
		if got.DiscountPercent != 20 || !got.Available {
			t.Errorf("discount_percent = %d, available = %v", got.DiscountPercent, got.Available)
		}
	})

	t.Run("laptop in HUF", func(t *testing.T) {
		var got ProductResponse
		getJSON(t, srv, "/products/1?currency=HUF", http.StatusOK, &got)

		if got.UnitPrice != (Money{Amount: 278074, Currency: "HUF"}) {
			t.Errorf("unit_price = %+v", got.UnitPrice)
		}
		if got.ListPrice != nil {
			t.Errorf("list_price = %+v, want omitted", got.ListPrice)
		}
	})

	t.Run("out of stock tee is still returned", func(t *testing.T) {
		var got ProductResponse
		getJSON(t, srv, "/products/3", http.StatusOK, &got)

		if got.Available || got.Stock != 0 {
			t.Errorf("available = %v, stock = %d", got.Available, got.Stock)
		}
		if got.Quantity != 1 || got.UnitPrice.Currency != "RON" {
			t.Errorf("quantity = %d, currency = %s, want defaults 1 and RON", got.Quantity, got.UnitPrice.Currency)
		}
	})
}

func TestProductErrors(t *testing.T) {
	srv := newTestServer(t)

	tests := []struct {
		path   string
		status int
		code   string
	}{
		{"/products/abc", http.StatusBadRequest, "invalid_id"},
		{"/products/0", http.StatusBadRequest, "invalid_id"},
		{"/products/1?currency=USD", http.StatusBadRequest, "invalid_currency"},
		{"/products/1?qty=0", http.StatusBadRequest, "invalid_qty"},
		{"/products/1?qty=11", http.StatusBadRequest, "invalid_qty"},
		{"/products/1?qty=two", http.StatusBadRequest, "invalid_qty"},
		{"/products/99", http.StatusNotFound, "product_not_found"},
		{"/products/99?currency=USD", http.StatusBadRequest, "invalid_currency"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			var got ErrorResponse
			getJSON(t, srv, tt.path, tt.status, &got)
			if got.Error.Code != tt.code {
				t.Errorf("error.code = %q, want %q", got.Error.Code, tt.code)
			}
		})
	}
}

func TestRouting(t *testing.T) {
	srv := newTestServer(t)

	var health HealthResponse
	getJSON(t, srv, "/healthz", http.StatusOK, &health)
	if health.Status != "ok" {
		t.Errorf("healthz status = %q", health.Status)
	}

	getJSON(t, srv, "/does-not-exist", http.StatusNotFound, nil)

	res, err := http.Post(srv.URL+"/", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("POST / status = %d, want 405", res.StatusCode)
	}
}
