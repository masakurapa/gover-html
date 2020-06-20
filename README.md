# go-cover

```
$ go test -coverpkg ./...  -coverprofile=coverage.out ./test/...
$ go tool cover -html=coverage.out -o coverage.html
$ open coverage.html
```

source
https://github.com/golang/go/tree/0104a31b8fbcbe52728a08867b26415d282c35d2/src/cmd/cover

https://stackoverflow.com/questions/31413281/golang-coverprofile-output-format
```
name.go:line.column,line.column numberOfStatements count
```
