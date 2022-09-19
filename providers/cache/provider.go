package cache

type CacheProvider interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	SetEx(key string, value string, expiration int) error
}
