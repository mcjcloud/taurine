package llvm

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (m LlvmModule) puts(val *constant.CharArray, block *ir.Block) error {
	if m.globalFunctions["puts"] == nil {
		return fmt.Errorf("function not found: puts")
	}

	strMem := block.NewAlloca(val.Type())
	block.NewStore(val, strMem)

	bc := block.NewBitCast(strMem, types.I8Ptr)
	block.NewCall(m.globalFunctions["puts"], bc)

	return nil
}

func (m LlvmModule) printf(val *constant.CharArray, block *ir.Block) error {
	if _, ok := m.globalFunctions["printf"]; !ok {
		return fmt.Errorf("function not found: printf")
	}

	template := constant.NewCharArrayFromString("%s\x00")
	tempMem := block.NewAlloca(template.Type())
	block.NewStore(template, tempMem)
	tBc := block.NewBitCast(tempMem, types.I8Ptr)

	strMem := block.NewAlloca(val.Type())
	block.NewStore(val, strMem)
	strBc := block.NewBitCast(strMem, types.I8Ptr)

	block.NewCall(m.globalFunctions["printf"], tBc, strBc)

	return nil
}
