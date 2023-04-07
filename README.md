# Base Package Set

![latest version](https://img.shields.io/badge/moc-0.8.6-blue)

## How to use with `vessel`?

```dhall
-- vessel.dhall
{
  dependencies = [ "base-0.8.6" ],
  compiler = Some "0.8.6"
}
```

```dhall
-- package-set.dhall
let base = https://github.com/internet-computer/base-package-set/releases/download/moc-0.8.6/package-set.dhall sha256:4a7734568f1c7e5dfe91d2ba802c9b6f218d7836904dea5a999a3096f6ef0d3c
in  base
```
