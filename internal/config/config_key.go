package config

type (
	// Config ...
	Config struct {
		Server   ServerConfig   `yaml:"server"`
		Database DatabaseConfig `yaml:"database"`
		API      APIConfig      `yaml:"api"`
		Swagger  SwaggerConfig  `yaml:"swagger"`
	}

	// ServerConfig ...
	ServerConfig struct {
		Port string `yaml:"port"`
	}

	// DatabaseConfig ...
	DatabaseConfig struct {
		Master string `yaml:"master"`
	}

	// APIConfig ...
	APIConfig struct {
		Auth string `yaml:"auth"`
	}

	SwaggerConfig struct {
		Host    string   `yaml:"host"`
		Schemes []string `yaml:"schemes"`
	}
)
