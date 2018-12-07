External resources access from `init()` function may have negative side-effects like cross components initialization dependency and other initialization order issues.

For example, there is a components tree like this:
```
- pkg1
- pkg2
  + pkg3
```

`lintinit` will report `pkg1` usages (global variables or function calls) from `pk2` `init()` function but allow to use everything from `pkg3`.

# Installing

```bash
go get -u github.com/astavonin/lintinit
```

# How to run

`intinit [directories]` - runs on package in current directory or `directories` recursively.
