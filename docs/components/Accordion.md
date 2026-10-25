# Accordion

Collapsible sections: `accordion.Accordion`, `Item{Value}`, `Trigger`, `Content` from `components/accordion`.

**Consumer provides:** items with `Value`, trigger text, content; `DefaultValue []string`, `Multiple`.

```templ
@accordion.Accordion(accordion.Props{DefaultValue: []string{"promote"}}) {
    @accordion.Item(accordion.ItemProps{Value: "promote"}) {
        @accordion.Trigger() { How are releases promoted? }
        @accordion.Content() { <p>Kargo promotes a freight from dev to staging automatically.</p> }
    }
}
```

- FAQs, runbook steps, advanced form options. Not for primary content people must see.
