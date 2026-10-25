# DropdownMenu

A menu of actions on one object: `dropdownmenu.DropdownMenu`, `dropdownmenu.Trigger(ctx)`, `Content{Side, Align}`, `Label`, `Item{Variant, Href, Disabled}`, `CheckboxItem`, `RadioGroup`, `RadioItem`, `Separator`, `Shortcut`, `Sub*` from `components/dropdownmenu`.

```templ
@dropdownmenu.DropdownMenu() {
    @button.Button(button.Props{Variant: button.VariantGhost, Size: button.SizeIcon, Attributes: dropdownmenu.Trigger(ctx)}) { @icon.Ellipsis() }
    @dropdownmenu.Content(dropdownmenu.ContentProps{Align: dropdownmenu.AlignEnd}) {
        @dropdownmenu.Item(dropdownmenu.ItemProps{Attributes: templ.Attributes{"hx-post": restartURL, "hx-swap": "none"}}) { @icon.RotateCcw() Restart pods }
        @dropdownmenu.Separator()
        @dropdownmenu.Item(dropdownmenu.ItemProps{Variant: dropdownmenu.ItemVariantDestructive}) { @icon.Trash2() Delete application }
    }
}
```

- Row actions in tables (ghost icon trigger, `Align: End`), and the overflow of a split `buttongroup`.
- Destructive items last, after a separator, and they open an `alertdialog`.
- Solid `popover` background: do not enable shadcn's translucent menus.
