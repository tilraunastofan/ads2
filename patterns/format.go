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
	neg := n < 0
	s := strconv.FormatUint(absU(n), 10)
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(r)
	}
	if neg {
		return minus + b.String()
	}
	return b.String()
}

func absU(n int64) uint64 {
	if n < 0 {
		return uint64(-(n + 1)) + 1
	}
	return uint64(n)
}

// Dec formats 1284.5 with one decimal as "1.284,5".
func Dec(f float64, decimals int) string {
	s := strconv.FormatFloat(math.Abs(f), 'f', decimals, 64)
	whole, frac, _ := strings.Cut(s, ".")
	n, _ := strconv.ParseInt(whole, 10, 64)
	out := Int(n)
	if frac != "" {
		out += "," + frac
	}
	if f < 0 && strings.Trim(s, "0.") != "" {
		return minus + out
	}
	return out
}

// Signed formats a change with an explicit sign: "+9,5" or "−0,8".
func Signed(f float64, decimals int) string {
	s := Dec(f, decimals)
	if strings.HasPrefix(s, minus) {
		return s
	}
	if strings.Trim(s, "0,.") == "" {
		return s
	}
	return "+" + s
}

// Pct formats a ratio already in percent: Pct(2.3, 1) is "2,3%".
func Pct(f float64, decimals int) string { return Dec(f, decimals) + "%" }

// SignedPct formats a change in percent: "+18%", "−0,8%".
func SignedPct(f float64, decimals int) string { return Signed(f, decimals) + "%" }
