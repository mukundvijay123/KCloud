package utils

import (
	"fmt"
	"log"
	"net/http"
)

func ParseFormData(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		}
		log.Printf("Error Parsing Form")
		http.Error(w, "Unable to parse For", http.StatusBadRequest)
		return err
	}
	return nil
}

// Function to check required fields
func CheckRequiredFields(w http.ResponseWriter, r *http.Request, requiredFields []string) error {
	for _, field := range requiredFields {
		value := r.FormValue(field)
		if value == "" {
			http.Error(w, fmt.Sprintf("Missing required field: %s", field), http.StatusBadRequest)
			return fmt.Errorf("missing required field: %s", field)
		}

	}
	return nil
}
