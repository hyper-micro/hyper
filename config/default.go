package config

var defaultConfig *Config

func InitDefault() error {
	defaultConfig = New()
	return defaultConfig.LoadPaths("./conf")
}

func Default() *Config {
	if defaultConfig == nil {
		panic("config: uninitialized")
	}
	return defaultConfig
}
