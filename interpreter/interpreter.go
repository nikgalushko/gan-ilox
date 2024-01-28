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
	stmts []internal.Stmt[internal.Literal]
	err   error
}

func New(env *env.Environment, stmts []internal.Stmt[internal.Literal]) *Interpreter {
	return &Interpreter{env: env, stmts: stmts}
}

func (i *Interpreter) Interpret() ([]any, error) {
	var ret []any
	for _, s := range i.stmts {
		v, err := i.Exec(s)
		if err != nil {
			return nil, err
		}
		if !v.IsNil() {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (i *Interpreter) eval(e internal.Expr[internal.Literal]) (internal.Literal, error) {
	ret := e.Accept(i)
	return ret, i.err
}

func (i *Interpreter) Exec(s internal.Stmt[internal.Literal]) (internal.Literal, error) {
	ret := s.Accept(i)
	return ret, i.err
}

func (i *Interpreter) VisitSetExpr(e internal.SetExpr[internal.Literal]) internal.Literal {
	obj, err := i.eval(e.Object)
	if err != nil {
		return internal.LiteralNil
	}

	if !obj.IsClassInstance() {
		i.err = errors.New("only instances have fields")
		return internal.LiteralNil
	}

	value, err := i.eval(e.Value)
	if err == nil {
		obj.AsClassInstance().Set(e.Name, value)
	}

	return internal.LiteralNil
}

func (i *Interpreter) VisitClassStmt(c internal.ClassStmt[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	methods := make(map[string]internal.Literal)
	for _, s := range c.Methods {
		methods[s.Name] = internal.NewLiteralUserFunction(s.Parameters, s.Body)
	}
	i.env.Define(c.Name, internal.NewLiteralClass(c.Name, methods))

	return internal.LiteralNil
}

func (i *Interpreter) VisitReturnStmt(s internal.RreturnStmt[internal.Literal]) internal.Literal {
	var ret internal.Literal
	if s.Expression != nil {
		ret, i.err = i.eval(s.Expression)
	} else {
		ret = internal.LiteralNil
	}
	ret = ret.AsReturnResult()

	return ret
}

func (i *Interpreter) VisitForSmt(s internal.ForStmt[internal.Literal]) internal.Literal {
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
		if s.Condition == nil {
			return true
		}

		cond, err := i.eval(s.Condition)
		if err != nil {
			i.err = err
			return false
		}

		return cond.AsBool()
	}

	for evalCond() {
		ret, err := i.Exec(s.Body)
		if err != nil {
			i.err = err
			break
		}

		if ret.IsReturnResult() {
			return ret
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

func (i *Interpreter) VisitIfStmt(s internal.IfStmt[internal.Literal]) internal.Literal {
	conditionResult, err := i.eval(s.Condition)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	var ret internal.Literal
	if conditionResult.AsBool() {
		ret, _ = i.Exec(s.If)
	} else if s.Else != nil {
		ret, _ = i.Exec(s.Else)
	}

	return ret
}

func (i *Interpreter) VisitElseStmt(s internal.ElseStmt[internal.Literal]) internal.Literal {
	var ret internal.Literal
	if s.If != nil {
		ret, _ = i.Exec(s.If)
	} else {
		ret, _ = i.Exec(s.Block)
	}

	return ret
}

func (i *Interpreter) VisitFuncStmt(s internal.FuncStmt[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	i.env.Define(s.Name, internal.NewLiteralUserFunction(
		s.Parameters,
		s.Body,
	))

	return internal.LiteralNil
}

func (i *Interpreter) VisitVarStmt(s internal.VarStmt[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	value := internal.LiteralNil
	name := s.Name

	if s.Expression != nil {
		v, err := i.eval(s.Expression)
		if err == nil {
			value = v
			i.err = err
		}
	}

	i.env.Define(name, value)

	return internal.LiteralNil
}

func (i *Interpreter) VisitPrintStmt(s internal.PrintStmt[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	fmt.Println(val.String())

	return internal.LiteralNil
}

func (i *Interpreter) VisitStmtExpression(s internal.StmtExpression[internal.Literal]) internal.Literal {
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

func (i *Interpreter) VisitBlockStmt(s internal.BlockStmt[internal.Literal]) internal.Literal {
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
		ret, err := i.Exec(s)
		if err != nil {
			i.err = err
			return internal.LiteralNil
		}

		if ret.IsReturnResult() {
			return ret
		}
	}

	return internal.LiteralNil
}

func (i *Interpreter) VisitCallExpr(e internal.Call[internal.Literal]) internal.Literal {
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
		args = append(args, a)
	}

	var ret internal.Literal
	if callee.IsFunction() {
		f := callee.AsFunction()
		prevEnv := i.env
		funEnv := env.NewWithParent(prevEnv)
		i.env = funEnv
		defer func() {
			i.env = prevEnv
		}()

		for idx := range args {
			i.env.Define(f.ArgumentsName[idx], args[idx])
		}

		ret, err = f.Call(args, i)
	} else if callee.IsClass() {
		// TODO: refactor it
		c := callee.AsClass()
		prevEnv := i.env
		funEnv := env.NewWithParent(prevEnv)
		i.env = funEnv
		defer func() {
			i.env = prevEnv
		}()

		if c.Initializer != nil {
			initializer := c.Initializer.AsFunction()
			for idx := range args {
				i.env.Define(initializer.ArgumentsName[idx], args[idx])
			}

			ret, err = c.Call(args, i)
		} else {
			ret, err = c.Call(nil, i)
		}
	} else {
		err = errors.New("this type is not callable")
	}

	if err != nil {
		i.err = err
		ret = internal.LiteralNil
	}

	return ret
}

func (i *Interpreter) VisitLogicalExpr(e internal.Logical[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.eval(e.Left)
	if err != nil {
		return internal.LiteralNil
	}
	leftResult := val
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

func (i *Interpreter) VisitAssignmentExpr(e internal.Assignment[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	val, err := i.eval(e.Expression)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	i.env.Assign(e.Name, val)

	return val
}

func (i *Interpreter) VisitVariableExpr(e internal.Variable[internal.Literal]) internal.Literal {
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

func (i *Interpreter) VisitBinaryExpr(expression internal.Binary[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	left := expression.Left.Accept(i)
	right := expression.Right.Accept(i)

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

func (i *Interpreter) VisitGroupingExpr(expression internal.Grouping[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}
	return expression.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expression internal.LiteralExpr[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}
	return expression.Value
}

func (i *Interpreter) VisitUnaryExpr(expression internal.Unary[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	val := expression.Right.Accept(i)

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

func (i *Interpreter) VisitGetExpr(e internal.GetExpr[internal.Literal]) internal.Literal {
	if i.err != nil {
		return internal.LiteralNil
	}

	v, err := i.eval(e.Expression)
	if err != nil {
		return internal.LiteralNil
	}

	if !v.IsClassInstance() {
		i.err = fmt.Errorf("only instances have property: %w", err)
		return internal.LiteralNil
	}

	ret, err := v.AsClassInstance().Get(e.Name)
	if err != nil {
		i.err = err
		return internal.LiteralNil
	}

	return ret
}
