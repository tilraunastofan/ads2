package patterns

// labels are the few words the patterns themselves render (skip link, icon-button
// names), in the page language, so an Icelandic page never reads "Toggle theme".
type labels struct {
	SkipToContent string
	DarkTheme     string // name of the theme toggle; aria-pressed carries on/off
	ToggleSidebar string
}

var (
	labelsEN = labels{SkipToContent: "Skip to content", DarkTheme: "Dark theme", ToggleSidebar: "Toggle sidebar"}
	labelsIS = labels{SkipToContent: "Fara beint í efni", DarkTheme: "Dökkt þema", ToggleSidebar: "Sýna eða fela valmynd"}
)

// labelsFor returns the labels for lang ("is", "is-IS", "en"…), English otherwise,
// choosing the language the same way Date does.
func labelsFor(lang string) labels {
	if icelandic(lang) {
		return labelsIS
	}
	return labelsEN
}
