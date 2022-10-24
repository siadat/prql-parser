package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/siadat/prql-parser/parser"
	"github.com/siadat/prql-parser/scanner"
)

func TestErrors(tt *testing.T) {
	var testCases = []struct {
		src  string
		want string
	}{
		{
			src:  `from table1, table2`,
			want: `failed to parse a transform, unexpected COMMA(",") at 11`,
		},
		{
			src:  `select 1 + ++1`, // extra plus sign
			want: `expected integer or float, got ADD("+") at 12`,
		},
		{
			src:  `select 1 + --1`, // extra minus sign
			want: `expected integer or float, got SUB("-") at 12`,
		},
		{
			src: `
			from table1
			select [1, a b]
			`,
			want: `unexpected token IDENTIFIER("b") at 32`,
		},
	}

	for _, tc := range testCases {
		var p = parser.NewParser()
		p.SetDebug(true)
		var src = tc.src
		var got, gotErr = p.Parse(strings.NewReader(src))
		fmt.Println("Parse got", gotErr)
		src = formatSrc(src, true)
		if gotErr == nil {
			tt.Fatalf("expected an error, got nil\nsrc:\n%s\ngot:\n%#v", src, got)
		}

		var cmpOpt = cmp.FilterValues(func(p1, p2 scanner.Pos) bool { return p1 == IgnorePos || p2 == IgnorePos || p1 == p2 }, cmp.Ignore())

		if diff := cmp.Diff(tc.want, gotErr.Error(), cmpOpt); diff != "" {
			tt.Fatalf("mismatching errors\nsrc:\n%s\ndiff guide:\n  - want\n  + got\ndiff:\n%s", src, diff)
		}
	}
}
