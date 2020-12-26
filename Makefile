build:
	go build -o gover-html gover.go

gotest:
	go test -count=1 ./internal/...

cover:
	go test -coverprofile=coverage.out ./internal/...

bench:
	go test -bench=. ./internal/... -benchmem
