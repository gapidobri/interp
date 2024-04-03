package interpreter

import (
	"interp/errors"
	"interp/expr"
	. "interp/token"
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
	case Greater:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) > right.(float64), nil
	case GreaterEqual:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) >= right.(float64), nil
	case Less:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) < right.(float64), nil
	case LessEqual:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) <= right.(float64), nil
	case BangEqual:
		return left != right, nil
	case EqualEqual:
		return left == right, nil
	case Minus:
		err := i.checkNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) - right.(float64), nil
	case Plus:
		if i.isFloat(left) && i.isFloat(right) {
			return left.(float64) + right.(float64), nil
		}
		if i.isString(left) && i.isString(right) {
			return left.(string) + right.(string), nil
		}

		return nil, errors.NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	case Slash:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		if right.(float64) == 0 {
			return nil, errors.NewRuntimeError(expr.Operator, "Can not divide by zero.")
		}

		return left.(float64) / right.(float64), nil
	case Star:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) * right.(float64), nil
	}

	return nil, nil
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

	if expr.Operator.Type == Or {
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
	case Bang:
		return !i.isTruthy(right), nil
	case Minus:
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
