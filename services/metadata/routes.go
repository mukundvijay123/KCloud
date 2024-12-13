package metadata

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mukundvijay123/KCloud/types"
	"github.com/mukundvijay123/KCloud/utils"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/resources/company", h.CompanyHandler).Methods(http.MethodPost)
}

func (h *Handler) CompanyHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	// Reading the request body
	rawBody, err := readRequestBody(w, r)
	if err != nil {
		return
	}

	// Rewind body for parsing form
	r.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Parse the form in request body
	if err := parseFormData(w, r); err != nil {
		return
	}

	// Enforce compulsory fields
	requiredFields := []string{"company_name", "company_username", "company_password", "action"}
	if err := checkRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	// Extract form values
	companyName := r.FormValue("company_name")
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")
	action := r.FormValue("action")

	// Validate action
	if action != "create" && action != "delete" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	// Creating company object for processing
	var companyToBeProcessed types.Company
	companyToBeProcessed.CompanyName = companyName
	companyToBeProcessed.Username = companyUsername
	companyToBeProcessed.CompanyPassword = companyPassword

	// Process based on action
	var processErr error
	if action == "create" {
		processErr = ProvisionCompany(&companyToBeProcessed, h.db)
	} else {
		processErr = DeleteCompany(&companyToBeProcessed, h.db)
	}

	if processErr != nil {
		http.Error(w, processErr.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Company operation successful"))
}

func parseFormData(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
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
		value := r.FormValue(field)
		if value == "" {
			http.Error(w, fmt.Sprintf("Missing required field: %s", field), http.StatusBadRequest)
			return fmt.Errorf("missing required field: %s", field)
		}

		// Validate username and password
		if field == "company_username" || field == "company_password" {
			if !utils.IsValidName(value) {
				http.Error(w, fmt.Sprintf("Invalid field: %s", field), http.StatusBadRequest)
				return fmt.Errorf("invalid field: %s", field)
			}
		}
	}
	return nil
}

func readRequestBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return nil, err
	}
	return rawBody, nil
}
