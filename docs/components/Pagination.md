# Pagination

Page through long server-side lists: `pagination.Pagination`, `Content`, `Item`, `Link{Href, IsActive}`, `Previous`, `Next`, `Ellipsis` from `components/pagination`.

```templ
@pagination.Pagination(pagination.Props{Attributes: templ.Attributes{"hx-boost": "true", "hx-target": "#audit-table", "hx-select": "#audit-table"}}) {
    @pagination.Content() {
        @pagination.Item() { @pagination.Previous(pagination.PreviousProps{Href: page.PrevURL, Disabled: !page.HasPrev}) }
        for _, n := range page.Numbers {
            @pagination.Item() { @pagination.Link(pagination.LinkProps{Href: page.URL(n), IsActive: n == page.Current}) { { strconv.Itoa(n) } } }
        }
        @pagination.Item() { @pagination.Next(pagination.NextProps{Href: page.NextURL, Disabled: !page.HasNext}) }
    }
}
```

- Real `Href`s so it works without JS; htmx boosts them and swaps only the table.
- Put it under the `table`, right-aligned in a `card.Footer` or centered below.
