package contracts

// Source provides access to configuration values.
//
// This abstraction allows infra/config to load from env today, and from other
// sources (files, flags, secrets, etc.) in the future without changing loaders.
type Source interface {
	// Lookup returns the raw value for a given key.
	// ok is false when the key is not set.
	Lookup(key string) (value string, ok bool)
}
