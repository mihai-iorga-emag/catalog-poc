package catalog

import (
	"cmp"
	"errors"
	"slices"
)

var ErrNotFound = errors.New("product not found")

// Store is a read-only, in-memory catalog. Nothing mutates it after
// construction, so it is safe for concurrent use.
type Store struct {
	products []Product
}

func NewStore(products []Product) *Store {
	sorted := slices.Clone(products)
	slices.SortFunc(sorted, func(a, b Product) int { return cmp.Compare(a.ID, b.ID) })
	return &Store{products: sorted}
}

// List returns every product, ordered by ID.
func (s *Store) List() []Product {
	return slices.Clone(s.products)
}

func (s *Store) Get(id int) (Product, error) {
	for _, p := range s.products {
		if p.ID == id {
			return p, nil
		}
	}
	return Product{}, ErrNotFound
}

// SeedProducts is the fixed catalog the PoC serves.
func SeedProducts() []Product {
	return []Product{
		{
			ID:          1,
			Slug:        "nimbus-14-laptop",
			Name:        "Nimbus 14 Laptop",
			Description: "14-inch ultralight laptop for everyday work.",
			Category:    CategoryElectronics,
			PriceMinor:  349999,
			Stock:       12,
			Rating:      4.6,
			Specs: LaptopSpecs{
				Brand:          "Nimbus",
				CPU:            "8-core 3.2 GHz",
				RAMGB:          16,
				StorageGB:      512,
				ScreenInches:   14,
				WarrantyMonths: 24,
			},
		},
		{
			ID:              2,
			Slug:            "the-pragmatic-gopher",
			Name:            "The Pragmatic Gopher",
			Description:     "A practical guide to writing maintainable Go services.",
			Category:        CategoryBooks,
			PriceMinor:      8990,
			DiscountPercent: 20,
			Stock:           3,
			Rating:          4.8,
			Specs: BookSpecs{
				Author:   "Jane Gopher",
				Pages:    312,
				ISBN:     "978-0-00-000000-2",
				Language: "en",
				Format:   "paperback",
			},
		},
		{
			ID:          3,
			Slug:        "basic-cotton-tee",
			Name:        "Basic Cotton Tee",
			Description: "Regular-fit crew neck t-shirt.",
			Category:    CategoryFashion,
			PriceMinor:  5900,
			Stock:       0,
			Rating:      4.1,
			Specs: ApparelSpecs{
				Sizes:    []string{"S", "M", "L", "XL"},
				Color:    "navy",
				Material: "100% cotton",
			},
		},
	}
}
