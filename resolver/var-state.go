package resolver

import "interp/token"

type varState struct {
	defined  bool
	resolved bool
	token    token.Token
}

func (v *varState) define() {
	v.defined = true
}

func (v *varState) resolve() {
	v.resolved = true
}
