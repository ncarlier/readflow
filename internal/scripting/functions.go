package scripting

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/readflow/pkg/defaults"
	"github.com/skx/evalfilter/v2/object"
	"golang.org/x/net/html/charset"
)

type fnType = func(args []object.Object) object.Object

var VOID = &object.Void{}
var NULL = &object.Null{}

const MAX_RESPONSE_SIZE = 1 << 20 // 1Mb

func ErrorString(err error) *object.String {
	return &object.String{Value: err.Error()}
}

func (i *Interpreter) buildSingleArgFunction(op OperationName) fnType {
	return func(args []object.Object) object.Object {
		if len(args) != 1 {
			return VOID
		}
		arg := args[0].Inspect()
		operations := *i.operations
		operations = append(operations, *NewOperation(op, arg))
		i.operations = &operations
		return VOID
	}
}

func (i *Interpreter) buildNoArgFunction(op OperationName) fnType {
	return func(args []object.Object) object.Object {
		operations := *i.operations
		operations = append(operations, *NewOperation(op))
		i.operations = &operations
		return VOID
	}
}

// fnNoOp do nothing
func fnNoOp(args []object.Object) object.Object {
	return VOID
}

// fnPrint is the implementation of our `print` function.
func (i *Interpreter) fnPrint(args []object.Object) object.Object {
	for _, e := range args {
		i.logger.Debug().Str("fn", "print").Msg(e.Inspect())
	}
	return VOID
}

// fnPrintf is the implementation of our `printf` function.
func (i *Interpreter) fnPrintf(args []object.Object) object.Object {
	// We expect 1+ arguments
	if len(args) < 1 {
		return NULL
	}
	// Type-check
	if args[0].Type() != object.STRING {
		return NULL
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
	return VOID
}

// fnFetch is the implementation of our `fetch` function.
func (i *Interpreter) fnFetch(args []object.Object) object.Object {
	// we expect 1 argument
	if len(args) == 0 {
		return NULL
	}
	// type-check
	if args[0].Type() != object.STRING {
		return NULL
	}
	// get the target URL
	rawurl := args[0].(*object.String).Value
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return ErrorString(err)
	}
	if u.Scheme != "https" {
		return NULL
	}

	// do HTTP request
	req, err := http.NewRequest("GET", u.String(), http.NoBody)
	if err != nil {
		return ErrorString(err)
	}
	req.Header.Set("User-Agent", defaults.UserAgent)
	i.logger.Debug().Str("fn", "fetch").Str("url", rawurl).Msg("do HTTP request...")
	res, err := defaults.HTTPClient.Do(req)
	if err != nil {
		return ErrorString(err)
	}
	defer res.Body.Close()

	// validate HTTP response
	if res.StatusCode != 200 {
		return ErrorString(fmt.Errorf("invalid status code: %s", res.Status))
	}
	contentType := res.Header.Get("Content-type")
	if !strings.HasPrefix(contentType, "text/html") {
		return ErrorString(fmt.Errorf("invalid content-type: %s", contentType))
	}

	// read body response
	body, err := charset.NewReader(res.Body, contentType)
	if err != nil {
		return ErrorString(err)
	}
	body = io.LimitReader(body, MAX_RESPONSE_SIZE)
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, body); err != nil {
		return ErrorString(err)
	}

	return &object.String{
		Value: buf.String(),
	}
}
