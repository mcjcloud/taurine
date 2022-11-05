package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func print(val value.Value, block *ir.Block) {
	asmFn := types.NewFunc(types.I64 /*, types.I64, types.I64, types.I64, types.I64*/)
	inlineAsm := ir.NewInlineAsm(types.NewPointer(asmFn), "syscall", "=rgpm,{rax},{rdi},{rsi},{rdx}")
	inlineAsm.SideEffect = true

	strPtr := getPointer(block, val)
	// strLen := length(block, val)

	block.NewCall(
		inlineAsm,
		constant.NewInt(types.I64, 0x2000004), // syscall number
		constant.NewInt(types.I64, 1),         // stdout
		strPtr,
		constant.NewInt(types.I64, int64(11)),
	)
}

func getPointer(block *ir.Block, src value.Value) value.Value {
	if _, ok := src.Type().(*types.PointerType); ok {
		l := block.NewGetElementPtr(pointerType(src), src, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
		return block.NewLoad(pointerType(l), l)
	}
	return block.NewExtractValue(src, 1)
}

// pointerType returns a pointer type for a given value
func pointerType(src value.Value) types.Type {
	return src.Type().(*types.PointerType).ElemType
}

func length(block *ir.Block, src value.Value) value.Value {
	if _, ok := src.Type().(*types.PointerType); ok {
		l := block.NewGetElementPtr(pointerType(src), src, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		return block.NewLoad(pointerType(l), l)
	}
	return block.NewExtractValue(src, 0)
}
