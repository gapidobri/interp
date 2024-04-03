package interpreter

import (
	"fmt"
	"interp/errors"
	"interp/expr"
	"interp/token"
)

func (i *Interpreter) VisitBinaryExpr(expr *expr.Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Greater:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) > right.(float64), nil
	case token.GreaterEqual:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) >= right.(float64), nil
	case token.Less:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) < right.(float64), nil
	case token.LessEqual:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) <= right.(float64), nil
	case token.BangEqual:
		return left != right, nil
	case token.EqualEqual:
		return left == right, nil
	case token.Minus:
		err := i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) - right.(float64), nil
	case token.Plus:
		if i.isFloat(left) && i.isFloat(right) {
			return left.(float64) + right.(float64), nil
		}
		if i.isString(left) && i.isString(right) {
			return left.(string) + right.(string), nil
		}

		return nil, errors.NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	case token.Slash:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		if right.(float64) == 0 {
			return nil, errors.NewRuntimeError(expr.Operator, "Can not divide by zero.")
		}

		return left.(float64) / right.(float64), nil
	case token.Star:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) * right.(float64), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr *expr.Call) (any, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	var arguments []any
	for _, argument := range expr.Arguments {
		value, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, value)
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, errors.NewRuntimeError(expr.Paren, "Can only call functions and classes.")
	}

	if len(arguments) != function.arity() {
		return nil, errors.NewRuntimeError(
			expr.Paren,
			fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments)),
		)
	}

	return function.call(i, arguments)
}

func (i *Interpreter) VisitGroupingExpr(expr *expr.Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *expr.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *expr.Logical) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.Or {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitUnaryExpr(expr *expr.Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.Bang:
		return !i.isTruthy(right), nil
	case token.Minus:
		return -right.(float64), nil
	}

	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr *expr.Variable) (any, error) {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *expr.Assign) (any, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	err = i.environment.Assign(expr.Name, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
