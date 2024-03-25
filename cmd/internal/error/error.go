package error

import (
	"fmt"

	"github.com/iamBharatManral/atom.git/cmd/internal/result"
)

func TypeMismatchError(first, second any) result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintf("error: mismatch types %v(%T) and %v(%T)", first, first, second, second),
	}
}

func UnsupportedTypeError(left any, operator string) result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintf("error: unsupported type %T for %s", left, operator),
	}
}

func UnsupportedOperatorError(operator string) result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintf("error: unsupported operator %s", operator),
	}
}

func UnsupportedTokensError() result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintln("error: unsupported tokens"),
	}
}

func DivisonByZeroError() result.Result {
	return result.Result{
		Type:  "error",
		Value: "error: division by zero",
	}
}

func SyntaxError(message string) result.Result {
	return result.Result{
		Type:  "error",
		Value: message,
	}
}

func UndefinedError(symbol string) result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintf("error: undefined symbol '%s'", symbol),
	}
}

func UnsupportedOperation(msg string) result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintf("error: %s", msg),
	}
}
func NotEnoughArguments(msg string) result.Result {
	return result.Result{
		Type:  "error",
		Value: fmt.Sprintf("error: %s", msg),
	}
}
