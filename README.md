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
$ go install github.com/masakurapa/gover-html@latest
```

## Usage

```sh
$ go test -coverprofile=coverage.out ./...
$ gover-html
$ open coverage.html
```

## Options

### -i

Coverage profile for input. default is `coverage.out`.

### -input-files

Specify multiple input coverage profiles.

If "input-files" is specified, "input" value is not used.

You can specify multiple items separated by commas.

### -o

File for output. default is `coverage.html`.

### -include

Output only the specified file or directory.

You can specify multiple items separated by commas.

**If "exclude" is also specified, "exclude" option takes precedence.**

**Must specify a relative path from the package root directory!!**

#### example

In the case of the following directory tree

```
github.com/masakurapa/gover-html
└── internal
    ├── cover
    │     ├── cover.go
    │     └── func.go
    ├── html
    │     ├── html.go
    │     ├── template.go
    │     └── tree
    │         └── tree.go
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

Exclude specified files or directories from output.

You can specify multiple items separated by commas.

**If "include" is also specified, this option takes precedence.**

**Must specify a relative path from the package root directory!!**

#### example

In the case of the following directory tree

```
github.com/masakurapa/gover-html
└── internal
    ├── cover
    │     ├── cover.go
    │     └── func.go
    ├── html
    │     ├── html.go
    │     ├── template.go
    │     └── tree
    │         └── tree.go
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

### -exclude-func

Exclude specified function from output.

You can specify multiple items separated by commas.

#### example

```sh
# Output excluding the all `Example` function!!
$ gover-html -exclude-func 'Example'

# Output from the `internal/cover` package without the `Example` function!!
$ gover-html -exclude-func '(internal/cover).Example'

# Output from the `internal/cover/cover.go` without the `Example` function!!
$ gover-html -exclude-func '(internal/cover/cover.go).Example'

# Output except for the `Example` function in the `Sample` structure of the `internal/cover` package!!
$ gover-html -exclude-func '(internal/cover.Sample).Example'

# If you specify a structure, wildcards can be used in the path.
$ gover-html -exclude-func '(*.Sample).Example'

# When specifying more than one
$ gover-html -exclude 'Example,(internal/cover).Example,(internal/cover.Sample).Example'
```

### -theme
HTML color theme to output (`dark` or `light`)

if the value is invalid, it will be `dark`.

## Config file
`gover-html` looks for config files in the following paths from the current working directory.

- `.gover.yml`

See [.gover.yml](https://github.com/masakurapa/gover-html/blob/master/.gover.yml) for a example configuration file.
