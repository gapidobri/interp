package interpreter

import (
	"fmt"
	"interp/errors"
	. "interp/token"
	"strings"
)

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

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	if f, ok := object.(float64); ok && f == 0 {
		return false
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
	return errors.NewRuntimeError(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator Token, left any, right any) error {
	if i.isFloat(left) && i.isFloat(right) {
		return nil
	}
	return errors.NewRuntimeError(operator, "Operands must be numbers.")
}
