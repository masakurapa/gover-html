build:
	go build -o gover-html gover.go

test:
	go test -count=1 ./internal/...

cover:
	go test -coverprofile=coverage.out ./internal/...

bench:
	go test -bench=. ./internal/... -benchmem

make_covers: build
	${MAKE} cover
	rm -rf covers
	mkdir covers
	./gover-html -o covers/dark_theme.html -theme dark
	./gover-html -o covers/light_theme.html -theme light
	./gover-html -o covers/include_1.html -include "internal/cover"
	./gover-html -o covers/include_2.html -include "internal/cover,./internal/html/tree/"
	./gover-html -o covers/include_3.html -include "internal/cover/cover.go"
	./gover-html -o covers/exclude_1.html -exclude "internal/cover"
	./gover-html -o covers/exclude_2.html -exclude "internal/cover,./internal/html/tree/"
	./gover-html -o covers/exclude_3.html -exclude "internal/cover/cover.go"
	./gover-html -o covers/exclude_func_1.html -exclude-func "IsOutputTarget"
	./gover-html -o covers/exclude_func_2.html -exclude-func "(internal/cover/filter).IsOutputTarget"
	./gover-html -o covers/exclude_func_3.html -exclude-func "(internal/cover/filter.filter).IsOutputTargetFunc"
	./gover-html -o covers/exclude_func_4.html -exclude-func "(internal/cover/filter/filter.go).IsOutputTarget"
	./gover-html -o covers/exclude_func_5.html -exclude-func "(internal/cover/filter/filter.go.filter).IsOutputTargetFunc"

.PHONY: test
