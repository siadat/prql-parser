package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kr/pretty"
	"github.com/siadat/prql-parser/ast"
	"github.com/siadat/prql-parser/parser"
	"github.com/siadat/prql-parser/scanner"
	"github.com/siadat/prql-parser/token"
)

func TestParseExpression(tt *testing.T) {
	var testCases = []struct {
		want ast.Expr
		src  string
	}{
		{
			src:  `123`,
			want: ast.Integer{Value: 123},
		},
		{
			src: `1 * 2`,
			want: ast.BinaryExpr{
				X:  ast.Integer{Value: 1},
				Y:  ast.Integer{Value: 2},
				Op: token.MUL,
			},
		},
		{
			src: `1 + 2 * 3 * 4 + 5 # == 1 + ((2 * (3 * 4)) + 5)`,
			want: ast.BinaryExpr{
				X: ast.Integer{Value: 1},
				Y: ast.BinaryExpr{
					X: ast.BinaryExpr{
						X: ast.Integer{Value: 2},
						Y: ast.BinaryExpr{
							X:  ast.Integer{Value: 3},
							Y:  ast.Integer{Value: 4},
							Op: token.MUL,
						},
						Op: token.MUL,
					},
					Y:  ast.Integer{Value: 5},
					Op: token.ADD,
				},
				Op: token.ADD,
			},
		},
		{
			src: `1 * 2 + 3 + 4 * 5 # == (1 * 2) + (3 + (4 * 5))`,
			want: ast.BinaryExpr{
				X: ast.BinaryExpr{
					X:  ast.Integer{Value: 1},
					Y:  ast.Integer{Value: 2},
					Op: token.MUL,
				},
				Y: ast.BinaryExpr{
					X: ast.Integer{Value: 3},
					Y: ast.BinaryExpr{
						X:  ast.Integer{Value: 4},
						Y:  ast.Integer{Value: 5},
						Op: token.MUL,
					},
					Op: token.ADD,
				},
				Op: token.ADD,
			},
		},
	}

	for _, tc := range testCases {
		var p = parser.NewParser()
		p.SetDebug(true)
		var src = tc.src
		var got, err = p.ParseExpr(strings.NewReader(src))
		src = formatSrc(src, true)
		if err != nil {
			tt.Fatalf("test case failed\nsrc:\n%s\nerr: %v", src, err)
		}

		var cmpOpt = cmp.FilterValues(func(p1, p2 scanner.Pos) bool { return p1 == IgnorePos || p2 == IgnorePos || p1 == p2 }, cmp.Ignore())

		if diff := cmp.Diff(tc.want, got, cmpOpt); diff != "" {
			fmt.Printf("got: %# v\n", pretty.Formatter(got))
			tt.Fatalf("mismatching results\nsrc:\n%s\ndiff guide:\n  - want\n  + got\ndiff:\n%s", src, diff)
		}
	}
}
