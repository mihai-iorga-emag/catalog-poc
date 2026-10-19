package catalog

// Category decides which specs type a product carries.
type Category string

const (
	CategoryElectronics Category = "electronics"
	CategoryBooks       Category = "books"
	CategoryFashion     Category = "fashion"
)

// Product is a catalog entry. Money is kept in RON minor units (bani) so that
// discounts and totals are computed on integers; conversion to other
// currencies happens only when a response is built.
type Product struct {
	ID              int
	Slug            string
	Name            string
	Description     string
	Category        Category
	PriceMinor      int64 // list price in RON bani
	DiscountPercent int   // 0 when the product is not on promotion
	Stock           int
	Rating          float64
	Specs           any // LaptopSpecs, BookSpecs or ApparelSpecs, depending on Category
}

// UnitPrice is the price of one unit in RON bani after the discount.
// Integer division rounds the discounted price down.
func (p Product) UnitPrice() int64 {
	if p.DiscountPercent <= 0 {
		return p.PriceMinor
	}
	return p.PriceMinor * int64(100-p.DiscountPercent) / 100
}

type LaptopSpecs struct {
	Brand          string  `json:"brand"`
	CPU            string  `json:"cpu"`
	RAMGB          int     `json:"ram_gb"`
	StorageGB      int     `json:"storage_gb"`
	ScreenInches   float64 `json:"screen_inches"`
	WarrantyMonths int     `json:"warranty_months"`
}

type BookSpecs struct {
	Author   string `json:"author"`
	Pages    int    `json:"pages"`
	ISBN     string `json:"isbn"`
	Language string `json:"language"`
	Format   string `json:"format"` // paperback, hardcover or ebook
}

type ApparelSpecs struct {
	Sizes    []string `json:"sizes"`
	Color    string   `json:"color"`
	Material string   `json:"material"`
}
