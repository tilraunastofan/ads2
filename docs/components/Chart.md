# Chart

Server-rendered SVG charts with shadcn's chart API: `chart.Container{Config}`, `BarChart`, `LineChart`, `AreaChart`, `PieChart`, `RadarChart` with `CartesianGrid`, `XAxis`, `YAxis`, `Bar`, `Line`, `Area`, `Tooltip`, `Legend` from `components/chart` (`shadcn-templ add chart`).

**Consumer provides:** `[]chart.Datum` rows (`map[string]any`), a `chart.Config` naming each series with a label and a colour, and the series components.

```templ
// Labels only: colours go on the series. A Config Color makes the container emit an
// inline <style>, which the showcase's CSP (style-src 'self') blocks.
var deployConfig = chart.Config{
    {Key: "deploys", Label: "Deployments"},
    {Key: "rollbacks", Label: "Rollbacks"},
}

templ DeployChart(rows []chart.Datum) {
    @chart.Container(chart.ContainerProps{ID: "deploys", Config: deployConfig, Class: "h-[200px] w-full"}) {
        @chart.BarChart(chart.BarChartProps{Data: rows, AccessibilityLayer: true}) {
            @chart.CartesianGrid(chart.CartesianGridProps{Vertical: chart.Bool(false)})
            @chart.XAxis(chart.XAxisProps{DataKey: "week", TickMargin: 8})
            @chart.YAxis(chart.YAxisProps{Width: 32, TickFormatter: tickInt})
            @chart.Tooltip()
            @chart.Legend()
            @chart.Bar(chart.BarProps{DataKey: "deploys", Fill: "var(--chart-1)", Radius: 4})
            @chart.Bar(chart.BarProps{DataKey: "rollbacks", Fill: "var(--chart-2)", Radius: 4})
        }
    }
}
```

- `TickFormatter` takes `any`: Go ints arrive as `int`, JSON-shaped data as `float64`. `tickInt` in `ui/pages/metrics.templ` handles both with `patterns.Int` / `patterns.Dec`.
- Colours: `chart-1` … `chart-5` in order (blue, gold, plum, slate, light slate), set on each series as `Fill` / `Stroke: "var(--chart-N)"`; the legend and tooltip swatches follow. Leave `Config.Color` empty under a strict CSP, and never use hex.
- Always a legend or direct labels, and a card title that names the unit. Format ticks and tooltips with `patterns.Dec` / `Pct` (Arion format).
- Put charts in a `card`; refresh a chart by swapping its card with htmx (`hx-get` + `hx-trigger`), never by client-side data fetching.
- Up to five series; beyond that, use a table.
- The preview draws the chart statically; the component renders the same SVG on the server.
