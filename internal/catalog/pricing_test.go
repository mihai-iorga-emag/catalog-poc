package catalog

import "testing"

func TestUnitPriceRoundsDiscountDown(t *testing.T) {
	p := Product{PriceMinor: 999, DiscountPercent: 15}
	// 999 * 0.85 = 849.15 bani
	if got := p.UnitPrice(); got != 849 {
		t.Fatalf("UnitPrice() = %d, want 849", got)
	}
}

func TestUnitPriceWithoutDiscount(t *testing.T) {
	p := Product{PriceMinor: 5900}
	if got := p.UnitPrice(); got != 5900 {
		t.Fatalf("UnitPrice() = %d, want 5900", got)
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		bani int64
		to   Currency
		want float64
	}{
		{8990, RON, 89.90},
		{8990, EUR, 18.09},
		{8990, HUF, 7143},
		{349999, EUR, 704.2},
	}
	for _, tt := range tests {
		if got := Convert(tt.bani, tt.to); got != tt.want {
			t.Errorf("Convert(%d, %s) = %v, want %v", tt.bani, tt.to, got, tt.want)
		}
	}
}

func TestParseCurrency(t *testing.T) {
	for in, want := range map[string]Currency{"RON": RON, "eur": EUR, " huf ": HUF} {
		got, err := ParseCurrency(in)
		if err != nil || got != want {
			t.Errorf("ParseCurrency(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "USD", "EURO"} {
		if _, err := ParseCurrency(in); err == nil {
			t.Errorf("ParseCurrency(%q) succeeded, want error", in)
		}
	}
}
