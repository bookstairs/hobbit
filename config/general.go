package config

// Configuration represents node's configuration.
type Configuration interface {
	// TOML returns configuration string as the toml format.
	TOML() string
}
