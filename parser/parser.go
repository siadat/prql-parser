package parser

import (
	"fmt"
	"io"

	"github.com/siadat/prql-parser/ast"
	"github.com/siadat/prql-parser/scanner"
	"github.com/siadat/prql-parser/token"
)

type Parser struct {
	src     io.Reader
	root    *ast.Root
	scanner *scanner.Scanner
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(src io.Reader) (*ast.Root, error) {
	p.scanner = scanner.NewScanner(src)
	p.scanner.SetSkipWhitespace(true)
	var nodes, err = p.parseTransforms()
	p.root = &ast.Root{
		Transforms: nodes,
	}

	return p.root, err
}

func (p *Parser) parseTransforms() ([]ast.Node, error) {
	var nodes []ast.Node
	for {
		var node, err = p.parseTransform(0)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
		if p.scanner.Eof() {
			return nodes, nil
		}
	}
}

func (p *Parser) parseTransform(indent int) (ast.Node, error) {
	var t, err = p.scanner.NextToken()
	if err != nil {
		return nil, err
	}
	if t.Typ != token.IDENTIFIER {
		return nil, fmt.Errorf("expected identifier, got %s", t)
	}

	switch t.Lit {
	case "from":
		return p.parseFromTransform()
	case "select":
		return p.parseSelectTransform()
	default:
		return nil, fmt.Errorf("unexpected identifier %s", t)
	}
}

func expectType(t scanner.Token, typ token.Token) error {
	if t.Typ != typ {
		return fmt.Errorf("expected %s, got %s", typ, t)
	}
	return nil
}

func expectExact(t scanner.Token, typ token.Token, lit string) error {
	if t.Typ != typ {
		return fmt.Errorf("expected %s, got %s", typ, t)
	}
	if lit == "*" {
		// any lit
		return nil
	}
	if t.Lit != lit {
		return fmt.Errorf("expected %s with lit %q, got %s", typ, lit, t)
	}
	return nil
}

func (p *Parser) parseFromTransform() (ast.Node, error) {
	// from table1
	// from e1 = table1
	var t1 = p.scanner.CurrToken()
	if err := expectExact(t1, token.IDENTIFIER, "from"); err != nil {
		return nil, err
	}

	var t2, err2 = p.scanner.NextToken()
	if err2 != nil {
		return nil, err2
	}
	if err := expectType(t1, token.IDENTIFIER); err != nil {
		return nil, err
	}

	var t3, err3 = p.scanner.NextToken()
	if err3 != nil {
		return nil, err3
	}
	if t3.Typ == token.NEWLINE || t3.Typ == token.EOF {
		return ast.FromTransform{
			Table: ast.Ident{Name: t2.Lit, Pos: t2.Pos},
		}, nil
	}

	if err := expectExact(t3, token.ASSIGN, "="); err != nil {
		return nil, err
	}

	var t4, err4 = p.scanner.NextToken()
	if err4 != nil {
		return nil, err2
	}
	if err := expectType(t4, token.IDENTIFIER); err != nil {
		return nil, err
	}

	// move to next one
	var _, err5 = p.scanner.NextToken()
	if err5 != nil {
		return nil, err5
	}

	return ast.FromTransform{
		Alias: &ast.Ident{Name: t2.Lit, Pos: t2.Pos},
		Table: ast.Ident{Name: t4.Lit, Pos: t4.Pos},
	}, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	var t = p.scanner.CurrToken()
	if err := expectType(t, token.IDENTIFIER); err != nil {
		return nil, err
	}
	return ast.Column{
		Name: ast.Ident{Name: t.Lit, Pos: t.Pos},
	}, nil
}

func (p *Parser) parseList() (ast.List, error) {
	var list = ast.List{
		Items: nil,
	}

	var t1 = p.scanner.CurrToken()
	if err := expectExact(t1, token.LBRACK, "["); err != nil {
		return list, err
	}
	for {
		var t2 = p.scanner.CurrToken()
		if t2.Typ == token.RBRACK && t2.Lit == "]" {
			break
		}

		var expr, err = p.parseExpr()
		if err != nil {
			return list, err
		}
		list.Items = append(list.Items, expr)
	}
	if err := expectExact(t1, token.RBRACK, "]"); err != nil {
		return list, err
	}
	return list, nil
}

func (p *Parser) parseSelectTransform() (ast.Node, error) {
	// select column1
	// select [column1, column2]
	// select [column1 = f"{column1} {column2}", column2]
	// select column1 = f"{column1} {column2}"
	var t1 = p.scanner.CurrToken()
	if err := expectExact(t1, token.IDENTIFIER, "select"); err != nil {
		return nil, err
	}

	var t2, err2 = p.scanner.NextToken()
	if err2 != nil {
		return nil, err2
	}

	var list ast.List
	if t2.Typ == token.LBRACK {
		var err error
		list, err = p.parseList()
		if err != nil {
			return nil, err
		}
	} else {
		var expr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, expr)
	}

	// move to next one
	var _, err3 = p.scanner.NextToken()
	if err3 != nil {
		return nil, err3
	}

	return ast.SelectTransform{List: list}, nil
}
