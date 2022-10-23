package scanner_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/siadat/prql-parser/scanner"
	"github.com/siadat/prql-parser/token"
)

const IgnorePos = -1

func TestIdentifier(tt *testing.T) {
	var testCases = []struct {
		src  string
		want []scanner.Token
	}{
		{
			src: "abc  `table.column` $abc _abc",
			want: []scanner.Token{
				{token.IDENTIFIER, `abc`, 0},
				{token.WHITESPACE, `  `, 3},
				{token.IDENTIFIER, "`table.column`", 5},
				{token.WHITESPACE, ` `, 19},
				{token.IDENTIFIER, `$abc`, 20},
				{token.WHITESPACE, ` `, 24},
				{token.IDENTIFIER, `_abc`, 25},
			},
		},
		{
			src: `
from employees
select [id, first_name, age]
sort age
take 10
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", 0},
				{token.IDENTIFIER, `from`, 1},
				{token.WHITESPACE, ` `, 5},
				{token.IDENTIFIER, `employees`, 6},
				{token.NEWLINE, "\n", 15},

				{token.IDENTIFIER, `select`, 16},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `id`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `first_name`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `age`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `sort`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `age`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `take`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTEGER, `10`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from order               # This is a comment
filter status == "done"
sort [-amount]           # sort order
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},
				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `order`, IgnorePos},
				{token.WHITESPACE, `               `, IgnorePos},
				{token.COMMENT, `# This is a comment`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `filter`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `status`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.EQL, `==`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.STRING, `"done"`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `sort`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.SUB, `-`, IgnorePos},
				{token.IDENTIFIER, `amount`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.WHITESPACE, `           `, IgnorePos},
				{token.COMMENT, `# sort order`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from employees
derive [
  age_at_year_end = (@2022-12-31T00:00:00 - dob),
  first_check_in = start + +20 + 10days,
]
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},
				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `derive`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `age_at_year_end`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LPAREN, `(`, IgnorePos},
				{token.TIMESTAMP, `@2022-12-31T00:00:00`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.SUB, `-`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `dob`, IgnorePos},
				{token.RPAREN, `)`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `first_check_in`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `start`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ADD, `+`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ADD, `+`, IgnorePos},
				{token.INTEGER, `20`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ADD, `+`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTERVAL, `10days`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from employees
# Filter before aggregations
filter start_date > @2021-01-01
group country (
  aggregate [max_salary = max salary]
)
# And filter after aggregations!
filter max_salary > 100000
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},
				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.COMMENT, `# Filter before aggregations`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `filter`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `start_date`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.GTR, `>`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.DATE, `@2021-01-01`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `group`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `country`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LPAREN, `(`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `aggregate`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `max_salary`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `max`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `salary`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.RPAREN, `)`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.COMMENT, `# And filter after aggregations!`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `filter`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `max_salary`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.GTR, `>`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTEGER, `100000`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from web
# Just like Python
select url = f"http://www.{domain}.{tld}/{page}"
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},
				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `web`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.COMMENT, `# Just like Python`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `select`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `url`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.STRING, `f"http://www.{domain}.{tld}/{page}"`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from employees
group employee_id (
  sort month
  window rolling:12 (
    derive [trail_12_m_comp = sum paycheck]
  )
)
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},
				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `group`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employee_id`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LPAREN, `(`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `sort`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `month`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `window`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `rolling`, IgnorePos},
				{token.COLON, `:`, IgnorePos},
				{token.INTEGER, `12`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LPAREN, `(`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `    `, IgnorePos},
				{token.IDENTIFIER, `derive`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `trail_12_m_comp`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `sum`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `paycheck`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.RPAREN, `)`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.RPAREN, `)`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
func fahrenheit_from_celsius temp -> temp * 9/5 + 32

from weather
select temp_f = (fahrenheit_from_celsius temp_c)
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},
				{token.IDENTIFIER, `func`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `fahrenheit_from_celsius`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `temp`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ARROW, `->`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `temp`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.MUL, `*`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTEGER, `9`, IgnorePos},
				{token.QUO, `/`, IgnorePos},
				{token.INTEGER, `5`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ADD, `+`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTEGER, `32`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `weather`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `select`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `temp_f`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LPAREN, `(`, IgnorePos},
				{token.IDENTIFIER, `fahrenheit_from_celsius`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `temp_c`, IgnorePos},
				{token.RPAREN, `)`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
# Most recent employee in each role
# Quite difficult in SQL...
from employees
group role (
  sort join_date
  take 1
)
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},

				{token.COMMENT, `# Most recent employee in each role`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.COMMENT, `# Quite difficult in SQL...`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `group`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `role`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LPAREN, `(`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `sort`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `join_date`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.WHITESPACE, `  `, IgnorePos},
				{token.IDENTIFIER, `take`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTEGER, `1`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.RPAREN, ")", IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
derive db_version = s"version()"
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `derive`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `db_version`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.STRING, `s"version()"`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from employees
join benefits [employee_id]
join side:left p=positions [id==employee_id]
select [employee_id, role, vision_coverage]
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `join`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `benefits`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `employee_id`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `join`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `side`, IgnorePos},
				{token.COLON, `:`, IgnorePos},
				{token.IDENTIFIER, `left`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `p`, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.IDENTIFIER, `positions`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `id`, IgnorePos},
				{token.EQL, `==`, IgnorePos},
				{token.IDENTIFIER, `employee_id`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `select`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `employee_id`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `role`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `vision_coverage`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from users
filter last_login != null
filter deleted_at == null
derive channel = channel ?? "unknown"
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `users`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `filter`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `last_login`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.NEQ, `!=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `null`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `filter`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `deleted_at`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.EQL, `==`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `null`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `derive`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `channel`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `channel`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.COALESCE, `??`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.STRING, `"unknown"`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
prql dialect:mssql  # Will generate TOP rather than LIMIT

from employees
sort age
take 10
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `prql`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `dialect`, IgnorePos},
				{token.COLON, `:`, IgnorePos},
				{token.IDENTIFIER, `mssql`, IgnorePos},
				{token.WHITESPACE, `  `, IgnorePos},
				{token.COMMENT, `# Will generate TOP rather than LIMIT`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `sort`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `age`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `take`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.INTEGER, `10`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `
from e = employees
select e.first_name
`,
			want: []scanner.Token{
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `e`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.ASSIGN, `=`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},

				{token.IDENTIFIER, `select`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `e`, IgnorePos},
				{token.PERIOD, `.`, IgnorePos},
				{token.IDENTIFIER, `first_name`, IgnorePos},
				{token.NEWLINE, "\n", IgnorePos},
			},
		},
		{
			src: `from employees | filter department == "Product" | select [first_name, last_name]`,
			want: []scanner.Token{
				{token.IDENTIFIER, `from`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `employees`, IgnorePos},

				{token.WHITESPACE, ` `, IgnorePos},
				{token.PIPE, "|", IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},

				{token.IDENTIFIER, `filter`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `department`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.EQL, `==`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.STRING, `"Product"`, IgnorePos},

				{token.WHITESPACE, ` `, IgnorePos},
				{token.PIPE, "|", IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},

				{token.IDENTIFIER, `select`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.LBRACK, `[`, IgnorePos},
				{token.IDENTIFIER, `first_name`, IgnorePos},
				{token.COMMA, `,`, IgnorePos},
				{token.WHITESPACE, ` `, IgnorePos},
				{token.IDENTIFIER, `last_name`, IgnorePos},
				{token.RBRACK, `]`, IgnorePos},
			},
		},
	}

	for _, tc := range testCases {
		var src = tc.src
		var s = scanner.NewScanner(strings.NewReader(src))
		var got []scanner.Token
		var err error
		for {
			var t, gotErr = s.NextToken()
			if t.Typ == token.EOF {
				break
			}
			if gotErr != nil {
				err = gotErr
				break
			}
			got = append(got, t)
		}

		var cmpOpt = cmp.FilterValues(func(p1, p2 scanner.Pos) bool { return p1 == IgnorePos || p2 == IgnorePos || p1 == p2 }, cmp.Ignore())

		if diff := cmp.Diff((error)(nil), err); diff != "" {
			tt.Fatalf("case error failed to match src=%q (-want +got):\n%s", src, diff)
		}
		if diff := cmp.Diff(tc.want, got, cmpOpt); diff != "" {
			tt.Fatalf("case failed src=%q (-want +got):\n%s", src, diff)
		}
	}
}
