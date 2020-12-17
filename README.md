# gover-html

![github pages](https://github.com/masakurapa/gover-html/workflows/github%20pages/badge.svg)

This is a tool for outputting Golang coverage in HTML format.

## Features

|   | go cover | gover-html|
|---|---|---|
| file list | dropdown list | show directory tree to sidebar |
| coverage by package | no | show in sidebar |
| coverage by function | exec `go tool cover` with `-func` option | show at the top of the report |
| line number | no | yes |

## Example
- [masakurapa/gover-html (dark theme)](https://masakurapa.github.io/gover-html/gover-html_dark.html)
- [masakurapa/gover-html (light theme)](https://masakurapa.github.io/gover-html/gover-html_light.html)
- [go tool cover](https://masakurapa.github.io/gover-html/go-tool-cover.html)

## Installing

```sh
$ go get -u github.com/masakurapa/gover-html
```

## Usage

```sh
$ go test -coverprofile=coverage.out ./...
$ gover-html
$ open coverage.html
```

## Options
### -i
coverage profile for input. default is `coverage.out`.

### -o
file for output. default is `coverage.html`.

### -filter
output only the specified directories.

multiple directories can be specified separated by commas.

**Must specify a relative path from the package root directory!!**

#### example
In the case of the following directory tree

```
github.com/masakurapa/gover-html
└── internal
    ├── cover
    │   ├── cover.go
    │   └── func.go
    ├── html
    │   ├── html.go
    │   ├── template.go
    │   └── tree
    │       └── tree.go
    └── profile
        └── profile.go
```

You can filter the output packages by specifying values as follows

```sh
# Output only `internal/cover` package!!
$ gover-html -filter 'internal/cover'

# Output only `internal/cover` and `internal/html/tree` package!!
$ gover-html -filter './internal/cover,internal/html/tree/'

# This is no good :-C
$ gover-html -filter 'cover'
$ gover-html -filter 'html/tree'
```

### -theme
HTML color theme to output (`dark` or `light`)

if the value is invalid, it will be `dark`.
