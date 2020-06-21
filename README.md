# go-cover

## Sample
- [masakurapa/go-cover](https://masakurapa.github.io/go-cover/index.html)
- [go tool cover](https://masakurapa.github.io/go-cover/coverage.html)

## Command
```
$ go test -coverpkg ./internal/...  -coverprofile=coverage.out ./test/...
$ go tool cover -html=coverage.out -o coverage.html
$ go run cover.go
```
