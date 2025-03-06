package dgraph

type dgraphConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ConfigProvider interface {
	LoadConfig(container *dgraphConfig) error
}

func loadConfig(provider ConfigProvider) dgraphConfig {

	var config dgraphConfig

	err := provider.LoadConfig(&config)
	if err != nil {
		// TODO: handle error
	}
	return config
}
