package resolver

type FunctionType string

const (
	FunctionTypeNone     FunctionType = "none"
	FunctionTypeFunction FunctionType = "function"
	FunctionTypeLambda   FunctionType = "lambda"
)

type LoopType string

const (
	LoopTypeNone  LoopType = "none"
	LoopTypeWhile LoopType = "while"
)
