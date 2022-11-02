package main

import (
	"fmt"
	"net/http"
)

type otherMWPlugin string

func (otherMWPlugin) Process(r *http.Request) error {
	fmt.Println("Successfully loaded DIFFERENT plugin!")
	return fmt.Errorf("Successfully loaded DIFFERENT plugin!")
}

var MiddlewarePlugin otherMWPlugin
