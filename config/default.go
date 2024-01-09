package config

var defaultConfig Config

func InitDefault() error {
	conf, err := New(PathTypePath, "./conf")
	if err != nil {
		return err
	}
	defaultConfig = conf
	return nil
}

func Default() Config {
	if defaultConfig == nil {
		panic("config: uninitialized")
	}
	return defaultConfig
}
