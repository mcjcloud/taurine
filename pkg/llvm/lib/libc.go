package lib

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func CompileLibc(m Module) error {
	if err := compilePuts(m); err != nil {
		return err
	}
	if err := compilePrintf(m); err != nil {
		return err
	}
	return nil
}

func compilePuts(m Module) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover: %v", recover())
		}
	}()

	charParam := ir.NewParam("", types.I8Ptr)
	charParam.Attrs = append(charParam.Attrs, enum.ParamAttrNoCapture)

	puts := m.GetIRModule().NewFunc("puts", types.I32, charParam)
	puts.FuncAttrs = append(puts.FuncAttrs, enum.FuncAttrNoUnwind)

	m.AddGlobalFunction("puts", puts)

	return
}

func compilePrintf(m Module) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover: %v", recover())
		}
	}()

	charParam := ir.NewParam("", types.I8Ptr)
	charParam.Attrs = append(charParam.Attrs, enum.ParamAttrNoAlias)
	charParam.Attrs = append(charParam.Attrs, enum.ParamAttrNoCapture)

	printf := m.GetIRModule().NewFunc("printf", types.I32, charParam)
	printf.Sig.Variadic = true

	m.AddGlobalFunction("printf", printf)

	return
}
