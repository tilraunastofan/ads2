package patterns

import (
	"math"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct{ got, want string }{
		{Int(0), "0"},
		{Int(999), "999"},
		{Int(1000), "1.000"},
		{Int(1284500), "1.284.500"},
		{Int(-1284500), "−1.284.500"},
		{Dec(2.3, 1), "2,3"},
		{Dec(1284.5, 1), "1.284,5"},
		{Dec(-0.8, 1), "−0,8"},
		{Dec(-0.04, 1), "0,0"}, // rounds to zero: no minus
		{Dec(1234, 0), "1.234"},
		{Signed(9.5, 1), "+9,5"},
		{Signed(-0.8, 1), "−0,8"},
		{Signed(0, 1), "0,0"},
		{Pct(2.3, 1), "2,3%"},
		{SignedPct(18, 0), "+18%"},
		{SignedPct(-0.8, 1), "−0,8%"},
		{Int(math.MinInt64), "−9.223.372.036.854.775.808"},
		// Beyond int64: grouped as text, not clamped to math.MaxInt64.
		{Dec(1e20, 0), "100.000.000.000.000.000.000"},
		{Dec(-1e20, 1), "−100.000.000.000.000.000.000,0"},
		{Signed(1e19, 0), "+10.000.000.000.000.000.000"},
		// Not a number (a change from zero): an en dash, never "0" or "−0".
		{Dec(math.NaN(), 1), NotAvailable},
		{Dec(math.Inf(1), 0), NotAvailable},
		{Dec(math.Inf(-1), 0), NotAvailable},
		{Signed(math.Inf(1), 1), NotAvailable},
		{Signed(math.NaN(), 1), NotAvailable},
		{Pct(math.NaN(), 1), NotAvailable},
		{SignedPct(math.Inf(-1), 0), NotAvailable},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}
