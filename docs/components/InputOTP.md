# InputOTP

One-time-code entry: `inputotp.InputOTP{MaxLength, Pattern}`, `Group`, `Slot{Index}`, `Separator` from `components/inputotp` (`shadcn-templ add input-otp`). One real `<input autocomplete="one-time-code">` sits behind the slots, so paste and SMS autofill work.

**Consumer provides:** `MaxLength`, one `Slot` per digit, a `name` through `Attributes` (they land on the input), and a label or description.

```templ
<form hx-post="/approvals/42" hx-target="this" hx-swap="outerHTML" hx-status:422="target:this">
    @field.Field() {
        @field.Label(field.LabelProps{For: "otp"}) { Approval code }
        @inputotp.InputOTP(inputotp.Props{ID: "otp", MaxLength: 6, Pattern: "^[0-9]+$",
            Attributes: templ.Attributes{"name": "code"}}) {
            @inputotp.Group() { @inputotp.Slot(inputotp.SlotProps{Index: 0}) @inputotp.Slot(inputotp.SlotProps{Index: 1}) @inputotp.Slot(inputotp.SlotProps{Index: 2}) }
            @inputotp.Separator()
            @inputotp.Group() { @inputotp.Slot(inputotp.SlotProps{Index: 3}) @inputotp.Slot(inputotp.SlotProps{Index: 4}) @inputotp.Slot(inputotp.SlotProps{Index: 5}) }
        }
        @field.Description() { Enter the 6-digit code from your authenticator. }
    }
    @button.Button(button.Props{Type: button.TypeSubmit}) { Approve }
</form>
```

- Six digits in two groups of three. Submit with the button or Enter in the field.
- No auto-submit under the showcase CSP: `hx-trigger="input-otp-change[detail.value.length==6]"` is a trigger filter, htmx compiles it with `new Function`, and `script-src 'self'` blocks that. Without a filter the event also fires at 1–5 digits. If you need auto-submit, ship a small script that calls `form.requestSubmit()` when `e.detail.value.length` reaches `MaxLength`.
- Wrong code: return 422 with the form re-rendered, the input cleared, and `field.Error`.
- Invalid state: set `aria-invalid` to the **string** `"true"`. A Go `true` renders a bare `aria-invalid`, which ARIA reads as false and Tailwind's `aria-invalid:` variant ignores. Set it on the input (`Attributes`) and on each `Slot` (`SlotProps.Attributes`): the input is a sibling of the groups, so the slots' red border does not follow it. The showcase's `ariaInvalid(msg)` helper returns both.
- Under a strict CSP the component's runtime `<style>` (autofill and selection hiding, iOS and password-manager rules) is blocked and its rules are silently dropped. Typing, paste and the slots still work; copy the rules from `components/inputotp/inputotp.js` (`insertStyles`) into your app CSS if you need them.
- Use for step-up approval (production promotion, payment release), never as the only login factor UI.
