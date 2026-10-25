# Item

A flexible list row with media, text and actions: `item.Item{Variant, Size, Href}`, `Group`, `Media{Variant}`, `Content`, `Title`, `Description`, `Actions`, `Header`, `Footer`, `Separator` from `components/item`.

```templ
@item.Group() {
    for _, f := range reports {
        @item.Item(item.Props{Variant: item.VariantOutline, Href: f.URL}) {
            @item.Media(item.MediaProps{Variant: item.MediaVariantIcon}) { @icon.FileText() }
            @item.Content() { @item.Title() { { f.Name } } @item.Description() { { f.Meta } } }
        }
    }
}
```

- Document lists (reports, statements), settings rows, search results, activity feeds.
- `VariantOutline` on `background`, `VariantMuted` inside cards.
