package patterns

import (
	"strconv"
	"time"
)

// Dates follow the page language: Icelandic pages write "23. október 2026",
// English pages "23 Oct 2026". Ranges share what they can and join with an en
// dash: "20.–24. október 2026" / "20–24 Oct 2026". Times are always 24-hour.
// lang is the page's BCP 47 tag; anything starting with "is" is Icelandic,
// everything else English.

var isMonths = [...]string{"janúar", "febrúar", "mars", "apríl", "maí", "júní", "júlí", "ágúst", "september", "október", "nóvember", "desember"}

func icelandic(lang string) bool { return len(lang) >= 2 && lang[:2] == "is" }

func day(t time.Time) string { return strconv.Itoa(t.Day()) }

func monthName(t time.Time, lang string) string {
	if icelandic(lang) {
		return isMonths[t.Month()-1]
	}
	return t.Format("Jan")
}

// Date formats a day: "23. október 2026" (is) or "23 Oct 2026" (en).
func Date(t time.Time, lang string) string {
	if icelandic(lang) {
		return day(t) + ". " + monthName(t, lang) + " " + strconv.Itoa(t.Year())
	}
	return day(t) + " " + monthName(t, lang) + " " + strconv.Itoa(t.Year())
}

// DateTime formats a day and a 24-hour time: "23. október 2026 kl. 14:32" (is)
// or "23 Oct 2026, 14:32" (en).
func DateTime(t time.Time, lang string) string {
	if icelandic(lang) {
		return Date(t, lang) + " kl. " + t.Format("15:04")
	}
	return Date(t, lang) + ", " + t.Format("15:04")
}

// DateRange formats an inclusive range, sharing the month and year when it can.
// A zero to (an open range) or to == from formats as a single date.
func DateRange(from, to time.Time, lang string) string {
	if to.IsZero() || sameDay(from, to) {
		return Date(from, lang)
	}
	is := icelandic(lang)
	switch {
	case from.Year() == to.Year() && from.Month() == to.Month():
		if is {
			return day(from) + ".–" + Date(to, lang)
		}
		return day(from) + "–" + Date(to, lang)
	case from.Year() == to.Year():
		if is {
			return day(from) + ". " + monthName(from, lang) + " – " + Date(to, lang)
		}
		return day(from) + " " + monthName(from, lang) + " – " + Date(to, lang)
	}
	return Date(from, lang) + " – " + Date(to, lang)
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

// ISODate is the wire format for dates in forms and URLs: "2026-10-23".
func ISODate(t time.Time) string { return t.Format(time.DateOnly) }

// ParseISODate reads ISODate's format in the given location (time.Local if nil).
func ParseISODate(s string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.Local
	}
	return time.ParseInLocation(time.DateOnly, s, loc)
}
