package config

import (
	"fmt"
	"github.com/fatih/structs"
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

	return config, nil
}

func validateConfig(config *Config) error {
	// Return an error if any required (fields without omitempty) fields are missing
	hasMissingFields := structs.HasZero(config)
	if hasMissingFields {
		return fmt.Errorf("missing required fields in config file")
	}

	for _, route := range config.Routes {
		// Make sure route parameterization is OK
		for _, topic := range route.Topics {
			if err := ValidateTopic(topic); err != nil {
				fmt.Println(topic)
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
