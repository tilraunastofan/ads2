package patterns

import (
	"context"
	"strings"
	"testing"
)

func TestLabelsFollowPageLanguage(t *testing.T) {
	for lang, want := range map[string]string{"": "Dark theme", "en": "Dark theme", "is": "Dökkt þema", "is-IS": "Dökkt þema", "de": "Dark theme"} {
		if got := labelsFor(lang).DarkTheme; got != want {
			t.Errorf("labelsFor(%q) = %q, want %q", lang, got, want)
		}
	}
}

func TestIcelandicShellHasNoEnglishChrome(t *testing.T) {
	var b strings.Builder
	if err := Shell(ShellProps{Title: "Yfirlit", Lang: "is"}).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	html := b.String()
	for _, want := range []string{`lang="is"`, "Fara beint í efni", `aria-label="Dökkt þema"`, `aria-label="Sýna eða fela valmynd"`} {
		if !strings.Contains(html, want) {
			t.Errorf("Icelandic shell lacks %s", want)
		}
	}
	for _, en := range []string{"Skip to content", `aria-label="Dark theme"`, `aria-label="Toggle sidebar"`} {
		if strings.Contains(html, en) {
			t.Errorf("Icelandic shell contains English %s", en)
		}
	}
}
