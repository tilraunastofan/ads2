# Kbd

Keyboard key hint: `kbd.Kbd` and `kbd.Group` from `components/kbd`.

```templ
@kbd.Group() { @kbd.Kbd() { ⌘ } @kbd.Kbd() { K } }
```

- In tooltips, menus (`dropdownmenu.Shortcut`) and help text. Show both ⌘ and Ctrl variants where the audience is mixed.
