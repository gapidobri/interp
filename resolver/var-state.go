package resolver

import "interp/token"

type varState struct {
	defined  bool
	resolved bool
	token    token.Token
}

func (v *varState) define() {
	if v == nil {
		return
	}
	v.defined = true
}

func (v *varState) resolve() {
	if v == nil {
		return
	}
	v.resolved = true
}
