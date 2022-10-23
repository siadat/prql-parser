package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/siadat/prql-parser/parser"
)

func main() {
	var p = parser.NewParser()
	var got, parseErr = p.Parse(os.Stdin)
	if parseErr != nil {
		fmt.Println(parseErr)
		return
	}
	fmt.Printf("%# v\n", pretty.Formatter(got))
}
