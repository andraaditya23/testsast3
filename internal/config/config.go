package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
)

const (
	envDevelopment = "development"
	envStaging     = "staging"
	envProduction  = "production"
)

type option struct {
	configFile string
}

// Init ...
func Init(opts ...Option) error {
	opt := &option{
		configFile: getDefaultConfigFile(),
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	out, err := ioutil.ReadFile(opt.configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(out, &config)
}

// Option ...
type Option func(*option)

// WithConfigFile ...
func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

func getDefaultConfigFile() string {

	configPath := "./files/etc/example/example.development.yaml"
	namespace, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")

	env := string(namespace)
	if os.Getenv("GOPATH") == "" {
		configPath = "./example.development.yaml"
	}

	if env != "" {
		switch env {
		case envStaging:
			configPath = "./example.staging.yaml"
		case envProduction:
			configPath = "./example.production.yaml"
		default:
			configPath = "./example.development.yaml"
		}
	}
	return configPath
}

// Get ...
func Get() *Config {
	return config
}
