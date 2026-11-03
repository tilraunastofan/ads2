package patterns

import (
	"math"
	"strconv"
	"strings"
)

// Number formats follow Arion's convention in every UI language: "." groups
// thousands, "," is the decimal mark, losses use a real minus sign (U+2212).
// Words and units follow the page language ("41 min" on English pages).

const minus = "−"

// Int formats 1284500 as "1.284.500".
func Int(n int64) string {
	s := group(strconv.FormatUint(absU(n), 10))
	if n < 0 {
		return minus + s
	}
	return s
}

func absU(n int64) uint64 {
	if n < 0 {
		return uint64(-(n + 1)) + 1
	}
	return uint64(n)
}

// NotAvailable is what Dec, Signed, Pct and SignedPct return for NaN or ±Inf, such
// as a change from zero (100 * (now-prev) / 0): an en dash, never a made-up "0".
const NotAvailable = "–"

func finite(f float64) bool { return !math.IsNaN(f) && !math.IsInf(f, 0) }

// Dec formats 1284.5 with one decimal as "1.284,5". The digits are grouped as
// text, so any finite float64 formats exactly as strconv rounds it, however large.
func Dec(f float64, decimals int) string {
	if !finite(f) {
		return NotAvailable
	}
	s := strconv.FormatFloat(math.Abs(f), 'f', decimals, 64)
	whole, frac, _ := strings.Cut(s, ".")
	out := group(whole)
	if frac != "" {
		out += "," + frac
	}
	if f < 0 && strings.Trim(s, "0.") != "" {
		return minus + out
	}
	return out
}

// group puts a "." between every three digits of an unsigned decimal string.
func group(digits string) string {
	var b strings.Builder
	for i, r := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Signed formats a change with an explicit sign: "+9,5" or "−0,8".
func Signed(f float64, decimals int) string {
	s := Dec(f, decimals)
	if s == NotAvailable || strings.HasPrefix(s, minus) {
		return s
	}
	if strings.Trim(s, "0,.") == "" {
		return s
	}
	return "+" + s
}

// Pct formats a ratio already in percent: Pct(2.3, 1) is "2,3%".
func Pct(f float64, decimals int) string { return withUnit(Dec(f, decimals), "%") }

// SignedPct formats a change in percent: "+18%", "−0,8%".
func SignedPct(f float64, decimals int) string { return withUnit(Signed(f, decimals), "%") }

// withUnit appends unit, except to NotAvailable: "–", not "–%".
func withUnit(s, unit string) string {
	if s == NotAvailable {
		return s
	}
	return s + unit
}
