# Sidebar

The app shell navigation: `sidebar.Provider`, `Sidebar`, `Header`, `Content`, `Group`, `GroupLabel`, `Menu`, `MenuItem`, `MenuButton{Href, IsActive, Tooltip}`, `MenuBadge`, `Footer`, `Inset`, `Trigger` from `components/sidebar`.

**Consumer provides:** the menu tree, the active item (`IsActive` from the request path), `Collapsible` (`offcanvas`, `icon`, `none`), `Variant` (`sidebar`, `floating`, `inset`), and the page in `sidebar.Inset`.

```templ
@sidebar.Provider() {
    @sidebar.Sidebar(sidebar.Props{Collapsible: sidebar.CollapsibleIcon}) {
        @sidebar.Content() {
            @sidebar.Group() {
                @sidebar.GroupLabel() { Workloads }
                @sidebar.Menu() {
                    for _, it := range nav {
                        @sidebar.MenuItem() {
                            @sidebar.MenuButton(sidebar.MenuButtonProps{Href: it.Href, IsActive: it.Href == currentPath, Tooltip: it.Label}) {
                                @it.Icon
                                <span>{ it.Label }</span>
                            }
                        }
                    }
                }
            }
        }
    }
    @sidebar.Inset() { { children... } }
}
```

- Internal tools use the sidebar shell (`patterns.Shell` wires header, menu, inset and theme toggle); public pages use a top bar. Give `sidebar.Menu` `Class: "gap-1"` for breathing room between items (Shell does).
- Active item = `sidebar-accent` fill with `sidebar-accent-foreground` (Arion's blue-on-pale-blue active nav).
- The header shows the Arion star mark (`patterns.ArionMark("size-8")`, blue in light, white in dark) at 32px, never on a `sidebar-primary` tile.
- The collapsed state survives navigation only if the server renders it: `sidebar.js` stores it in the `sidebar_state` cookie, and wrapping your handler in `patterns.SidebarState` makes `patterns.Shell` render that state. Without the middleware every page load, including each boosted link, expands the sidebar again.
- Collapsed (`CollapsibleIcon`) the rail is `4rem`, not shadcn's `3rem`, so the 32px mark fits with its padding, and the menu icons are centred under it. `patterns.Shell` sets `--sidebar-width-icon: 4rem` on the Provider and `group-data-[collapsible=icon]:items-center` on the menu; do the same if you build the shell by hand.
- Boost menu links (`hx-boost="true"` on `sidebar.Inset` target) so only the inset swaps.
