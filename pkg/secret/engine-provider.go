package secret

// EngineProvider is an interface of a secret engine provider
type EngineProvider interface {
	// Seal secrets
	Seal(secrets *Secrets) error
	// unseal secrets
	UnSeal(secrets *Secrets) error
}

// NewSecretsEngineProvider create new secret engine provider
func NewSecretsEngineProvider(uri string) (EngineProvider, error) {
	if uri == "" {
		return nil, nil
	}
	return newLocalSecretsEngineProvider(uri)
}
