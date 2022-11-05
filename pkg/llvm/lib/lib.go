package lib

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// Module represents a module which may add new global functions
type Module interface {
	GetIRModule() *ir.Module
	AddGlobalFunction(string, value.Value) error
}
