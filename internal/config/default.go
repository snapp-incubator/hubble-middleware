package config

func Default() Config {
	return Config{
		API: API{Port: 1381},
	}
}
