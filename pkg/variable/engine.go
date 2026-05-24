package variable

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Engine resolves variable expressions in template strings.
type Engine struct {
	captured  map[string]any
	env       map[string]string
	sequences map[string]*int64
	mu        sync.Mutex
}

// NewEngine creates a new variable resolution engine.
func NewEngine(captured map[string]any, env map[string]string) *Engine {
	if captured == nil {
		captured = make(map[string]any)
	}
	if env == nil {
		env = make(map[string]string)
		// Load from OS environment
		for _, e := range os.Environ() {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
	}
	return &Engine{
		captured:  captured,
		env:       env,
		sequences: make(map[string]*int64),
	}
}

// Resolve resolves all {{...}} expressions in a template string.
func (e *Engine) Resolve(template string) (string, error) {
	t, err := ParseTemplate(template)
	if err != nil {
		return template, err
	}

	var result strings.Builder
	for _, part := range t.Parts {
		if part.Expr != nil {
			val, err := e.resolveExpression(part.Expr)
			if err != nil {
				return template, err
			}
			result.WriteString(val)
		} else {
			result.WriteString(part.Literal)
		}
	}

	return result.String(), nil
}

// ResolveMap resolves all {{...}} expressions in string values of a map recursively.
func (e *Engine) ResolveMap(m map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(m))
	for k, v := range m {
		resolved, err := e.resolveValue(v)
		if err != nil {
			return nil, fmt.Errorf("resolve key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

// Capture stores a value captured from a test result assertion.
func (e *Engine) Capture(name string, value any) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.captured[name] = value
}

// GetCaptured returns all captured values.
func (e *Engine) GetCaptured() map[string]any {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make(map[string]any, len(e.captured))
	for k, v := range e.captured {
		result[k] = v
	}
	return result
}

func (e *Engine) resolveValue(v any) (any, error) {
	switch val := v.(type) {
	case string:
		return e.Resolve(val)
	case map[string]any:
		return e.ResolveMap(val)
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			resolved, err := e.resolveValue(item)
			if err != nil {
				return nil, err
			}
			result[i] = resolved
		}
		return result, nil
	default:
		return v, nil
	}
}

func (e *Engine) resolveExpression(expr *Expression) (string, error) {
	if expr.IsTernary {
		return e.resolveTernary(expr.Ternary)
	}

	if expr.IsFunc {
		return e.resolveFunc(expr.FuncName, expr.FuncArgs)
	}

	// Simple variable lookup
	return e.resolveVar(expr.VarName), nil
}

func (e *Engine) resolveVar(name string) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Check captured variables
	if val, ok := e.captured[name]; ok {
		return fmt.Sprintf("%v", val)
	}

	// Check environment variables
	if val, ok := e.env[name]; ok {
		return val
	}

	// Return as-is if not found (might be a literal)
	return "{{" + name + "}}"
}

func (e *Engine) resolveFunc(name string, args []string) (string, error) {
	switch name {
	case "$uuid":
		return randomUUID(), nil
	case "$random_email":
		return randomEmail(), nil
	case "$random_string":
		length := ParseIntArg(firstArg(args), 8)
		return randomString(length), nil
	case "$random_int":
		min := ParseIntArg(argAt(args, 0), 0)
		max := ParseIntArg(argAt(args, 1), 100)
		return fmt.Sprintf("%d", randomInt(min, max)), nil
	case "$random_float":
		min := ParseFloatArg(argAt(args, 0), 0)
		max := ParseFloatArg(argAt(args, 1), 100)
		decimals := ParseIntArg(argAt(args, 2), 2)
		return fmt.Sprintf("%.*f", decimals, randomFloat(min, max, decimals)), nil
	case "$random_bool":
		return fmt.Sprintf("%t", randomBool()), nil
	case "$random_name":
		return randomName(), nil
	case "$random_first_name":
		return randomFirstName(), nil
	case "$random_last_name":
		return randomLastName(), nil
	case "$random_phone":
		return randomPhone(), nil
	case "$random_ipv4":
		return randomIPv4(), nil
	case "$random_ipv6":
		return randomIPv6(), nil
	case "$random_enum":
		return randomEnum(args), nil
	case "$random_hex":
		length := ParseIntArg(firstArg(args), 16)
		return randomHex(length), nil
	case "$random_alpha":
		length := ParseIntArg(firstArg(args), 8)
		return randomAlpha(length), nil
	case "$random_digits":
		length := ParseIntArg(firstArg(args), 6)
		return randomDigits(length), nil
	case "$random_sentence":
		wordCount := ParseIntArg(firstArg(args), 6)
		return randomSentence(wordCount), nil
	case "$random_url":
		return randomURL(), nil
	case "$timestamp":
		return timestampSeconds(), nil
	case "$timestamp_ms":
		return timestampMS(), nil
	case "$iso_date":
		return isoDate(), nil
	case "$date":
		format := firstArg(args)
		return dateWithFormat(format), nil
	case "$env":
		if len(args) > 0 {
			return e.env[args[0]], nil
		}
		return "", nil
	case "$seq":
		name := firstArg(args)
		if name == "" {
			name = "default"
		}
		return e.seq(name), nil
	default:
		// Handle $env.XXX pattern (parsed as single func name without parens)
		if strings.HasPrefix(name, "$env.") {
			key := strings.TrimPrefix(name, "$env.")
			return e.env[key], nil
		}
		// Unknown function, return as-is
		return "{{" + name + "}}", nil
	}
}

func (e *Engine) resolveTernary(ternary *TernaryExpr) (string, error) {
	// Evaluate condition
	condResult, err := e.evaluateCondition(ternary.Condition)
	if err != nil {
		return "", fmt.Errorf("evaluate condition: %w", err)
	}

	if condResult {
		return e.Resolve(ternary.TrueVal)
	}
	return e.Resolve(ternary.FalseVal)
}

func (e *Engine) evaluateCondition(cond string) (bool, error) {
	cond = strings.TrimSpace(cond)

	// Check for function that returns bool
	if strings.HasPrefix(cond, "$") {
		val, err := e.resolveFunc(cond, nil)
		if err != nil {
			return false, err
		}
		return val == "true", nil
	}

	// Check for comparison operators
	for _, op := range []string{"!=", "==", ">=", "<=", ">", "<"} {
		if idx := strings.Index(cond, op); idx != -1 {
			left := strings.TrimSpace(cond[:idx])
			right := strings.TrimSpace(cond[idx+len(op):])

			leftVal := e.resolveVar(left)
			rightVal := e.resolveVar(right)

			switch op {
			case "==":
				return leftVal == rightVal, nil
			case "!=":
				return leftVal != rightVal, nil
			case ">":
				return compareNumeric(leftVal, rightVal) > 0, nil
			case "<":
				return compareNumeric(leftVal, rightVal) < 0, nil
			case ">=":
				return compareNumeric(leftVal, rightVal) >= 0, nil
			case "<=":
				return compareNumeric(leftVal, rightVal) <= 0, nil
			}
		}
	}

	// Simple truthiness check
	val := e.resolveVar(cond)
	return val != "" && val != "false" && val != "0" && val != "{{"+cond+"}}", nil
}

func compareNumeric(a, b string) int {
	var fa, fb float64
	fmt.Sscanf(a, "%f", &fa)
	fmt.Sscanf(b, "%f", &fb)
	if fa < fb {
		return -1
	}
	if fa > fb {
		return 1
	}
	return 0
}

func (e *Engine) seq(name string) string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.sequences[name] == nil {
		e.sequences[name] = new(int64)
	}
	*e.sequences[name]++
	return fmt.Sprintf("%d", *e.sequences[name])
}

func firstArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

func argAt(args []string, index int) string {
	if index < len(args) {
		return args[index]
	}
	return ""
}

// ResolveValue resolves a single value (used by the HTTP executor for body, headers, etc.)
func (e *Engine) ResolveValue(v any) (any, error) {
	return e.resolveValue(v)
}

// ResolveJSON resolves variables in a JSON byte slice.
func (e *Engine) ResolveJSON(data []byte) ([]byte, error) {
	var m any
	if err := json.Unmarshal(data, &m); err != nil {
		// Not JSON, try as string template
		resolved, err := e.Resolve(string(data))
		if err != nil {
			return data, err
		}
		return []byte(resolved), nil
	}

	resolved, err := e.resolveValue(m)
	if err != nil {
		return data, err
	}

	return json.Marshal(resolved)
}
