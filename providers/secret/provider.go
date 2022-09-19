package secret

type SecretProvider interface {
	Get(name string) (string, error)
}
