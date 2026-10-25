# Empty

Empty state for a list, table or search: `empty.Empty`, `Header`, `Media{Variant}`, `Title`, `Description`, `Content` from `components/empty`.

```templ
@empty.Empty() {
    @empty.Header() {
        @empty.Media(empty.MediaProps{Variant: empty.MediaVariantIcon}) { @icon.Inbox() }
        @empty.Title() { No applications yet }
        @empty.Description() { Register a Git repository and Argo CD will create the first application. }
    }
    @empty.Content() { @button.Button(button.Props{Size: button.SizeSm}) { Add repository } }
}
```

- Say what is missing and the one next step. For "no search results" offer clearing filters.
