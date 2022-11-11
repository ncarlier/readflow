package scripting

// OperationName is the operation name
type OperationName uint

const (
	// OpDrop to drop an article
	OpDrop OperationName = iota
	// OpTriggerWebhook to trigger an outgoing webhook
	OpTriggerWebhook
	// OpSendNotification to send a notification to all user devices
	OpSendNotification
	// OpSetText to set article text
	OpSetText
	// OpSetTitle to set article title
	OpSetTitle
	// OpSetCategory to set article category
	OpSetCategory
)

// Operation object
type Operation struct {
	Name OperationName
	Args []string
}

// GetFirstArg retrn first operation argument
func (op Operation) GetFirstArg() string {
	if len(op.Args) > 0 {
		return op.Args[0]
	}
	return ""
}

// OperationStack is a stack of operation
type OperationStack []Operation

// Contains test if an operation is part of the stack
func (ops OperationStack) Contains(op OperationName) bool {
	for _, v := range ops {
		if v.Name == op {
			return true
		}
	}
	return false
}

// NewOperation create new operation
func NewOperation(name OperationName, args ...string) *Operation {
	return &Operation{
		Name: name,
		Args: args,
	}
}
