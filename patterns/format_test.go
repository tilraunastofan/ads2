package patterns

import "testing"

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
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}
