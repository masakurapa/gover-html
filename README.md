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
- [masakurapa/gover-html](https://masakurapa.github.io/gover-html/gover-html.html)
- [go tool cover](https://masakurapa.github.io/gover-html/go-tool-cover.html)

## Installing

```sh
$ go get -u github.com/masakurapa/gover-html
```

## Usage

```
$ go test -coverprofile=coverage.out ./...
$ gover-html
$ open coverage.html
```

## Options
### -i
coverage profile for input. default is `coverage.out`.

### -o
file for output. default is `coverage.html`.
