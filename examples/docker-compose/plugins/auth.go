package main

import (
	"net/http"
)

type authPlugin string

func (authPlugin) Process(r *http.Request) error {
	// Check if the auth token is actually valid
	// Make a call to login service
	// If auth token found and matches, return nil. Else: err
	return nil
}

var MiddlewarePlugin authPlugin
