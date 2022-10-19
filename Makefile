test:
	go test -count=1 -v ./...

prql.pest:
	wget --quiet 'https://raw.githubusercontent.com/prql/prql/main/prql-compiler/src/prql.pest'
