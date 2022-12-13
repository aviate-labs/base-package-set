# Base Package Set

![latest version](https://img.shields.io/badge/moc-0.7.4-blue)

## How to use with `vessel`?

```dhall
-- vessel.dhall
{
  dependencies = [ "base-0.7.4" ],
  compiler = Some "0.7.4"
}
```

```dhall
-- package-set.dhall
let base = https://github.com/internet-computer/base-package-set/releases/download/moc-0.7.4/package-set.dhall sha256:3a20693fc597b96a8c7cf8645fda7a3534d13e5fbda28c00d01f0b7641efe494
in  base
```
