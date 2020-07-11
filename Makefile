build:
	go build -o go-cover cover.go

cover:
	go test -coverprofile=coverage.out ./internal/...
