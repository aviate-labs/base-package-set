# Base Package Set

![latest version](https://img.shields.io/badge/moc-0.9.1-blue)

## How to use with `vessel`?

```dhall
-- vessel.dhall
{
  dependencies = [ "base-0.9.1" ],
  compiler = Some "0.9.1"
}
```

```dhall
-- package-set.dhall
let base = https://github.com/internet-computer/base-package-set/releases/download/moc-0.9.1/package-set.dhall sha256:eb7e8de9987ee129adcfeaf45d1d7d5363f2d206383724df71c0b7fe872eb437
in  base
```
