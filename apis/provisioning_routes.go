package provisioning

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mukundvijay123/KCloud/router"
	"github.com/mukundvijay123/KCloud/utils"
)

func company(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//If the method is POST
	case http.MethodPost:
		//limiting the sixe of the request
		r.Body = http.MaxBytesReader(w, r.Body, router.ProvisioningMaxBodySize)

		//Reading the request body
		rawBody, err := readRequestBody(w, r)
		if err != nil {
			return
		}
		//Rewind body for parsing form
		r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

		//Parse the form in request body
		if err := parseFormData(w, r); err != nil {
			return
		}

		//Enforce cumpolssory fields
		requiredFields := []string{"company_username", "company_password"}
		if err := checkRequiredFields(w, r, requiredFields); err != nil {
			return
		}

		// Parse the raw body (for action field)
		var requestBody map[string]string
		if err := parseRawBody(w, rawBody, &requestBody); err != nil {
			return
		}
		if val, ok := requestBody["action"]; !ok && (val == "delete" || val == "create") {
			http.Error(w, "invalid request", http.StatusBadRequest)
		}

		//creating company object for processing
		var companyToBeProcessed Company
		companyToBeProcessed.CompanyName = r.FormValue("comapany_username")
		companyToBeProcessed.CompanyPassword = r.FormValue("comapany_password")

	default:
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}

func readRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return nil, err
	}
	return rawBody, nil
}

func parseFormData(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
		}
		return err
	}
	return nil
}

func checkRequiredFields(w http.ResponseWriter, r *http.Request, requiredFields []string) error {
	for _, field := range requiredFields {
		if !utils.IsValidName(r.FormValue(field)) {
			http.Error(w, fmt.Sprintf("Missing required field: %s", field), http.StatusBadRequest)
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}

func parseRawBody(w http.ResponseWriter, rawBody []byte, requestBody *map[string]string) error {
	err := json.Unmarshal(rawBody, requestBody)
	if err != nil {
		http.Error(w, "Unable to parse raw body for action", http.StatusBadRequest)
		return err
	}
	return nil
}
