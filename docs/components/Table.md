# Table

Tabular data: `table.Table`, `Header`, `Body`, `Footer`, `Row`, `Head`, `Cell`, `Caption` from `components/table`.

**Consumer provides:** rows from the server, column headers, an optional caption (units!), optional footer totals.

```templ
@table.Table() {
    @table.Caption() { Key figures, m.kr. }
    @table.Header() { @table.Row() {
        @table.Head() { Item }
        @table.Head(table.HeadProps{Class: "text-right"}) { 3F 2026 }
    } }
    @table.Body(table.BodyProps{ID: "rows"}) {
        for _, r := range rows {
            @table.Row() {
                @table.Cell() { { r.Label } }
                @table.Cell(table.CellProps{Class: "text-right font-mono tabular-nums"}) { { isk(r.Value) } }
            }
        }
    }
}
```

- Put tables in a `card` on the `background` canvas. `patterns.DataCard` does this and insets the first and last columns, so cells need no `pl-4` / `pr-4`.
- Numbers: right-aligned, `font-mono tabular-nums` (`figure` style), Arion number format in every language (`14.902`, `+9,5%`; use `patterns.Int`, `Dec`, `Signed`), units in the caption, not every cell.
- Gains `text-positive` with `+`, losses `text-destructive` with a real minus `−`. Never colour alone. Identifiers and versions: `font-mono text-code`.
- Sorting, filtering and paging are server round trips: `hx-get` on headers/filters with `hx-target="#rows"`, `hx-push-url="true"`.
- Selected rows: `data-state="selected"` on `table.Row`.
