package main

import (
	"fmt"
)

type enrichPizzaPlugin string

func (enrichPizzaPlugin) Process(s string) (string, error) {
	fmt.Println("Successfully loaded enrich plugin")
	return "DONE", nil
}

var EnrichmentPlugin enrichPizzaPlugin
