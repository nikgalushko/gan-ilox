package interpreter

import (
	"errors"
	"fmt"

	"github.com/nikgalushko/gan-ilox/env"
	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type Interpreter struct {
	env   *env.Environment
	stmts []expr.Stmt
	err   error
}

func New(env *env.Environment, stmts []expr.Stmt) *Interpreter {
	return &Interpreter{env: env, stmts: stmts}
}

func (i *Interpreter) Interpret() ([]any, error) {
	var ret []any
	for _, s := range i.stmts {
		v, err := i.exec(s)
		if err != nil {
			return nil, err
		}
		if v != nil && v.(token.Literal) != token.LiteralNil {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (i *Interpreter) eval(e expr.Expr) (any, error) {
	ret := e.Accept(i)
	return ret, i.err
}

func (i *Interpreter) exec(s expr.Stmt) (any, error) {
	ret := s.Accept(i)
	return ret, i.err
}

func (i *Interpreter) VisitForSmt(s expr.ForStmt) any {
	if i.err != nil {
		return token.LiteralNil
	}

	if s.Initializer != nil {
		prevEnv := i.env
		forEnv := env.NewWithParent(prevEnv)
		i.env = forEnv
		defer func() {
			i.env = prevEnv
		}()

		_, err := i.exec(s.Initializer)
		if err != nil {
			return token.LiteralNil
		}
	}

	evalCond := func() bool {
		cond, err := i.eval(s.Condition)
		if err != nil {
			i.err = err
			return false
		}

		return cond.(token.Literal).AsBool()
	}

	for evalCond() {
		_, err := i.exec(s.Body)
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

	return token.LiteralNil
}

func (i *Interpreter) VisitIfStmt(s expr.IfStmt) any {
	conditionResult, err := i.eval(s.Condition)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	var ret any
	if conditionResult.(token.Literal).AsBool() {
		ret, _ = i.exec(s.If)
	} else if s.Else != nil {
		ret, _ = i.exec(s.Else)
	}

	return ret
}

func (i *Interpreter) VisitElseStmt(s expr.ElseStmt) any {
	var ret any
	if s.If != nil {
		ret, _ = i.exec(s.If)
	} else {
		ret, _ = i.exec(s.Block)
	}

	return ret
}

func (i *Interpreter) VisitVarStmt(s expr.VarStmt) any {
	if i.err != nil {
		return token.LiteralNil
	}

	value := token.LiteralNil
	name := s.Name.Lexeme

	if s.Expression != nil {
		v, err := i.eval(s.Expression)
		if err == nil {
			value = v.(token.Literal)
			i.err = err
		}
	}

	i.env.Define(name, value)

	return token.LiteralNil
}

func (i *Interpreter) VisitPrintStmt(s expr.PrintStmt) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	fmt.Println(val.(token.Literal).String())

	return token.LiteralNil
}

func (i *Interpreter) VisitStmtExpression(s expr.StmtExpression) any {
	if i.err != nil {
		return token.LiteralNil
	}

	ret, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		ret = token.LiteralNil
	}

	return ret
}

func (i *Interpreter) VisitBlockStmt(s expr.BlockStmt) any {
	if i.err != nil {
		return token.LiteralNil
	}

	prevEnv := i.env
	blockEnv := env.NewWithParent(prevEnv)
	i.env = blockEnv
	defer func() {
		i.env = prevEnv
	}()

	for _, s := range s.Stmts {
		_, err := i.exec(s)
		if err != nil {
			i.err = err
			return nil
		}
	}

	return nil
}

func (i *Interpreter) VisitLogicalExpr(e expr.Logical) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val, err := i.eval(e.Left)
	if err != nil {
		return token.LiteralNil
	}
	leftResult := val.(token.Literal)
	needToComputeRightExpression := false
	switch e.Operator {
	case token.Or:
		needToComputeRightExpression = !leftResult.AsBool()
	case token.And:
		needToComputeRightExpression = leftResult.AsBool()
	}

	if needToComputeRightExpression {
		val, err = i.eval(e.Right)
		if err != nil {
			i.err = err
			val = token.LiteralNil
		}
	}

	return val
}

func (i *Interpreter) VisitAssignmentExpr(e expr.Assignment) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val, err := i.eval(e.Expression)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	i.env.Assign(e.Name.Lexeme, val.(token.Literal))

	return val
}

func (i *Interpreter) VisitVariableExpr(e expr.Variable) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val, err := i.env.Get(e.Name.Lexeme)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	return val
}

func (i *Interpreter) VisitBinaryExpr(expression expr.Binary) any {
	if i.err != nil {
		return token.LiteralNil
	}

	left := expression.Left.Accept(i).(token.Literal)
	right := expression.Right.Accept(i).(token.Literal)

	var ret token.Literal
	switch expression.Operator.Type {
	case token.Minus:
		ret, i.err = sub(left, right)
	case token.Plus:
		ret, i.err = add(left, right)
	case token.Slash:
		ret, i.err = div(left, right)
	case token.Star:
		ret, i.err = mul(left, right)
	case token.Less:
		ret, i.err = less(left, right)
	case token.LessEqual:
		ret, i.err = lessOrEqual(left, right)
	case token.Greater:
		ret, i.err = graeater(left, right)
	case token.GreaterEqual:
		ret, i.err = graeaterOrEqual(left, right)
	case token.EqualEqual:
		ret, i.err = equal(left, right)
	}

	return ret
}

func (i *Interpreter) VisitGroupingExpr(expression expr.Grouping) any {
	if i.err != nil {
		return token.LiteralNil
	}
	return expression.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expression expr.Literal) any {
	if i.err != nil {
		return token.LiteralNil
	}
	return expression.Value
}

func (i *Interpreter) VisitUnaryExpr(expression expr.Unary) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val := expression.Right.Accept(i).(token.Literal)

	switch expression.Operator.Type {
	case token.Bang:
		return token.NewLiteralBool(!val.AsBool())
	case token.Minus:
		if val.IsInt() {
			return token.NewLiteralInt(-val.AsInt())
		} else if val.IsFloat() {
			return token.NewLiteralFloat(-val.AsFloat())
		}

		i.err = errors.New("Illegal operation") // TODO: craete more freandly error message
		return token.LiteralNil
	case token.BitwiseNot:
		if val.IsInt() {
			return token.NewLiteralInt(^val.AsInt())
		}
		i.err = errors.New("bitwise operator can be used only with integer number")
		return token.LiteralNil
	}

	panic("unreachable code")
}
