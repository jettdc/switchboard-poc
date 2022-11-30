package main

import (
	"fmt"
	"net/http"
)

type otherMWPlugin string

func (otherMWPlugin) Process(r *http.Request) error {
	fmt.Println("Successfully loaded other middleware plugin")
	return nil
}

var MiddlewarePlugin otherMWPlugin
