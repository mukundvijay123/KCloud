package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// This function is used to check if a schema of a daata of a device is correct while device registration
// takes json.RawMessage as input and evaluates schema follows rules
// Rules are
// 1. KeyLength not  greater than 10 characters
// 2. datatypes allowed are string,bol,float,json
// 3. there is no provision for storing units (like Centigrade) in the schema
// 4. isValidSchema and isValidValue are helper functions
func IsValidTelemetrySchema(data json.RawMessage) (bool, error) {
	var jsonObj map[string]interface{}

	if err := json.Unmarshal(data, &jsonObj); err != nil {
		log.Println("Invalid JSON")
		return false, fmt.Errorf("not a valid json")
	}
	return isValidSchema(jsonObj)

}

// helper to Is IsValidTelemetrySchema()
func isValidSchema(data map[string]interface{}) (bool, error) {
	for key, value := range data {
		if len(key) > 10 {
			log.Printf("%v key is longer than 10 characters", key)
			return false, fmt.Errorf("key length error")
		}
		if !isValidValue(value) {
			log.Printf("invalid key value")
			return false, fmt.Errorf("value error")
		}
	}
	return true, nil

}

// helper to Is IsValidTelemetrySchema()
func isValidValue(value interface{}) bool {

	switch v := value.(type) {
	case string:
		if v == "string" || v == "bool" || v == "timestamp" || v == "float" {
			return true
		}
	case map[string]interface{}:
		return true
	default:
		return false
	}
	return false
}
