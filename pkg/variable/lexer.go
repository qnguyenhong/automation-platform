package variable

import (
	"fmt"
	"strconv"
	"strings"
)

// TokenType represents the type of a token.
type TokenType int

const (
	TokenLiteral TokenType = iota
	TokenName           // variable name like "user_id"
	TokenFunc           // function call like "$random_int(1,10)"
	TokenTernary        // ternary expression
)

// Token represents a parsed token from a template string.
type Token struct {
	Type    TokenType
	Value   string
	Args    []string // for function calls
	Ternary *TernaryExpr // for ternary expressions
}

// TernaryExpr represents a condition ? trueVal : falseVal expression.
type TernaryExpr struct {
	Condition string
	TrueVal   string
	FalseVal  string
}

// Template represents a parsed template with literal parts and variable expressions.
type Template struct {
	Parts []TemplatePart
}

// TemplatePart is one piece of a template: literal text or a variable expression.
type TemplatePart struct {
	Literal string
	Expr    *Expression
}

// Expression is a variable or function expression to resolve.
type Expression struct {
	Raw        string
	IsFunc     bool
	FuncName   string
	FuncArgs   []string
	VarName    string
	IsTernary  bool
	Ternary    *TernaryExpr
}

// ParseTemplate parses a template string containing {{...}} expressions.
// Examples:
//
//	"Hello {{name}}" → literal "Hello " + var "name"
//	"{{$random_int(1,10)}}" → func "$random_int" args ["1","10"]
//	"{{$random_int(1,10) > 5 ? high : low}}" → ternary
func ParseTemplate(s string) (*Template, error) {
	t := &Template{}
	i := 0

	for i < len(s) {
		// Find next {{
		start := strings.Index(s[i:], "{{")
		if start == -1 {
			// Rest is literal
			t.Parts = append(t.Parts, TemplatePart{Literal: s[i:]})
			break
		}

		// Add literal before {{
		if start > 0 {
			t.Parts = append(t.Parts, TemplatePart{Literal: s[i : i+start]})
		}

		// Find matching }}
		exprStart := i + start + 2
		end := findClosingBrace(s, exprStart)
		if end == -1 {
			return nil, fmt.Errorf("unclosed {{ at position %d", i+start)
		}

		exprStr := strings.TrimSpace(s[exprStart:end])
		expr, err := parseExpression(exprStr)
		if err != nil {
			return nil, fmt.Errorf("invalid expression %q: %w", exprStr, err)
		}
		t.Parts = append(t.Parts, TemplatePart{Expr: expr})

		i = end + 2
	}

	return t, nil
}

// findClosingBrace finds the position of }} matching the opening {{.
func findClosingBrace(s string, start int) int {
	depth := 0
	for i := start; i < len(s)-1; i++ {
		if s[i] == '{' && s[i+1] == '{' {
			depth++
			i++
		} else if s[i] == '}' && s[i+1] == '}' {
			if depth == 0 {
				return i
			}
			depth--
			i++
		}
	}
	return -1
}

// parseExpression parses a single expression inside {{...}}.
func parseExpression(s string) (*Expression, error) {
	s = strings.TrimSpace(s)

	// Check for ternary: expr ? trueVal : falseVal
	if idx := findTernaryQuestion(s); idx != -1 {
		cond := strings.TrimSpace(s[:idx])
		rest := s[idx+1:]
		colonIdx := findTernaryColon(rest)
		if colonIdx == -1 {
			return nil, fmt.Errorf("invalid ternary expression: missing colon")
		}
		trueVal := strings.TrimSpace(rest[:colonIdx])
		falseVal := strings.TrimSpace(rest[colonIdx+1:])

		return &Expression{
			Raw:       s,
			IsTernary: true,
			Ternary: &TernaryExpr{
				Condition: cond,
				TrueVal:   trueVal,
				FalseVal:  falseVal,
			},
		}, nil
	}

	// Check for function call: $func(args)
	if strings.HasPrefix(s, "$") {
		parenIdx := strings.Index(s, "(")
		if parenIdx != -1 && s[len(s)-1] == ')' {
			funcName := s[:parenIdx]
			argsStr := s[parenIdx+1 : len(s)-1]
			args := parseFuncArgs(argsStr)
			return &Expression{
				Raw:      s,
				IsFunc:   true,
				FuncName: funcName,
				FuncArgs: args,
			}, nil
		}
		// $func without parens (like $uuid)
		return &Expression{
			Raw:      s,
			IsFunc:   true,
			FuncName: s,
		}, nil
	}

	// Simple variable name
	return &Expression{
		Raw:     s,
		VarName: s,
	}, nil
}

// parseFuncArgs splits function arguments by comma, respecting nested parentheses.
func parseFuncArgs(s string) []string {
	if s == "" {
		return nil
	}

	var args []string
	depth := 0
	start := 0

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				args = append(args, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	args = append(args, strings.TrimSpace(s[start:]))
	return args
}

// findTernaryQuestion finds the ? that's not inside parentheses.
func findTernaryQuestion(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth++
		case '?':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// findTernaryColon finds the : that's not inside parentheses.
func findTernaryColon(s string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ':':
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// ParseIntArg parses a string argument as an integer.
func ParseIntArg(s string, defaultVal int) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

// ParseFloatArg parses a string argument as a float64.
func ParseFloatArg(s string, defaultVal float64) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultVal
	}
	return v
}
