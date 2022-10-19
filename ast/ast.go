package ast

import "github.com/siadat/prql-parser/scanner"

type Ident struct {
	Name string
	Pos  scanner.Pos
}

type FromTransform struct {
	Alias *Ident
	Table Ident
}

type SelectTransform struct {
	List List
}

type List struct {
	Items []Expr
}

type Root struct {
	Transforms []Node
}

type Column struct {
	Name Ident
}

type Node interface {
	node()
}
type Expr interface {
	expr()
}

func (FromTransform) node()   {}
func (SelectTransform) node() {}
func (List) node()            {}
func (Root) node()            {}
func (Column) node()          {}

func (Column) expr() {}
