package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type authPlugin string

func (authPlugin) Process(r *http.Request) error {
	fmt.Println("Extracting token")
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return fmt.Errorf("No token provided", err)
	}

	fmt.Println("Verifying token")
	token := tokenCookie.Value

	res, exists := os.LookupEnv("AUTH_ADDRESS")
	if !exists {
		res = "localhost:8081"
	}

	values := map[string]string{"token": token}
	json_data, err := json.Marshal(values)

	if err != nil {
		return fmt.Errorf("Could not marshal token", err)
	}

	resp, err := http.Post("http://"+res+"/validate", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("Validation request failed", err)
	}

	return nil
}

var MiddlewarePlugin authPlugin
