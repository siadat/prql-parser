test:
	go list ./... | while read -r pkg; do \
		go test -failfast -count=1 -v "$$pkg" || exit 1; \
	done
	@echo "All tests passed"

prql.pest:
	wget --quiet 'https://raw.githubusercontent.com/prql/prql/main/prql-compiler/src/prql.pest'

main.wasm:
	GOOS=js GOARCH=wasm go build -o main.wasm ./cmd/prql-parser/

clean:
	rm main.wasm
