package main

import (
	"fmt"
	"net/http"
)

type authPlugin string

func (authPlugin) Process(r *http.Request) error {
	fmt.Println("Successfully loaded plugin!")
	return fmt.Errorf("Successfully loaded plugin!")
}

var MiddlewarePlugin authPlugin
