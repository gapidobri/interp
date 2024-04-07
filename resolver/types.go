package resolver

type FunctionType string

const (
	FunctionTypeNone        FunctionType = "none"
	FunctionTypeFunction    FunctionType = "function"
	FunctionTypeInitializer FunctionType = "initializer"
	FunctionTypeMethod      FunctionType = "method"
)

type ClassType string

const (
	ClassTypeNone  ClassType = "none"
	ClassTypeClass ClassType = "class"
)
