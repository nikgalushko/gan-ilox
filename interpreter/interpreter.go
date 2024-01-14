package interpreter

import (
	"errors"
	"fmt"

	"github.com/nikgalushko/gan-ilox/env"
	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/token/kind"
)

type Interpreter struct {
	env   *env.Environment
	stmts []internal.Stmt
	err   error
}

func New(env *env.Environment, stmts []internal.Stmt) *Interpreter {
	return &Interpreter{env: env, stmts: stmts}
}

func (i *Interpreter) Interpret() ([]any, error) {
	var ret []any
	for _, s := range i.stmts {
		v, err := i.Exec(s)
		if err != nil {
			return nil, err
		}
		if v != nil && !v.(internal.Literal).IsNil() {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (i *Interpreter) eval(e internal.Expr) (any, error) {
	ret := e.Accept(i)
	return ret, i.err
}

func (i *Interpreter) Exec(s internal.Stmt) (any, error) {
	ret := s.Accept(i)
	return ret, i.err
}

func (i *Interpreter) VisitForSmt(s internal.ForStmt) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	if s.Initializer != nil {
		prevEnv := i.env
		forEnv := env.NewWithParent(prevEnv)
		i.env = forEnv
		defer func() {
			i.env = prevEnv
		}()

		_, err := i.Exec(s.Initializer)
		if err != nil {
			return internal.LiteralNil
		}
	}

	evalCond := func() bool {
		cond, err := i.eval(s.Condition)
		if err != nil {
			i.err = err
			return false
		}

		return cond.(internal.Literal).AsBool()
	}

	for evalCond() {
		_, err := i.Exec(s.Body)
		if err != nil {
			i.err = err
			break
		}

		if s.Step != nil {
			_, err = i.eval(s.Step)
			if err != nil {
				i.err = err
				break
			}
		}
	}

	return internal.LiteralNil
}

func (i *Interpreter) VisitIfStmt(s internal.IfStmt) any {
	conditionResult, err := i.eval(s.Condition)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	var ret any
	if conditionResult.(internal.Literal).AsBool() {
		ret, _ = i.Exec(s.If)
	} else if s.Else != nil {
		ret, _ = i.Exec(s.Else)
	}

	return ret
}

func (i *Interpreter) VisitElseStmt(s internal.ElseStmt) any {
	var ret any
	if s.If != nil {
		ret, _ = i.Exec(s.If)
	} else {
		ret, _ = i.Exec(s.Block)
	}

	return ret
}

func (i *Interpreter) VisitVarStmt(s internal.VarStmt) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	value := internal.LiteralNil
	name := s.Name

	if s.Expression != nil {
		v, err := i.eval(s.Expression)
		if err == nil {
			value = v.(internal.Literal)
			i.err = err
		}
	}

	i.env.Define(name, value)

	return internal.LiteralNil
}

func (i *Interpreter) VisitPrintStmt(s internal.PrintStmt) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	fmt.Println(val.(internal.Literal).String())

	return internal.LiteralNil
}

func (i *Interpreter) VisitStmtExpression(s internal.StmtExpression) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	ret, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		ret = internal.LiteralNil
	}

	return ret
}

func (i *Interpreter) VisitBlockStmt(s internal.BlockStmt) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	prevEnv := i.env
	blockEnv := env.NewWithParent(prevEnv)
	i.env = blockEnv
	defer func() {
		i.env = prevEnv
	}()

	for _, s := range s.Stmts {
		_, err := i.Exec(s)
		if err != nil {
			i.err = err
			return nil
		}
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(e internal.Call) any {
	callee, err := i.eval(e.Callee)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	var args []internal.Literal
	for _, e := range e.Arguments {
		a, err := i.eval(e)
		if err != nil {
			i.err = err
			return internal.LiteralNil
		}
		args = append(args, a.(internal.Literal))
	}

	var ret any
	if callee.(internal.Literal).IsFunction() {
		ret, err = callee.(internal.Literal).AsFunction().Call(args, i)
	} else {
		err = errors.New("this type is not callable")
	}

	if err != nil {
		i.err = err
		ret = internal.LiteralNil
	}

	return ret
}

func (i *Interpreter) VisitLogicalExpr(e internal.Logical) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.eval(e.Left)
	if err != nil {
		return internal.LiteralNil
	}
	leftResult := val.(internal.Literal)
	needToComputeRightExpression := false
	switch e.Operator {
	case kind.Or:
		needToComputeRightExpression = !leftResult.AsBool()
	case kind.And:
		needToComputeRightExpression = leftResult.AsBool()
	}

	if needToComputeRightExpression {
		val, err = i.eval(e.Right)
		if err != nil {
			i.err = err
			val = internal.LiteralNil
		}
	}

	return val
}

func (i *Interpreter) VisitAssignmentExpr(e internal.Assignment) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.eval(e.Expression)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	i.env.Assign(e.Name, val.(internal.Literal))

	return val
}

func (i *Interpreter) VisitVariableExpr(e internal.Variable) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.env.Get(e.Name)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	return val
}

func (i *Interpreter) VisitBinaryExpr(expression internal.Binary) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	left := expression.Left.Accept(i).(internal.Literal)
	right := expression.Right.Accept(i).(internal.Literal)

	var ret internal.Literal
	switch expression.Operator {
	case kind.Minus:
		ret, i.err = sub(left, right)
	case kind.Plus:
		ret, i.err = add(left, right)
	case kind.Slash:
		ret, i.err = div(left, right)
	case kind.Star:
		ret, i.err = mul(left, right)
	case kind.Less:
		ret, i.err = less(left, right)
	case kind.LessEqual:
		ret, i.err = lessOrEqual(left, right)
	case kind.Greater:
		ret, i.err = graeater(left, right)
	case kind.GreaterEqual:
		ret, i.err = graeaterOrEqual(left, right)
	case kind.EqualEqual:
		ret, i.err = equal(left, right)
	}

	return ret
}

func (i *Interpreter) VisitGroupingExpr(expression internal.Grouping) any {
	if i.err != nil {
		return internal.LiteralNil
	}
	return expression.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expression internal.LiteralExpr) any {
	if i.err != nil {
		return internal.LiteralNil
	}
	return expression.Value
}

func (i *Interpreter) VisitUnaryExpr(expression internal.Unary) any {
	if i.err != nil {
		return internal.LiteralNil
	}

	val := expression.Right.Accept(i).(internal.Literal)

	switch expression.Operator {
	case kind.Bang:
		return internal.NewLiteralBool(!val.AsBool())
	case kind.Minus:
		if val.IsInt() {
			return internal.NewLiteralInt(-val.AsInt())
		} else if val.IsFloat() {
			return internal.NewLiteralFloat(-val.AsFloat())
		}

		i.err = errors.New("Illegal operation") // TODO: craete more freandly error message
		return internal.LiteralNil
	case kind.BitwiseNot:
		if val.IsInt() {
			return internal.NewLiteralInt(^val.AsInt())
		}
		i.err = errors.New("bitwise operator can be used only with integer number")
		return internal.LiteralNil
	}

	panic("unreachable code")
}
