package chart

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// nest renders parent with kids as its { children... }.
func nest(parent templ.Component, kids ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			for _, k := range kids {
				if err := k.Render(ctx, w); err != nil {
					return err
				}
			}
			return nil
		})
		return parent.Render(templ.WithChildren(ctx, body), w)
	})
}

func renderLine(t *testing.T, line LineProps) string {
	t.Helper()
	config := Config{{Key: "now", Label: "Today"}} // no Color: a strict CSP blocks its <style>
	c := nest(Container(ContainerProps{ID: "t", Config: config}),
		nest(LineChart(LineChartProps{Data: []Datum{{"h": "1", "now": 3}, {"h": "2", "now": 5}}}),
			Legend(),
			Line(line),
		),
	)
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// The legend swatch of a Line follows its Stroke, like Recharts'
// getLegendItemColor, so a chart without Config colours still has one.
func TestLineLegendSwatchFollowsStroke(t *testing.T) {
	for _, tc := range []struct {
		name, stroke, want string
	}{
		{"stroke", "var(--chart-1)", "background-color:var(--chart-1)"},
		{"no stroke", "", "background-color:var(--color-now)"},
		{"stroke none", "none", "background-color:var(--color-now)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			html := renderLine(t, LineProps{DataKey: "now", Stroke: tc.stroke})
			if !strings.Contains(html, tc.want) {
				t.Errorf("legend swatch: want %q in\n%s", tc.want, html)
			}
		})
	}
}
