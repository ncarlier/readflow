package scripting

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
	"github.com/skx/evalfilter/v2"
)

// Interpreter is a script interpreter
type Interpreter struct {
	eval       *evalfilter.Eval
	operations *OperationStack
	mu         sync.Mutex
	logger     zerolog.Logger
}

// NewInterpreter create new script interpreter
func NewInterpreter(script string, logger zerolog.Logger) (*Interpreter, error) {
	eval := evalfilter.New(script)
	if err := eval.Prepare(); err != nil {
		return nil, fmt.Errorf("unable to compile provided script: %w", err)
	}

	operations := OperationStack{}
	interpreter := &Interpreter{
		eval:       eval,
		operations: &operations,
		logger:     logger,
	}

	// init the interpreter
	interpreter.init()

	return interpreter, nil
}

func (i *Interpreter) init() {
	// deactivate unwanted functions
	i.eval.AddFunction("getenv", fnNoOp)
	// alter builtins functions
	i.eval.AddFunction("print", i.fnPrint)
	i.eval.AddFunction("printf", i.fnPrintf)
	// add custom functions
	i.eval.AddFunction("triggerWebhook", i.buildSingleArgFunction(OpTriggerWebhook))
	i.eval.AddFunction("sendNotification", i.buildNoArgFunction(OpSendNotification))
	i.eval.AddFunction("setCategory", i.buildSingleArgFunction(OpSetCategory))
	i.eval.AddFunction("setTitle", i.buildSingleArgFunction(OpSetTitle))
	i.eval.AddFunction("setText", i.buildSingleArgFunction(OpSetTitle))
	i.eval.AddFunction("setHTML", i.buildSingleArgFunction(OpSetHTML))
	i.eval.AddFunction("markAsRead", i.buildNoArgFunction(OpMarkAsRead))
	i.eval.AddFunction("markAsToRead", i.buildNoArgFunction(OpMarkAsToRead))
	i.eval.AddFunction("disableGlobalNotification", i.buildNoArgFunction(OpDisableGlobalNotification))
}

// Exec a script by the interpreter
func (i *Interpreter) Exec(ctx context.Context, input ScriptInput) (OperationStack, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.eval.SetContext(ctx)
	i.operations = &OperationStack{}
	if result, err := i.eval.Run(input); err != nil {
		return nil, fmt.Errorf("unable to execute script: %w", err)
	} else if !result {
		operations := append(*i.operations, *NewOperation(OpDrop))
		i.operations = &operations
	}
	return *i.operations, nil
}
