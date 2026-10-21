package expr

// Expr is a node of an arithmetic expression tree.
type Expr interface {
	Eval(env map[string]float64) float64
	String() string
}

// Num is a constant.
type Num float64

// Var is a variable whose value comes from the environment.
type Var string

// Add is L + R.
type Add struct{ L, R Expr }

// Mul is L * R.
type Mul struct{ L, R Expr }

// Neg is -X.
type Neg struct{ X Expr }

// Eval returns the constant.
func (n Num) Eval(env map[string]float64) float64 {
	// TODO
	return 0
}

// String formats the constant with strconv's shortest 'g' form.
func (n Num) String() string {
	// TODO
	return ""
}

// Eval returns env[v], or 0 if it is missing.
func (v Var) Eval(env map[string]float64) float64 {
	// TODO
	return 0
}

// String returns the variable name.
func (v Var) String() string {
	// TODO
	return ""
}

// Eval returns L + R.
func (a Add) Eval(env map[string]float64) float64 {
	// TODO
	return 0
}

// String returns "(L + R)".
func (a Add) String() string {
	// TODO
	return ""
}

// Eval returns L * R.
func (m Mul) Eval(env map[string]float64) float64 {
	// TODO
	return 0
}

// String returns "(L * R)".
func (m Mul) String() string {
	// TODO
	return ""
}

// Eval returns -X.
func (n Neg) Eval(env map[string]float64) float64 {
	// TODO
	return 0
}

// String returns "-X".
func (n Neg) String() string {
	// TODO
	return ""
}

// Simplify returns an equivalent, simplified expression.
func Simplify(e Expr) Expr {
	// TODO
	return nil
}

// Vars returns the sorted, distinct variable names used in e.
func Vars(e Expr) []string {
	// TODO
	return nil
}

// Size returns the number of nodes in e.
func Size(e Expr) int {
	// TODO
	return 0
}
