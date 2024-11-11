package utils

import (
	"encoding/json"
	"fmt"
)

// JSONFormat formats the provided data as pretty-printed JSON with indentation
// Returns a formatted JSON string or an error if formatting fails
func JSONFormat(data interface{}) (string, error) {
	// Attempt to marshal data into JSON with indentation
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// Return a more informative error if JSON marshalling fails
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}

	// Convert the bytes to a string and return
	return string(bytes), nil
}
