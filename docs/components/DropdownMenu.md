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
- Destructive items last, after a separator, and they open an `alertdialog`. Render the dialog next to the menu with an `ID` and give the item `templ.Attributes{"data-base-ui-click-trigger": true, "data-templ-controls": id, "aria-haspopup": "dialog"}`: the dialogs' `TriggerFor(id)` minus `data-slot`, `type` and `tabindex`, which would override the item's own (the showcase's `opens`). A dialog that needs more input (Scale…) works the same way.
- Solid `popover` background: do not enable shadcn's translucent menus.
