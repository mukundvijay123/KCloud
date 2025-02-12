package metadata

import (
	"database/sql"
	"encoding/json"
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
	router.HandleFunc("/company/delete", h.CompanyDeleteHandler).Methods(http.MethodPost)
	router.HandleFunc("/group/create", h.GrpCreateHandler).Methods(http.MethodPost)
	router.HandleFunc("/group/delete", h.GrpDeleteHandler).Methods(http.MethodPost)
	router.HandleFunc("/device/create", h.DeviceCreateHandler).Methods(http.MethodPost)
	router.HandleFunc("/device/delete", h.DeviceDeleteHandler).Methods(http.MethodPost)
}

func (h *Handler) CompanyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//Limit Request Body Size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	//ParseFormData
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}

	//Enforcing cumposory fields
	requiredFields := []string{"company_username"}
	if err := utils.CheckRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	//extract values
	companyUsername := r.FormValue("company_username")
	if !utils.IsValidName(companyUsername) {
		http.Error(w, "Invalid field: company_username", http.StatusBadRequest)
		return
	}
	// Creating company object for processing
	companyToBeProcessed := types.Company{
		Username: companyUsername,
	}

	// Initiating process
	if err := DeleteCompany(&companyToBeProcessed, h.db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Company operation successful"))

}

func (h *Handler) GrpCreateHandler(w http.ResponseWriter, r *http.Request) {
	//Limit Request Body Size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	//ParseFormData
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}

	//Enforcing cumposory fields
	requiredFields := []string{"company_username", "group_name"}
	if err := utils.CheckRequiredFields(w, r, requiredFields); err != nil {
		return
	}
	companyUsername := r.FormValue("company_username")
	groupName := r.FormValue("group_name")
	//Validate form information
	if !utils.IsValidName(companyUsername) || !utils.IsValidName(groupName) {
		http.Error(w, "Invalid field(s)", http.StatusBadRequest)
		return
	}

	comapanyToBeProcessed := types.Company{
		Username: companyUsername,
	}

	grpToBeProcessed := types.Grp{
		GroupName: groupName,
	}

	if err := ProvisionGroup(&comapanyToBeProcessed, &grpToBeProcessed, h.db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Group operation successful"))

}

func (h *Handler) GrpDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//Limit Request Size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	//Parse FormData
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}
	requiredFields := []string{"company_username", "group_name"}
	if err := utils.CheckRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	//Extract form values
	companyUsername := r.FormValue("company_username")
	groupName := r.FormValue("group_name")

	//Validate form information
	if !utils.IsValidName(companyUsername) || !utils.IsValidName(groupName) {
		http.Error(w, "Invalid field(s)", http.StatusBadRequest)
		return
	}
	comapanyToBeProcessed := types.Company{
		Username: companyUsername,
	}

	grpToBeProcessed := types.Grp{
		GroupName: groupName,
	}
	//Delete group
	if err := DeleteGroup(&comapanyToBeProcessed, &grpToBeProcessed, h.db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Group operation successful"))

}

func (h *Handler) DeviceCreateHandler(w http.ResponseWriter, r *http.Request) {
	//Limit Request Size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	//Parse FormData
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}
	requiredFields := []string{"company_username", "group_name", "device_name", "telemetry_data_schema", "device_description", "device_type", "longitude", "latitude"}
	if err := utils.CheckRequiredFields(w, r, requiredFields); err != nil {
		return
	}
	companyUsername := r.FormValue("company_username")
	groupName := r.FormValue("group_name")
	deviceName := r.FormValue("device_name")
	telemetryDataSchema := json.RawMessage(r.FormValue("telemetry_data_schema"))
	deviceDescription := r.FormValue("device_description")
	deviceType := r.FormValue("device_type")
	longitude, latitude, err := utils.LocationValidator(r.FormValue("longitude"), r.FormValue("latitude")) //Parsing and extracting latitude and longitude from the form
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Validating more form information
	if !utils.IsValidName(companyUsername) || !utils.IsValidName(groupName) || !utils.IsValidName(deviceName) || !utils.IsNotEmptySring(deviceType) || !utils.IsNotEmptySring(deviceDescription) {
		http.Error(w, "Invalid/Empty device description or device type", http.StatusBadRequest)
		return
	}

	companyToBeProcessed := types.Company{
		Username: companyUsername,
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

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Device operation successful"))

}

func (h *Handler) DeviceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	//Limit Request Size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	//Parse FormData
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}
	requiredFields := []string{"company_username", "group_name", "device_name"}
	// Enforce compulsory fields
	if err := utils.CheckRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	//Extracting Information from the device
	companyUsername := r.FormValue("company_username")
	groupName := r.FormValue("group_name")
	deviceName := r.FormValue("device_name")

	//Validating Form info
	if !utils.IsValidName(companyUsername) || !utils.IsValidName(groupName) || !utils.IsValidName(deviceName) {
		http.Error(w, "invalid field(s)", http.StatusBadRequest)
		return
	}

	companyToBeProcessed := types.Company{
		Username: companyUsername,
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

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Device operation successful"))

}
