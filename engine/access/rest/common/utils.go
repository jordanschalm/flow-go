package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// SliceToMap converts a slice of strings into a map where each string
// in the slice becomes a key in the map with the value set to true.
func SliceToMap(values []string) map[string]bool {
	valueMap := make(map[string]bool, len(values))
	for _, v := range values {
		valueMap[v] = true
	}
	return valueMap
}

func ParseBody(raw io.Reader, dst interface{}) error {
	dec := json.NewDecoder(raw)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("request body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("request body contains unknown field %s", fieldName)
		case errors.Is(err, io.EOF):
			return fmt.Errorf("request body must not be empty")
		default:
			return err
		}
	}

	if dst == nil {
		return fmt.Errorf("request body must not be empty")
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("request body must only contain a single JSON object")
	}

	return nil
}

// ParseInterfacesToString converts a slice of interface{} to a slice of strings.
//
// No errors are expected during normal operations.
func ParseInterfacesToString(value interface{}) ([]string, error) {
	// Check if value is []string directly
	if strSlice, ok := value.([]string); ok {
		return strSlice, nil
	}

	// Check if value is []interface{}
	interfaceSlice, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' must be an array", value)
	}

	// Convert []interface{} to []string
	result := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("'%s' must be an array of strings", value)
		}
		result[i] = str
	}

	return result, nil
}
