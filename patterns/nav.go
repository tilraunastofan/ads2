package patterns

import "strings"

// activeHref is the Href of the item the current path belongs to: the longest Href,
// among all items and their sub-items, that is the path or a path-segment prefix of
// it ("/apps" for "/apps/payments", never for "/apps-old"). "/" matches only itself,
// or it would match every page. "" when nothing matches.
func activeHref(groups []NavGroup, cur string) string {
	best := ""
	consider := func(href string) {
		if len(href) > len(best) && within(cur, href) {
			best = href
		}
	}
	for _, g := range groups {
		for _, it := range g.Items {
			consider(it.Href)
			for _, c := range it.Children {
				consider(c.Href)
			}
		}
	}
	return best
}

func within(cur, href string) bool {
	switch {
	case href == "":
		return false
	case cur == href:
		return true
	case href == "/":
		return false
	}
	return strings.HasPrefix(cur, strings.TrimSuffix(href, "/")+"/")
}

// itemActive: the item is the active one, or one of its sub-items is. A collapsed
// sidebar hides sub-items, so the parent shows where the user is.
func itemActive(it NavItem, active string) bool {
	if active == "" {
		return false
	}
	if it.Href == active {
		return true
	}
	for _, c := range it.Children {
		if c.Href == active {
			return true
		}
	}
	return false
}
