package secret

// Action upon a secret
type Action uint

const (
	// Seal action
	Seal Action = iota
	// UnSeal action
	UnSeal
)

// EngineProvider is an interface of a secret engine provider
type EngineProvider interface {
	// Seal secrets
	Seal(secrets *Secrets) error
	// unseal secrets
	UnSeal(secrets *Secrets) error
	// Apply action on secrets
	Apply(action Action, secrets *Secrets) error
}

// NewSecretsEngineProvider create new secret engine provider
func NewSecretsEngineProvider(uri string) (EngineProvider, error) {
	if uri == "" {
		return nil, nil
	}
	return newLocalSecretsEngineProvider(uri)
}
