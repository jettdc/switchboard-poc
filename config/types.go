package config

type Config struct {
	Server ServerConfig  `yaml:"server"`
	Routes []RouteConfig `yaml:"routes"`
}

type ServerConfig struct {
	Host    string       `yaml:"host"`
	Port    uint32       `yaml:"port"`
	Pubsub  PubSubConfig `yaml:"pubsub"`
	SSL     *SSLConfig   `yaml:"ssl,omitempty"`
	EnvFile string       `yaml:"env-file"`
}

type PubSubConfig struct {
	Provider string `yaml:"provider"`
}

type RouteConfig struct {
	Endpoint string         `yaml:"endpoint"`
	Topics   []string       `yaml:"topics"`
	Plugins  *PluginsConfig `yaml:"plugins,omitempty"`
}

type PluginsConfig struct {
	MiddlewarePaths []string            `yaml:"middleware"`
	EnrichmentPaths []string            `yaml:"message-enrichment"`
	Middleware      []*MiddlewarePlugin `yaml:"-"`
	Enrichment      []*EnrichmentPlugin `yaml:"-"`
}

type SSLConfig struct {
	Mode     string `yaml:"mode"`
	CertPath string `yaml:"cert"`
	KeyPath  string `yaml:"key"`
}
