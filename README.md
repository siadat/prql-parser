This is a WIP parser in Go for [PRQL](https://github.com/prql/prql).
Only a subset of the grammar is supported at the moment.

* Run in Go playground https://go.dev/play/p/AIUpm5vMMy1
* Run in terminal: `echo 'from table1' | go run ./cmd/prql-parser`
* See the tests:
  * [/parser/parser_test.go](/parser/parser_test.go)
  * [/parser/expression_test.go](/parser/expression_test.go)
  * [/parser/errors_test.go](/parser/errors_test.go)
  * [/scanner/scanner_test.go](/scanner/scanner_test.go)
* Send a PR :)
