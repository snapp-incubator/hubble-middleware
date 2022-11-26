package config

func Default() Config {
	return Config{
		API: API{Port: 1381},
		Dex: Dex{Secret: "HUBBLE_CLIENT_SECRET"},
	}
}
