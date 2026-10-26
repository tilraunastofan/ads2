package patterns

import (
	"testing"
	"time"
)

func d(y int, m time.Month, day int) time.Time { return time.Date(y, m, day, 0, 0, 0, 0, time.UTC) }

func TestDates(t *testing.T) {
	at := time.Date(2026, 10, 23, 14, 32, 0, 0, time.UTC)
	cases := []struct{ got, want string }{
		{Date(d(2026, 10, 23), "is"), "23. október 2026"},
		{Date(d(2026, 10, 23), "en"), "23 Oct 2026"},
		{Date(d(2026, 5, 1), "is-IS"), "1. maí 2026"},
		{Date(d(2026, 8, 9), "en-GB"), "9 Aug 2026"},
		{DateTime(at, "is"), "23. október 2026 kl. 14:32"},
		{DateTime(at, "en"), "23 Oct 2026, 14:32"},
		{DateRange(d(2026, 10, 20), d(2026, 10, 24), "en"), "20–24 Oct 2026"},
		{DateRange(d(2026, 10, 20), d(2026, 10, 24), "is"), "20.–24. október 2026"},
		{DateRange(d(2026, 9, 28), d(2026, 10, 3), "en"), "28 Sep – 3 Oct 2026"},
		{DateRange(d(2026, 9, 28), d(2026, 10, 3), "is"), "28. september – 3. október 2026"},
		{DateRange(d(2025, 12, 29), d(2026, 1, 4), "en"), "29 Dec 2025 – 4 Jan 2026"},
		{DateRange(d(2026, 10, 23), d(2026, 10, 23), "en"), "23 Oct 2026"},
		{DateRange(d(2026, 10, 23), time.Time{}, "is"), "23. október 2026"},
		{ISODate(d(2026, 10, 3)), "2026-10-03"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
	p, err := ParseISODate("2026-10-03", time.UTC)
	if err != nil || !p.Equal(d(2026, 10, 3)) {
		t.Errorf("ParseISODate = %v, %v", p, err)
	}
}
