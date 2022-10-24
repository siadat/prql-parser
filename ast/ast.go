package ast

import (
	"github.com/siadat/prql-parser/scanner"
	"github.com/siadat/prql-parser/token"
)

type Ident struct {
	Name string
	Pos  scanner.Pos
}

type FromTransform struct {
	Alias *Ident
	Table Ident
}

type SelectTransform struct {
	List ExprList
}

type DeriveTransform struct {
	List ExprList
}

type BinaryExpr struct {
	X  Expr
	Y  Expr
	Op token.Token
}

type UnaryExpr struct {
	X  Expr
	Op token.Token
}

type ParenExpr struct {
	X Expr
}

type AssignExpr struct {
	Name string
	Expr Expr
}

type ExprList struct {
	Items []Expr
}

type Root struct {
	Transforms []Node
}

type String struct {
	Value string
}

type Integer struct {
	Value int
}

type Date struct {
	Year, Month, Day int
}

type Time struct {
	Hour, Minute, Second int
}

type Timestamp struct {
	Year, Month, Day     int
	Hour, Minute, Second int
}
type Interval struct {
	Count int
	Unit  string
}

type Float struct {
	Value float64
}

type Column struct {
	Name Ident
}

type Node interface {
	node()
}

type Expr interface {
	node()
	expr()
}

func (FromTransform) node()   {}
func (SelectTransform) node() {}
func (DeriveTransform) node() {}
func (ExprList) node()        {}
func (Root) node()            {}
func (Column) node()          {}
func (Integer) node()         {}
func (Date) node()            {}
func (Time) node()            {}
func (Timestamp) node()       {}
func (Interval) node()        {}
func (String) node()          {}
func (Float) node()           {}
func (BinaryExpr) node()      {}
func (UnaryExpr) node()       {}
func (ParenExpr) node()       {}
func (AssignExpr) node()      {}

func (Column) expr()     {}
func (Integer) expr()    {}
func (Date) expr()       {}
func (Time) expr()       {}
func (Timestamp) expr()  {}
func (Interval) expr()   {}
func (String) expr()     {}
func (Float) expr()      {}
func (BinaryExpr) expr() {}
func (UnaryExpr) expr()  {}
func (ParenExpr) expr()  {}
func (AssignExpr) expr() {}
