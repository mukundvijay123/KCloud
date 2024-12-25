package metadata

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mukundvijay123/KCloud/types"
	"github.com/mukundvijay123/KCloud/utils"
)

// Handler struct for metadata service
// Anything the handler function needs access to should be enclosed in this struct
// For example a db connection ,redis conection,logger object etc
type Handler struct {
	db *sql.DB
}

// Function to get a new handler
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// Function to register routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/resources/company", h.CompanyHandler).Methods(http.MethodPost)
	router.HandleFunc("/resources/group", h.GrpHandler).Methods(http.MethodPost)
	router.HandleFunc("/resources/device", h.DeviceHandler).Methods(http.MethodPost)
}

// Handler function for endpoint "/resources/company"
func (h *Handler) CompanyHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	// Parse the form in request body
	if err := parseFormData(w, r); err != nil {
		return
	}

	// Extract the action from the form
	action := r.FormValue("action")
	if action != "create" && action != "delete" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	// Define compulsory fields for each action
	var requiredFields []string
	if action == "create" {
		requiredFields = []string{"company_name", "company_username", "company_password"}
	} else if action == "delete" {
		requiredFields = []string{"company_username", "company_password"}
	}

	// Enforce compulsory fields
	if err := checkRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	// Extract form values
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")

	if !utils.IsValidName(companyUsername) {
		http.Error(w, "Invalid field: company_username", http.StatusBadRequest)
		return
	}

	// If action is "create", also validate company_name
	if action == "create" {
		companyName := r.FormValue("company_name")
		if !utils.IsValidName(companyName) {
			http.Error(w, "Invalid company name", http.StatusBadRequest)
			return
		}
		// Creating company object for processing
		companyToBeProcessed := types.Company{
			CompanyName:     companyName,
			Username:        companyUsername,
			CompanyPassword: companyPassword,
		}

		// Process creation
		if err := ProvisionCompany(&companyToBeProcessed, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if action == "delete" {
		// Creating company object for processing
		companyToBeProcessed := types.Company{
			Username:        companyUsername,
			CompanyPassword: companyPassword,
		}

		// Process deletion
		if err := DeleteCompany(&companyToBeProcessed, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Company operation successful"))

}

// Handler for endpoint "/resources/group"
func (h *Handler) GrpHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	// Parse the form in request body
	if err := parseFormData(w, r); err != nil {
		return
	}

	// Extract the action from the form
	action := r.FormValue("action")
	if action != "create" && action != "delete" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	//Define compulsory fields
	var requiredFields []string
	if action == "create" {
		requiredFields = []string{"company_username", "company_password", "group_name"}
	} else if action == "delete" {
		requiredFields = []string{"company_username", "company_password", "group_name"}
	}

	//Enforce compulsory fiels for each action
	if err := checkRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	//Extract form values
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")
	groupName := r.FormValue("group_name")

	//Validate form information
	if !utils.IsValidName(companyUsername) || !utils.IsValidName(groupName) || !utils.IsNotEmptySring(companyPassword) {
		fmt.Println(utils.IsNotEmptySring(companyPassword), utils.IsNotEmptySring(companyPassword), utils.IsValidName(groupName))
		http.Error(w, "Invalid field(s)", http.StatusBadRequest)
		return
	}

	//actiion => Create
	if action == "create" {
		comapanyToBeProcessed := types.Company{
			Username:        companyUsername,
			CompanyPassword: companyPassword,
		}

		grpToBeProcessed := types.Grp{
			GroupName: groupName,
		}
		//Create group
		if err := ProvisionGroup(&comapanyToBeProcessed, &grpToBeProcessed, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if action == "delete" { //Delete action
		comapanyToBeProcessed := types.Company{
			Username:        companyUsername,
			CompanyPassword: companyPassword,
		}

		grpToBeProcessed := types.Grp{
			GroupName: groupName,
		}
		//Delete group
		if err := DeleteGroup(&comapanyToBeProcessed, &grpToBeProcessed, h.db); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Group operation successful"))

}

// Handler function for endpoint "/resource/device"
func (h *Handler) DeviceHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	// Parse the form in request body
	if err := parseFormData(w, r); err != nil {
		return
	}

	// Extract the action from the form
	action := r.FormValue("action")
	if action != "create" && action != "delete" {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	// Define compulsory fields for each action
	var requiredFields []string
	if action == "create" {
		requiredFields = []string{"company_username", "company_password", "group_name", "device_name", "telemetry_data_schema", "device_description", "device_type", "longitude", "latitude"}
	} else if action == "delete" {
		requiredFields = []string{"company_username", "company_password", "group_name", "device_name"}
	}

	// Enforce compulsory fields
	if err := checkRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	//Extracting Information from the device
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")
	groupName := r.FormValue("group_name")
	deviceName := r.FormValue("device_name")

	//Validating information from the form
	if !utils.IsValidName(companyUsername) || !utils.IsValidName(groupName) || !utils.IsValidName(deviceName) || !utils.IsNotEmptySring(companyPassword) {
		http.Error(w, "invalid field(s)", http.StatusBadRequest)
		return
	}
	//Create a device
	if action == "create" {
		telemetryDataSchema := json.RawMessage(r.FormValue("telemetry_data_schema"))
		deviceDescription := r.FormValue("device_description")
		deviceType := r.FormValue("device_type")
		longitude, latitude, err := utils.LocationValidator(r.FormValue("longitude"), r.FormValue("latitude")) //Parsing and extracting latitude and longitude from the form
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//Validating more form information
		if !utils.IsNotEmptySring(deviceType) || !utils.IsNotEmptySring(deviceDescription) {
			http.Error(w, "Invalid/Empty device description or device type", http.StatusBadRequest)
			return
		}

		companyToBeProcessed := types.Company{
			Username:        companyUsername,
			CompanyPassword: companyPassword,
		}

		groupToBeProcessed := types.Grp{
			GroupName: groupName,
		}

		deviceToBeProcessed := types.Device{
			DeviceName:          deviceName,
			DeviceType:          deviceType,
			DeviceDescription:   deviceDescription,
			TelemetryDataSchema: telemetryDataSchema,
			Longitude:           longitude,
			Latitude:            latitude,
		}

		//Creating a device
		if err := ProvisionDevice(&groupToBeProcessed, &companyToBeProcessed, &deviceToBeProcessed, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if action == "delete" { //Delete action
		companyToBeProcessed := types.Company{
			Username:        companyUsername,
			CompanyPassword: companyPassword,
		}

		groupToBeProcessed := types.Grp{
			GroupName: groupName,
		}

		deviceToBeProcessed := types.Device{
			DeviceName: deviceName,
		}

		//Deleting the device
		if err := DeleteDevice(&companyToBeProcessed, &groupToBeProcessed, &deviceToBeProcessed, h.db); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Device operation successful"))
}

// Below functions are helper functions for the handler functions
// Function to parse form data
func parseFormData(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(types.MaxMetadataRequestSize) // setting max size
	if err != nil {
		if err.Error() == "http: request body too large" {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		} else {
			fmt.Print("Problem!!")
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
		}
		return err
	}
	return nil
}

// Function to check required fields
func checkRequiredFields(w http.ResponseWriter, r *http.Request, requiredFields []string) error {
	for _, field := range requiredFields {
		value := r.FormValue(field)
		if value == "" {
			http.Error(w, fmt.Sprintf("Missing required field: %s", field), http.StatusBadRequest)
			return fmt.Errorf("missing required field: %s", field)
		}

	}
	return nil
}
