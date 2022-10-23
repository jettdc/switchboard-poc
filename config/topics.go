package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func ParameterizeTopic(topic string, wildcards gin.Params) (string, error) {
	var parameterizedTopic = topic
	wildcard, i, valid := findWildcard(topic)

	// Parameterize each wildcard in the topic until there are no more
	for i != -1 {
		if !valid {
			return "", fmt.Errorf("invalid wildcards found in topic: %s", topic)
		}

		// Find what value to replace it with, if exists
		replaceWith, ok := "", false
		for _, wc := range wildcards {
			if wc.Key == wildcard[1:] {
				replaceWith, ok = wc.Value, true
				break
			}
		}

		// No value was found to replace the wildcard with
		if !ok {
			return "", fmt.Errorf("missing mapping for wildcard %s in topic %s", wildcard[1:], topic)
		}

		parameterizedTopic = strings.Replace(parameterizedTopic, wildcard, replaceWith, 1)

		// Find the next
		wildcard, i, valid = findWildcard(parameterizedTopic)
	}
	return parameterizedTopic, nil
}

// Parameterize topic with placeholders to see if there's any errors.
func ValidateTopic(topic string) error {
	var parameterizedTopic = topic
	wildcard, i, valid := findWildcard(topic)

	// Parameterize each wildcard in the topic until there are no more
	for i != -1 {
		if !valid {
			return fmt.Errorf("invalid wildcards found in topic: %s", topic)
		}

		parameterizedTopic = strings.Replace(parameterizedTopic, wildcard, "VALID_PARAM", 1)

		// Find the next
		wildcard, i, valid = findWildcard(parameterizedTopic)
	}

	return nil
}

// From gin github, adapted to not accept * characters as wildcards
// Search for a wildcard segment and check the name for invalid characters.
// Returns -1 as index, if no wildcard was found.
func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (catch-all)
		if c != ':' {
			continue
		}

		// Find end and check for invalid characters
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}
