package llvm

import (
	"errors"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mcjcloud/taurine/pkg/ast"
)

func (m LlvmModule) compileExpression(exp ast.Expression, block *ir.Block) (value.Value, error) {
	switch exp := exp.(type) {
	case *ast.StringLiteral:
		return m.compileStringLiteral(exp, block)
	case *ast.BooleanLiteral:
		return m.compileBooleanLiteral(exp, block)
	case *ast.NumberLiteral:
		return m.compileNumberLiteral(exp, block)
	case *ast.IntegerLiteral:
		return m.compileIntegerLiteral(exp, block)
	case *ast.Identifier:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.OperationExpression:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.VariableDecleration:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.AssignmentExpression:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.FunctionCall:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.GroupExpression:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.ArrayExpression:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.FunctionLiteral:
		// TODO
		return nil, errors.New("unimplemented")
	case *ast.ObjectLiteral:
		// TODO
		return nil, errors.New("unimplemented")
	default:
		return nil, fmt.Errorf("unrecognized expression: %s", exp)
	}
}

// compileStringLiteral adds the string literal to global scope and returns a reference
func (m LlvmModule) compileStringLiteral(exp *ast.StringLiteral, block *ir.Block) (value.Value, error) {
	nullTerminated := append([]byte(exp.Value), 0x00)
	return constant.NewCharArrayFromString(string(nullTerminated)), nil
}

func (m LlvmModule) compileBooleanLiteral(exp *ast.BooleanLiteral, block *ir.Block) (value.Value, error) {
	if exp.Value {
		return constant.True, nil
	}
	return constant.False, nil
}

func (m LlvmModule) compileNumberLiteral(exp *ast.NumberLiteral, block *ir.Block) (value.Value, error) {
	return constant.NewFloat(types.Float, exp.Value), nil
}

func (m LlvmModule) compileIntegerLiteral(exp *ast.IntegerLiteral, block *ir.Block) (value.Value, error) {
	if !exp.Value.IsInt64() {
		return nil, fmt.Errorf("cannot represent %s as int64", exp.Value)
	}
	return constant.NewInt(types.I64, exp.Value.Int64()), nil
}

// func (m LlvmModule) compileVariableDecleration(exp *ast.VariableDecleration, block *ir.Block) (value.Value, error) {

// }
