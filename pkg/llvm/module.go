package llvm

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

type LlvmModule struct {
	*ir.Module
	globalFunctions map[string]value.Value

	ir []byte
}

func (m LlvmModule) GetIRModule() *ir.Module {
	return m.Module
}

func (m LlvmModule) AddGlobalFunction(name string, v value.Value) error {
	if _, ok := m.globalFunctions[name]; ok {
		return fmt.Errorf("global function redeclared with name '%s'", name)
	}
	m.globalFunctions[name] = v

	return nil
}

// read LLVM IR from this module
func (m LlvmModule) Read(p []byte) (n int, err error) {
	for n = 0; n < len(p); n += 1 {
		if n >= len(m.ir) {
			break
		}
		p[n] = m.ir[n]
	}
	return
}

// write LLVM IR to the module
func (m LlvmModule) Write(p []byte) (n int, err error) {
	m.ir = make([]byte, len(p))
	copy(m.ir, p)

	return len(m.ir), nil
}
