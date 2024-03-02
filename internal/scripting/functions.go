package scripting

import (
	"fmt"

	"github.com/skx/evalfilter/v2/object"
)

type fnType = func(args []object.Object) object.Object

func (i *Interpreter) buildSingleArgFunction(op OperationName) fnType {
	return func(args []object.Object) object.Object {
		if len(args) != 1 {
			return &object.Void{}
		}
		arg := args[0].Inspect()
		operations := *i.operations
		operations = append(operations, *NewOperation(op, arg))
		i.operations = &operations
		return &object.Void{}
	}
}

func (i *Interpreter) buildNoArgFunction(op OperationName) fnType {
	return func(args []object.Object) object.Object {
		operations := *i.operations
		operations = append(operations, *NewOperation(op))
		i.operations = &operations
		return &object.Void{}
	}
}

// fnNoOp do nothing
func fnNoOp(args []object.Object) object.Object {
	return &object.Void{}
}

// fnPrint is the implementation of our `print` function.
func (i *Interpreter) fnPrint(args []object.Object) object.Object {
	for _, e := range args {
		i.logger.Debug().Str("fn", "print").Msg(e.Inspect())
	}
	return &object.Void{}
}

// fnPrintf is the implementation of our `printf` function.
func (i *Interpreter) fnPrintf(args []object.Object) object.Object {
	// We expect 1+ arguments
	if len(args) < 1 {
		return &object.Null{}
	}
	// Type-check
	if args[0].Type() != object.STRING {
		return &object.Null{}
	}
	// Get the format-string.
	fs := args[0].(*object.String).Value
	// Convert the arguments to something go's sprintf
	// code will understand.
	argLen := len(args)
	fmtArgs := make([]interface{}, argLen-1)
	// Here we convert and assign.
	for i, v := range args[1:] {
		fmtArgs[i] = v.ToInterface()
	}
	// Call the helper
	out := fmt.Sprintf(fs, fmtArgs...)
	i.logger.Debug().Str("fn", "printf").Msg(out)
	return &object.Void{}
}
