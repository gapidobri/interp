package interpreter

type Class struct {
	Name    string
	methods map[string]*Function
}

func NewClass(name string, methods map[string]*Function) *Class {
	return &Class{name, methods}
}

func (c *Class) String() string {
	return c.Name
}

func (c *Class) findMethod(name string) *Function {
	return c.methods[name]
}

func (c *Class) call(interpreter *Interpreter, arguments []any) (any, error) {
	instance := NewInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		_, err := initializer.bind(instance).call(interpreter, arguments)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (c *Class) arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}
