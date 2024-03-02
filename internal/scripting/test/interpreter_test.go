package test

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/ncarlier/readflow/internal/scripting"
	"github.com/ncarlier/readflow/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestInterpreterSimpleScript(t *testing.T) {
	script := `
if ("bar" in Tags) {
	print("sending notification because Tags contains bar");
	sendNotification();
	return false;
}
return true;
`

	interpreter, err := scripting.NewInterpreter(script, logger.With().Logger())
	assert.Nil(t, err)
	assert.NotNil(t, interpreter)
	input := scripting.ScriptInput{
		Tags: []string{"foo", "bar"},
	}
	operations, err := interpreter.Exec(context.Background(), input)
	assert.Nil(t, err)
	assert.Len(t, operations, 2)
	operation := operations[0]
	assert.Equal(t, scripting.OpSendNotification, operation.Name, "invalid operation")
	operation = operations[1]
	assert.Equal(t, scripting.OpDrop, operation.Name, "invalid operation")

	input.Tags = []string{"foo"}
	operations, err = interpreter.Exec(context.Background(), input)
	assert.Nil(t, err)
	assert.Len(t, operations, 0)
}

func TestInterpreterRaceCondition(t *testing.T) {
	script := `
printf("title=%s", Title);
index = int(Title)
if (index %2 == 0) {
	return true;
}
return false;
`

	interpreter, err := scripting.NewInterpreter(script, logger.With().Logger())
	assert.Nil(t, err)
	assert.NotNil(t, interpreter)

	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)
	results := make([]scripting.OperationStack, count)

	for i := 0; i < count; i++ {
		go func(index int) {
			input := scripting.ScriptInput{Title: strconv.Itoa(index)}
			ops, _ := interpreter.Exec(context.Background(), input)
			results[index] = ops
			wg.Done()
		}(i)
	}

	wg.Wait()
	for i := 0; i < count; i++ {
		result := results[i]
		assert.Len(t, result, i%2, "invalid test case #%d", i)
	}
}
