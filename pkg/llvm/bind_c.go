package llvm

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (m LlvmModule) compileLibC() error {
	if err := m.compilePuts(); err != nil {
		return err
	}
	return nil
}

func (m LlvmModule) compilePuts() (err error) {
	defer func() {
		err = fmt.Errorf("%v", recover())
	}()
	charParam := ir.NewParam("char", types.I8Ptr)
	charParam.Attrs = append(charParam.Attrs, enum.ParamAttrNoCapture)
	puts := ir.NewFunc("puts", types.I32, charParam)
	puts.FuncAttrs = append(puts.FuncAttrs, enum.FuncAttrNoUnwind)
	m.Funcs = append(m.Funcs, puts)
	m.globalFunctions["puts"] = puts
	return
}

func (m LlvmModule) puts(val value.Value, block *ir.Block) {
	// TODO: replace hardcoded 12 with string size
	s := val.String()

	strMem := block.NewAlloca(types.NewArray(uint64(len(s)), types.I8))
	block.NewStore(val, strMem)
	bc := block.NewBitCast(strMem, types.I8Ptr)
	block.NewCall(m.globalFunctions["puts"], bc)
}
