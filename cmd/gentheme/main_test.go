package main

import (
	"os"
	"strings"
	"testing"
)

const tokenPath = "../../tokens/tokens.json"

func TestContrastRequirements(t *testing.T) {
	raw, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatal(err)
	}
	_, tk, err := generate(raw, "tokens/tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range checkContrast(tk) {
		t.Error(f)
	}
}

// The committed theme must be what the tokens generate.
func TestThemeIsCurrent(t *testing.T) {
	raw, _ := os.ReadFile(tokenPath)
	css, _, err := generate(raw, "tokens/tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	cur, err := os.ReadFile("../../assets/css/arion-theme.css")
	if err != nil {
		t.Fatal(err)
	}
	if string(cur) != css {
		t.Error("assets/css/arion-theme.css is out of date: run `task theme`")
	}
}

func TestMissingDarkValueFails(t *testing.T) {
	doc := `{"color":{"themes":[{"id":"light"},{"id":"dark"}],"tokens":[{"name":"x","value":{"light":"#ffffff"}}]}}`
	if _, _, err := generate([]byte(doc), "x"); err == nil || !strings.Contains(err.Error(), `"dark"`) {
		t.Errorf("want an error naming the dark theme, got %v", err)
	}
}

func TestEmptyThemesFail(t *testing.T) {
	if _, _, err := generate([]byte(`{"color":{"themes":[],"tokens":[]}}`), "x"); err == nil {
		t.Error("want an error for empty themes")
	}
}

func TestRatio(t *testing.T) {
	black, white := [3]float64{}, [3]float64{255, 255, 255}
	if r := ratio(black, white); r < 20.9 || r > 21.1 {
		t.Errorf("black on white = %.2f, want 21", r)
	}
}
