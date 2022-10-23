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
			src: `from table1 # comment`,
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
			src: "from table1\n\n \n  \n# comment1 \n   # comment2 \n \t select column1 # comment3",
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.FromTransform{
						Table: ast.Ident{Name: "table1", Pos: 5},
					},
					ast.SelectTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.Column{Name: ast.Ident{Name: "column1", Pos: IgnorePos}},
							},
						},
					},
				},
			},
		},
		{
			src: `select [column1, column2]`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.SelectTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.Column{Name: ast.Ident{Name: "column1", Pos: IgnorePos}},
								ast.Column{Name: ast.Ident{Name: "column2", Pos: IgnorePos}},
							},
						},
					},
				},
			},
		},
		{
			src: `select [
			  column1,
			  column2,
			  123,
			  1.23,
			  "hello world",
			  @2022-12-31,
			  @01:02:03,
			  @2022-12-31T01:02:03,
			  123seconds
			]`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.SelectTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.Column{Name: ast.Ident{Name: "column1", Pos: IgnorePos}},
								ast.Column{Name: ast.Ident{Name: "column2", Pos: IgnorePos}},
								ast.Integer{Value: 123},
								ast.Float{Value: 1.23},
								ast.String{Value: `"hello world"`},
								ast.Date{Year: 2022, Month: 12, Day: 31},
								ast.Time{Hour: 1, Minute: 2, Second: 3},
								ast.Timestamp{Year: 2022, Month: 12, Day: 31, Hour: 1, Minute: 2, Second: 3},
								ast.Interval{Count: 123, Unit: "seconds"},
							},
						},
					},
				},
			},
		},
		{
			src: `select [
			  column1,
			  column2, # trailing comma
			]`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.SelectTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.Column{Name: ast.Ident{Name: "column1", Pos: IgnorePos}},
								ast.Column{Name: ast.Ident{Name: "column2", Pos: IgnorePos}},
							},
						},
					},
				},
			},
		},
		{
			src: `derive x = 5`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.DeriveTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.AssignExpr{Name: "x", Expr: ast.Integer{Value: 5}},
							},
						},
					},
				},
			},
		},
		{
			src: `
			select [
			  1, 1 * 2, # 2 expressions in one line
			  +1 + -2.1, # signed numbers
			  expr1 = 1 + 2 * 3 * 4 + 5 # == 1 + ((2 * (3 * 4)) + 5),
			  expr2 = 1 * 2 + 3 + 4 * 5 # == (1 * 2) + (3 + (4 * 5)),
			]
			`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.SelectTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.Integer{Value: 1},
								ast.BinaryExpr{
									X:  ast.Integer{Value: 1},
									Y:  ast.Integer{Value: 2},
									Op: token.MUL,
								},
								ast.BinaryExpr{
									X:  ast.Integer{Value: +1},
									Y:  ast.Float{Value: -2.1},
									Op: token.ADD,
								},
								ast.AssignExpr{
									Name: "expr1",
									Expr: ast.BinaryExpr{
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
								ast.AssignExpr{
									Name: "expr2",
									Expr: ast.BinaryExpr{
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
							},
						},
					},
				},
			},
		},
		{
			src: `
			select [
			  column1,
			  x - 1,
			  1 - x,
			  (1),
			  (1 + 2),
			  y + (1),
			  (1) + x,
			  z = ((z*2) + 1),
			]
			`,
			want: &ast.Root{
				Transforms: []ast.Node{
					ast.SelectTransform{
						List: ast.ExprList{
							Items: []ast.Expr{
								ast.Column{Name: ast.Ident{Name: "column1", Pos: IgnorePos}},
								ast.BinaryExpr{
									X:  ast.Column{Name: ast.Ident{Name: "x", Pos: IgnorePos}},
									Y:  ast.Integer{Value: 1},
									Op: token.SUB,
								},
								ast.BinaryExpr{
									X:  ast.Integer{Value: 1},
									Y:  ast.Column{Name: ast.Ident{Name: "x", Pos: IgnorePos}},
									Op: token.SUB,
								},
								ast.ParenExpr{
									X: ast.Integer{Value: 1},
								},
								ast.ParenExpr{
									X: ast.BinaryExpr{
										X:  ast.Integer{Value: 1},
										Y:  ast.Integer{Value: 2},
										Op: token.ADD,
									},
								},
								ast.BinaryExpr{
									X:  ast.Column{Name: ast.Ident{Name: "y", Pos: IgnorePos}},
									Y:  ast.ParenExpr{X: ast.Integer{Value: 1}},
									Op: token.ADD,
								},
								ast.BinaryExpr{
									X:  ast.ParenExpr{X: ast.Integer{Value: 1}},
									Y:  ast.Column{Name: ast.Ident{Name: "x", Pos: IgnorePos}},
									Op: token.ADD,
								},
								ast.AssignExpr{
									Name: "z",
									Expr: ast.ParenExpr{
										X: ast.BinaryExpr{
											X: ast.ParenExpr{
												X: ast.BinaryExpr{
													X:  ast.Column{Name: ast.Ident{Name: "z", Pos: IgnorePos}},
													Y:  ast.Integer{Value: 2},
													Op: token.MUL,
												},
											},
											Y:  ast.Integer{Value: 1},
											Op: token.ADD,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		var p = parser.NewParser()
		p.SetDebug(true)
		var src = tc.src
		var got, err = p.Parse(strings.NewReader(src))
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

func formatSrc(src string, showWhitespaces bool) string {
	var prefix = "   | "
	if showWhitespaces {
		// src = strings.ReplaceAll(src, " ", "₋")
		src = strings.ReplaceAll(src, "\t", "␣")
		src = strings.Join(strings.Split(src, "\n"), "⏎\n"+prefix)
		src = prefix + src + "·"
		return src
	}
	return src
}
