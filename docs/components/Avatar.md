# Avatar

A person or service identity: `avatar.Avatar{Size}`, `Image{Src, Alt}`, `Fallback`, `Badge`, `Group`, `GroupCount` from `components/avatar`.

```templ
@avatar.Avatar(avatar.Props{Size: avatar.SizeSm}) {
    @avatar.Image(avatar.ImageProps{Src: u.PhotoURL, Alt: u.Name})
    @avatar.Fallback() { { u.Initials } }
}
```

- Always provide a `Fallback` with initials (Icelandic letters included: Þ, Æ, Ð). On-prem photo services are often unreachable.
- `Group` + `GroupCount` for approvers and on-call rotations.
