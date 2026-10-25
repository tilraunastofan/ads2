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

- Internal tools use the sidebar shell; public pages use a top bar.
- Active item = `sidebar-accent` fill with `sidebar-accent-foreground` (Arion's blue-on-pale-blue active nav).
- The header shows the Arion star mark (`Arion SVG/logo.svg`, served from `/assets/img/arion-mark.svg` in the app) at 32px, never on a `sidebar-primary` tile.
- Boost menu links (`hx-boost="true"` on `sidebar.Inset` target) so only the inset swaps.
