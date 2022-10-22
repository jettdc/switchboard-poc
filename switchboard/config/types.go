package config

type Config struct {
	Server ServerConfig  `yaml:"server"`
	Routes []RouteConfig `yaml:"routes"`
}

type ServerConfig struct {
	Host string     `yaml:"host"`
	Port uint32     `yaml:"port"`
	SSL  *SSLConfig `yaml:"ssl,omitempty"`
}

type PubSubConfig struct {
	ConnectionString string `yaml:"connection-string"`
}

type RouteConfig struct {
	Endpoint string        `yaml:"endpoint"`
	Topics   []string      `yaml:"topics"`
	Plugins  PluginsConfig `yaml:"plugins,omitempty"`
}

type PluginsConfig struct {
	MiddlewarePaths []string `yaml:"middleware"`
	EnrichmentPaths []string `yaml:"message-enrichment"`
}

type SSLConfig struct {
	Mode     string `yaml:"mode"`
	CertPath string `yaml:"cert"`
	KeyPath  string `yaml:"key"`
}
