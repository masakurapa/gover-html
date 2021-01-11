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

### -include
output only the specified file or directory.

you can specify multiple items separated by commas.

**if "exclude" is also specified, "exclude" option takes precedence.**

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

You can filter the output by specifying values as follows.

```sh
# Output only `internal/cover`!!
$ gover-html -include 'internal/cover'

# Output only `internal/cover` and `internal/html/tree`!!
$ gover-html -include './internal/cover,internal/html/tree/'

# This is no good :-C
$ gover-html -include '/internal/cover'
$ gover-html -include 'cover'
$ gover-html -include 'html/tree'
```

### -exclude
output expect the specified file or directory.

You can specify multiple items separated by commas.

**if "include" is also specified, this option takes precedence.**

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

You can filter the output by specifying values as follows.

```sh
# Output excluding `internal/cover`!!
$ gover-html -exclude 'internal/cover'

# Output excluding `internal/cover` and `internal/html/tree`!!
$ gover-html -exclude './internal/cover,internal/html/tree/'

# This is no good :-C
$ gover-html -exclude '/internal/cover'
$ gover-html -exclude 'cover'
$ gover-html -exclude 'html/tree'
```

If you specify `include` and `exclude` at the same time, the output looks like this.

```sh
$ gover-html -include 'internal/html' -exclude 'internal/html/tree/'
github.com/masakurapa/gover-html
└── internal
    └── html
        ├── html.go
        └── template.go
```

### -theme
HTML color theme to output (`dark` or `light`)

if the value is invalid, it will be `dark`.

## Config file
`gover-html` looks for config files in the following paths from the current working directory.

- `.gover.yml`

See [.gover.yml](https://github.com/masakurapa/gover-html/blob/master/.gover.yml) for a example configuration file.

