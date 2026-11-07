package icon

import (
	"strings"
	"testing"
)

// Lucide's data keeps a newline and indent between elements; the rendered
// SVG must not, or View source shows stray "<path" lines (multi-element icons).
func TestIcon_NoWhitespaceBetweenElements(t *testing.T) {
	for _, name := range []string{"folder-git-2", "sun", "x"} {
		var b strings.Builder
		if err := Icon(name)().Render(t.Context(), &b); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		svg := b.String()
		if strings.ContainsAny(svg, "\n\t") || strings.Contains(svg, ">  <") || strings.Contains(svg, "> <") {
			t.Errorf("%s has whitespace between elements: %q", name, svg)
		}
		if !strings.Contains(svg, `<path d=`) {
			t.Errorf("%s lost its paths: %q", name, svg)
		}
	}
}
