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

// OperationStack is a stack of operation
type OperationStack []Operation

// NewOperation create new operation
func NewOperation(name OperationName, args ...string) *Operation {
	return &Operation{
		Name: name,
		Args: args,
	}
}
