package evaluator

import (
	"errors"
	"fmt"

	"github.com/mcjcloud/taurine/pkg/ast"
)

func evaluateExpression(exp ast.Expression, scope *Scope) (ast.Expression, error) {
	switch t := exp.(type) {
	case *ast.OperationExpression:
		return evaluateOperation(t, scope)
	case *ast.Identifier:
		return scope.Get(t.Name), nil
	case *ast.VariableDecleration:
		return evaluateVariableDecleration(t, scope)
	case *ast.AssignmentExpression:
		return evaluateAssignmentExpression(t, scope)
	case *ast.FunctionCall:
		return evaluateFunctionCall(t, scope)
	case *ast.GroupExpression:
		return evaluateExpression(t.Expression, scope)
	case *ast.ArrayExpression:
		return evaluateArrayExpression(t, scope)
	case *ast.FunctionLiteral:
		return evaluateFunctionLiteral(t, scope)
	case *ast.ObjectLiteral:
		return evaluateObjectLiteral(t, scope)
	default:
		return exp, nil
	}
}

func evaluateVariableDecleration(decl *ast.VariableDecleration, scope *Scope) (ast.Expression, error) {
	val, err := evaluateExpression(decl.Value, scope)
	if err != nil {
		return nil, err
	}

	if scope.Variables[decl.Symbol] != nil {
		return nil, fmt.Errorf("variable '%s' already exists", decl.Symbol)
	}
	scope.Set(decl.Symbol, val)

	return val, nil
}

func evaluateAssignmentExpression(asn *ast.AssignmentExpression, scope *Scope) (ast.Expression, error) {
	// make sure the identifier exists
	if scope.Get(asn.Identifier.Name) == nil {
		return nil, fmt.Errorf("'%s' was not declared", asn.Identifier.Name)
	}

	val, err := evaluateExpression(asn.Value, scope)
	if err != nil {
		return nil, err
	}

	// update the scope and return the evaluated value
	scope.Set(asn.Identifier.Name, val)
	return val, nil
}

func evaluateArrayExpression(arr *ast.ArrayExpression, scope *Scope) (ast.Expression, error) {
	exp := &ast.ArrayExpression{Expressions: make([]ast.Expression, len(arr.Expressions))}
	for i, el := range arr.Expressions {
		val, err := evaluateExpression(el, scope)
		if err != nil {
			return nil, err
		}
		exp.Expressions[i] = val
	}
	return exp, nil
}

func evaluateFunctionCall(call *ast.FunctionCall, scope *Scope) (ast.Expression, error) {
	// TODO: make this cleaner, maybe move built-in functions someplace else
	if id, ok := call.Function.(*ast.Identifier); ok && id.Name == "len" {
		if len(call.Arguments) != 1 {
			return nil, errors.New("len takes only one argument")
		}
		return builtInLen(call.Arguments[0], scope)
	} else if ok && id.Name == "int" {
		return builtInInt(call.Arguments[0], scope)
	}

	// must be a non-built-in function
	fn, err := evaluateExpression(call.Function, scope)
	if err != nil {
		return nil, err
	}

	// expect that the expression evaluates to ScopedFunction
	scopedFn, ok := fn.(*ScopedFunction)
	if !ok {
		return nil, errors.New("called expression did not evaluate to function")
	}

	// check that the number of parameters are correct
	if len(scopedFn.Function.Parameters) != len(call.Arguments) {
		return nil, fmt.Errorf("expected '%d' arguments but got '%d' for call to '%s'", len(scopedFn.Function.Parameters), len(call.Arguments), call.Function)
	}

	// evaluate arguments and populate scope
	for i, argExp := range call.Arguments {
		exp, err := evaluateExpression(argExp, scope)
		if err != nil {
			return nil, err
		}

		dType := ast.Symbol(scopedFn.Function.Parameters[i].SymbolType)
		arg, err := conformDataType(dType, exp)
		if err != nil {
			return nil, err
		}

		scopedFn.Scope.Set(scopedFn.Function.Parameters[i].Symbol, arg)
	}
	// execute statements
	if err := executeStatement(scopedFn.Function.Body, scopedFn.Scope); err != nil {
		return nil, err
	}
	return scopedFn.Scope.ReturnValue, nil
}

func evaluateFunctionLiteral(fnVal *ast.FunctionLiteral, scope *Scope) (ast.Expression, error) {
	// if evaluating a FunctionLiteral, wrap it in the current scope
	// this allows that scope to be accessed during execution
	sf := &ScopedFunction{
		Scope:    scope,
		Function: fnVal,
	}

	// if there is a symbol name, store the function in scope
	// TODO: eventually I should distinguish between functinos and anon functions..
	// right now, you could name a variable function and it could be stored twice
	if fnVal.Symbol != "" {
		scope.Set(fnVal.Symbol, sf)
	}
	return sf, nil
}

func evaluateObjectLiteral(objExp *ast.ObjectLiteral, scope *Scope) (ast.Expression, error) {
	// if evaluating an object literal, evaluate each of it's properties
	for k, v := range objExp.Value {
		newExp, err := evaluateExpression(v, scope)
		if err != nil {
			return nil, err
		}
		objExp.Value[k] = newExp
	}
	return objExp, nil
}
