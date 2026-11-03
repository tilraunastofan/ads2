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
            @chart.Tooltip(chart.TooltipProps{Content: chart.TooltipContentProps{Formatter: tooltipRow}})
            @chart.Legend()
            @chart.Bar(chart.BarProps{DataKey: "deploys", Fill: "var(--chart-1-quiet)", Radius: []float64{4, 4, 0, 0}}) {
                for i := range rows {
                    @chart.Cell(chart.CellProps{Fill: weekFill(i, len(rows))}) // chart-1 for the latest
                }
                @chart.LabelList(chart.LabelListProps{DataKey: "latest", Position: "top", Class: "fill-foreground font-mono"})
            }
            @chart.Bar(chart.BarProps{DataKey: "rollbacks", Fill: "var(--chart-4)", Radius: []float64{4, 4, 0, 0}})
        }
    }
}
```

- Numbers: the chart's own tooltip formats with `toLocaleString("en-US")` (`1,284`). Always pass a tooltip `Formatter` that returns `patterns.ChartTooltipRow(name, color, patterns.ChartNumber(v, 0))` (`deployTooltipRow` in `ui/pages/metrics.templ`), and format ticks and labels with `patterns.ChartNumber`, which takes the `any` values the chart hands over (`int` from Go data, `float64` from JSON-shaped data).
- Bars: `Radius: []float64{4, 4, 0, 0}`, rounded at the data end and square on the baseline. A single number rounds all four corners.
- Colours: the leading series in `chart-1`, other series in neutral `chart-4` / `chart-5`; `chart-2` (gold) and `chart-3` (plum) only for series of equal weight. Set them on each series as `Fill` / `Stroke: "var(--chart-N)"`; the legend and tooltip swatches follow. Leave `Config.Color` empty under a strict CSP, and never use hex.
- Highlight one period with `chart.Cell` children on the leading bar: full `var(--chart-1)` for it and `var(--chart-1-quiet)` for the rest (the showcase's `weekFill`). Give the series itself `chart-1-quiet`, so its legend swatch matches most of its bars, label the highlighted bar directly with a `LabelList` whose data key is empty for every other row, and say so in the card description ("latest week highlighted and labelled").
- Tick labels mono: add `[&_.recharts-cartesian-axis-tick_text]:font-mono` to the `chart.Container` `Class`. The main chart of a page is 280–360px tall (`h-[300px]` on Metrics).
- Always a legend or direct labels, and a card title that names the unit.
- Put charts in a `card`; refresh a chart by swapping its card with htmx (`hx-get` + `hx-trigger`), never by client-side data fetching.
- Up to five series; beyond that, use a table.
- Every chart has its data as a table too, for screen readers and exact values: a collapsed `accordion` under the chart with a `table.Table` (`sr-only` caption, numbers right-aligned in mono). See `WeeklyTable` in `ui/pages/metrics.templ`.
- The preview draws the chart statically; the component renders the same SVG on the server.
