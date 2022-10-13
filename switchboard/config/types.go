package config

type Config struct {
	Server ServerConfig  `yaml:"server"`
	PubSub PubSubConfig  `yaml:"pubsub"`
	Routes []RouteConfig `yaml:"routes"`
}

type ServerConfig struct {
	Port uint32 `yaml:"port"`
}

type PubSubConfig struct {
	ConnectionString string `yaml:"connection-string"`
}

type RouteConfig struct {
	Endpoint string        `yaml:"endpoint"`
	Topic    string        `yaml:"topic"`
	Plugins  PluginsConfig `yaml:"plugins,omitempty"`
}

type PluginsConfig struct {
	MiddlewarePaths []string `yaml:"middleware"`
	EnrichmentPaths []string `yaml:"message-enrichment"`
}
