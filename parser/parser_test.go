package parser_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/siadat/prql-parser/ast"
	"github.com/siadat/prql-parser/parser"
	"github.com/siadat/prql-parser/scanner"
)

const IgnorePos = -1

func TestParser(tt *testing.T) {
	var testCases = []struct {
		src  string
		want *ast.Root
	}{
		{
			src: `from table1`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.FromTransform{
						Table: ast.Ident{Name: "table1", Pos: 5},
					},
				},
			},
		},
		{
			src: `from e1 = table1`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.FromTransform{
						Alias: &ast.Ident{Name: "e1", Pos: 5},
						Table: ast.Ident{Name: "table1", Pos: 10},
					},
				},
			},
		},
		{
			// TODO: skip extra newlines
			src: `from table1
			select column1`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.FromTransform{
						Table: ast.Ident{Name: "table1", Pos: 5},
					},
					ast.SelectTransform{
						List: ast.List{
							Items: []ast.Expr{
								ast.Column{Name: ast.Ident{Name: "column1", Pos: IgnorePos}},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		var p = parser.NewParser()
		var src = tc.src
		var got, err = p.Parse(strings.NewReader(src))
		if err != nil {
			tt.Fatalf("test case failed src=%q: %v", src, err)
		}

		var cmpOpt = cmp.FilterValues(func(p1, p2 scanner.Pos) bool { return p1 == IgnorePos || p2 == IgnorePos || p1 == p2 }, cmp.Ignore())

		if diff := cmp.Diff(tc.want, got, cmpOpt); diff != "" {
			tt.Fatalf("case failed src=%q (-want +got):\n%s", src, diff)
		}
	}
}
