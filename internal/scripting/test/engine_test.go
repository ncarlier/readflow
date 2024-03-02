package test

import (
	"context"
	"testing"

	"github.com/ncarlier/readflow/internal/scripting"
	"github.com/stretchr/testify/assert"
)

func TestEngineSimpleScript(t *testing.T) {
	script := `
if (Title == "foo") {
	printf("sending notification because Title = %s", Title);
	sendNotification();
	return false;
}
return true;
`
	engine := scripting.NewScriptEngine(10)
	assert.NotNil(t, engine)

	input := scripting.ScriptInput{
		Title: "foo",
	}

	operations, err := engine.Exec(context.Background(), script, input)
	assert.Nil(t, err)
	assert.Len(t, operations, 2)

	input.Title = "bar"
	operations, err = engine.Exec(context.Background(), script, input)
	assert.Nil(t, err)
	assert.Len(t, operations, 0)
}
