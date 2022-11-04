package main

import (
	"fmt"
	"net/http"
)

type authPlugin string

func (authPlugin) Process(r *http.Request) error {
	fmt.Println("Successfully loaded auth plugin")
	return nil
}

var MiddlewarePlugin authPlugin
