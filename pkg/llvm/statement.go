package llvm

import (
	"errors"
	"fmt"

	"github.com/llir/llvm/ir/constant"

	"github.com/llir/llvm/ir"
	"github.com/mcjcloud/taurine/pkg/ast"
)

func (m LlvmModule) compileStatement(stmt ast.Statement, block *ir.Block) error {
	switch stmt := stmt.(type) {
	case *ast.EtchStatement:
		return m.compileEtchStatement(stmt, block)
	case *ast.ReadStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.ExpressionStatement:
		_, err := m.compileExpression(stmt.Expression, block)
		return err
	case *ast.BlockStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.IfStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.ForLoopStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.WhileLoopStatement:
		// TODO
		return errors.New("unimplemented")
	case *ast.ReturnStatement:
		// TODO
		return errors.New("unimplemented")
	default:
		return fmt.Errorf("unknown statement %s", stmt)
	}
}

func (m LlvmModule) compileEtchStatement(stmt *ast.EtchStatement, block *ir.Block) error {
	var values []*constant.CharArray
	for _, exp := range stmt.Expressions {
		v, err := m.compileExpression(exp, block)
		if err != nil {
			return err
		}
		strVal, err := wrapInCharArray(v)
		if err != nil {
			return err
		}
		values = append(values, strVal)
	}

	space := constant.NewCharArrayFromString(" \x00")
	for i, v := range values {
		m.printf(v, block)
		if i < len(values)-1 {
			m.printf(space, block)
		}
	}

	newline := constant.NewCharArrayFromString("\n\x00")
	m.printf(newline, block)

	return nil
}
