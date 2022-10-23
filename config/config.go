package config

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/jettdc/switchboard/u"
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot find specified config file at %s", path)
	}
	defer file.Close()

	config := &Config{}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("invalid yaml in config file: %s", err.Error())
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	u.Logger.Info("Successfully loaded the switchboard config.")

	return config, nil
}

func validateConfig(config *Config) error {
	if err := validateServerConfig(config.Server); err != nil {
		return err
	}

	for _, route := range config.Routes {
		// Make sure route parameterization is OK
		for _, topic := range route.Topics {
			// TODO: MAke sure that this makes sense
			if err := ValidateTopic(topic); err != nil {
				return err
			}
		}

		// Make sure that specified plugins exist and are valid
		if !structs.IsZero(route.Plugins) {
			if err := validatePlugins(route.Plugins); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateServerConfig(sc ServerConfig) error {
	// Makes sure that "server" exists
	if structs.IsZero(sc) {
		return fmt.Errorf("invalid server configuration")
	}

	if sc.Host == "" {
		return fmt.Errorf("missing host in config")
	}

	if structs.IsZero(sc.Pubsub) {
		return fmt.Errorf("missing pubsub configuration")
	}

	if sc.SSL != nil {
		if err := validateSSLConfig(sc.SSL); err != nil {
			return err
		}
	}

	return nil
}

func validateSSLConfig(sslConfig *SSLConfig) error {
	switch sslConfig.Mode {
	case "", "manual":
		if sslConfig.KeyPath == "" || sslConfig.CertPath == "" {
			return fmt.Errorf("must provide paths to ssl cert and key files with non-automatic ssl")
		}

		if _, err := os.Stat(sslConfig.KeyPath); err != nil {
			return fmt.Errorf("cannot find ssl key at path \"%s\"", sslConfig.KeyPath)
		}

		if _, err := os.Stat(sslConfig.CertPath); err != nil {
			return fmt.Errorf("cannot find ssl certificate at path \"%s\"", sslConfig.CertPath)
		}

		return nil
	case "auto":
		if sslConfig.KeyPath != "" || sslConfig.CertPath != "" {
			u.Logger.Warn("SSL mode configured to \"auto\" but cert or key was also provided. Defaulting to auto mode.")
		}
		return nil
	case "none":
		return nil
	default:
		return fmt.Errorf("invalid ssl mode \"%s\"", sslConfig.Mode)
	}
}

func validatePlugins(pluginsConfig PluginsConfig) error {
	if len(pluginsConfig.EnrichmentPaths) > 0 {
		enrichmentPluginsErr := validateEnrichmentPlugins(pluginsConfig.EnrichmentPaths)
		if enrichmentPluginsErr != nil {
			return enrichmentPluginsErr
		}
	}

	if len(pluginsConfig.MiddlewarePaths) > 0 {
		middlewarePluginsErr := validateMiddlewarePlugins(pluginsConfig.MiddlewarePaths)
		if middlewarePluginsErr != nil {
			return middlewarePluginsErr
		}
	}

	return nil
}

func validateEnrichmentPlugins(enrichmentPlugins []string) error {
	// Make sure the files exist
	for _, pluginPath := range enrichmentPlugins {
		_, err := os.Stat(pluginPath)
		if err != nil {
			return fmt.Errorf("cannot find enrichment plugin file at %s", pluginPath)
		}
	}

	// TODO: Validate that plugin is ok
	return nil
}

func validateMiddlewarePlugins(middlewarePlugins []string) error {
	// Make sure the files exist
	for _, pluginPath := range middlewarePlugins {
		_, err := os.Stat(pluginPath)
		if err != nil {
			return fmt.Errorf("cannot find middleware plugin file at %s", pluginPath)
		}
	}

	// TODO: Validate that plugin is ok
	return nil
}
