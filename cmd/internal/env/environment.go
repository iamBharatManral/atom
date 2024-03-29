package env

import "github.com/iamBharatManral/atom.git/cmd/internal/result"

type Environment struct {
	symbols map[string]result.Result
}

func New() *Environment {
	return &Environment{
		symbols: make(map[string]result.Result),
	}
}

func (e *Environment) Get(symbol string) (result.Result, bool) {
	value, ok := e.symbols[symbol]
	return value, ok
}

func (e *Environment) Set(symbol string, result result.Result) {
	e.symbols[symbol] = result
}

func (e *Environment) Symbols() map[string]result.Result {
	return e.symbols
}
