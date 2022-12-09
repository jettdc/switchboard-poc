package main

import (
	"fmt"
	"net/http"
)

type authPlugin string

func (authPlugin) Process(r *http.Request) error {
	// Check if the auth token is actually valid
	// Make a call to login service
	// If auth token found and matches, return nil. Else: err
	fmt.Println("Extracting token")
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return fmt.Errorf("No token provided", err)
	}

	fmt.Println("Verifying token")
	token := tokenCookie.Value
	fmt.Println(token)
	return nil
}

var MiddlewarePlugin authPlugin
