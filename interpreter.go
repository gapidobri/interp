package main

import (
	"fmt"
	"strings"
)

type RuntimeError struct {
	token   Token
	message string
}

func (r RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", r.message, r.token.Line)
}

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) Interpret(expression Expr) error {
	value, err := i.evaluate(expression)
	if err != nil {
		return err
	}

	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) stringify(object any) string {
	if object == nil {
		return "nil"
	}
	if i.isFloat(object) {
		text := fmt.Sprintf("%f", object)
		if strings.HasSuffix(text, ".0") {
			text = text[:len(text)-2]
		}
	}
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.accept(i)
}

func (i *Interpreter) visitBinaryExpr(expr *Binary) (any, error) {
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

		return nil, i.error(expr.Operator, "Operands must be two numbers or two strings.")
	case Slash:
		err := i.checkNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}

		if right.(float64) == 0 {
			return nil, i.error(expr.Operator, "Can not divide by zero.")
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

func (i *Interpreter) visitGroupingExpr(expr *Grouping) (any, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) visitLiteralExpr(expr *Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) visitUnaryExpr(expr *Unary) (any, error) {
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

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	switch object.(type) {
	case bool:
		return object.(bool)
	}
	return true
}

func (i *Interpreter) isFloat(object any) bool {
	_, ok := object.(float64)
	return ok
}

func (i *Interpreter) isString(object any) bool {
	_, ok := object.(string)
	return ok
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) error {
	if i.isFloat(operand) {
		return nil
	}
	return i.error(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator Token, left any, right any) error {
	if i.isFloat(left) && i.isFloat(right) {
		return nil
	}
	return i.error(operator, "Operands must be numbers.")
}

func (i *Interpreter) error(token Token, message string) RuntimeError {
	return RuntimeError{token, message}
}
