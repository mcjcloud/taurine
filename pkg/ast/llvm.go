package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type LlvmSSA interface{}

func LogLlvmError(err error) LlvmSSA {
	fmt.Println(err)
	return nil
}

var SYSCALL_TABLE map[string]int32 = map[string]int32{
	"sys_write": 4,
}

func (e *EtchStatement) GenLlvm(m *ir.Module, block *ir.Block) {
	asmFn := ir.NewInlineAsm(types.NewPointer(types.NewFunc(types.I64)), "syscall", "=r,{rax},{rdi},{rsi},{rdx}")
	asmFn.SideEffect = true

	// TODO build string

	var strPtr, strLen value.Value

	block.NewCall(
		asmFn,
		constant.NewInt(types.I64, 0x2000004), // syscall number
		constant.NewInt(types.I64, 1),         // stdout
		strPtr,
		strLen,
	)
}
