# Breadcrumb

Where the page sits in the hierarchy: `breadcrumb.Breadcrumb`, `List`, `Item`, `Link`, `Separator`, `Page`, `Ellipsis` from `components/breadcrumb`.

```templ
@breadcrumb.Breadcrumb() {
    @breadcrumb.List() {
        @breadcrumb.Item() { @breadcrumb.Link(breadcrumb.LinkProps{Href: "/apps"}) { Applications } }
        @breadcrumb.Separator()
        @breadcrumb.Item() { @breadcrumb.Page() { { app.Name } } }
    }
}
```

- Top of every detail page, above the page title (Arion places it above its serif hero).
- Add `hx-boost="true"` on the `<nav>` or body so links swap the page without a full reload.
