package env

import "github.com/iamBharatManral/atom.git/cmd/internal/result"

type Environment struct {
	symbols map[string]result.Result
	parent  *Environment
}

func New(parent *Environment) *Environment {
	return &Environment{
		symbols: make(map[string]result.Result),
		parent:  parent,
	}
}

func (e *Environment) Get(symbol string) (result.Result, bool) {
	if e == nil {
		return result.Result{}, false
	}
	value, ok := e.symbols[symbol]
	if ok {
		return value, ok
	}
	return e.parent.Get(symbol)
}

func (e *Environment) Set(symbol string, result result.Result) {
	e.symbols[symbol] = result
}

func (e *Environment) Symbols() map[string]result.Result {
	return e.symbols
}
