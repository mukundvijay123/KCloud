package metadatarouter

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mukundvijay123/KCloud/metadata"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type passwordChangeRequest struct {
	CompanyId   uuid.UUID `json:"company_id"`
	NewPassword string    `json:"new_password"`
	OldPassword string    `json:"old_password"`
}

func (m *MetadataRouter) AddRoutes() error {
	err := m.addCompanyRoutes()
	if err != nil {
		return err
	}
	return nil
}

// addCompanyRoutes adds the signup route
func (m *MetadataRouter) addCompanyRoutes() error {
	companySubRouter := m.Router.PathPrefix("/api/user").Subrouter()
	companySubRouter.HandleFunc("/signup", m.signupHandler).Methods("POST") // use Methods("POST")
	companySubRouter.HandleFunc("/login", m.loginHandler).Methods("POST")

	postLoginRouter := companySubRouter.NewRoute().Subrouter()
	postLoginRouter.Use(m.JWTMiddleWare.JWTMiddleware)
	postLoginRouter.HandleFunc("/deleteCompany", m.DeleteComapnyHandler).Methods("POST")
	postLoginRouter.HandleFunc("/changePassword", m.updatePasswordHandle).Methods("POST")
	postLoginRouter.HandleFunc("/createGroup", m.createGroupHandler).Methods("POST")
	postLoginRouter.HandleFunc("/deleteGroup", m.deleteGroupHandler).Methods("POST")
	postLoginRouter.HandleFunc("/getGroup", m.getGroupByIDHandler).Methods("GET")
	postLoginRouter.HandleFunc("/getGroups", m.getGroupsHandler).Methods("GET")
	postLoginRouter.HandleFunc("/createDevice", m.createDeviceHandler).Methods("POST")
	return nil
}

// signupHandler handles user signup and creates a company entry
func (m *MetadataRouter) signupHandler(w http.ResponseWriter, r *http.Request) {
	var company metadata.Company

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save company using metadata store
	if err := m.MdataStore.CreateCompany(&company); err != nil {
		http.Error(w, "Failed to create company: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Company created successfully",
		"id":      company.ID.String(), // assuming Company struct has an ID field
	})
}

func (m *MetadataRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Verify credentials
	ok, err := m.MdataStore.VerifyCompany(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Error verifying credentials", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Lookup company to get ID (for JWT payload)
	company, err := m.MdataStore.GetCompanyByUsername(req.Username)
	if err != nil {
		http.Error(w, "Failed to fetch company info", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := m.JWTMiddleWare.GenerateToken(company.ID.String())
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
	w.WriteHeader(http.StatusCreated)

}

func (m *MetadataRouter) DeleteComapnyHandler(w http.ResponseWriter, r *http.Request) {
	var company metadata.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if company.Username == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = m.MdataStore.DeleteCompany(&company)
	if err != nil {
		http.Error(w, "error deleting the company", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func (m *MetadataRouter) updatePasswordHandle(w http.ResponseWriter, r *http.Request) {
	var req passwordChangeRequest
	var company *metadata.Company
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	company, err = m.MdataStore.GetCompanyByID(req.CompanyId.String())
	if err != nil {
		http.Error(w, "Cant find company", http.StatusInternalServerError)
		return
	}
	valid, err := m.MdataStore.VerifyCompany(company.Username, req.OldPassword)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}

	err = m.MdataStore.UpdatePassword(company, req.NewPassword)
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

}

func (m *MetadataRouter) createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var grp metadata.Grp
	err := json.NewDecoder(r.Body).Decode(&grp)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if grp.CompanyID == uuid.Nil || grp.GroupName == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	grp.NoOfDevices = 0
	err = m.MdataStore.CreateGroup(&grp)
	if err != nil {
		http.Error(w, "Error creating a group", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (m *MetadataRouter) deleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	var grp metadata.Grp
	err := json.NewDecoder(r.Body).Decode(&grp)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if grp.CompanyID == uuid.Nil || grp.GroupName == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	grp.NoOfDevices = 0
	err = m.MdataStore.DeleteGroup(&grp)
	if err != nil {
		http.Error(w, "Error creating a group", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (m *MetadataRouter) getGroupsHandler(w http.ResponseWriter, r *http.Request) {
	companyID := r.URL.Query().Get("company_id")
	if companyID == "" {
		http.Error(w, "company_id is required", http.StatusBadRequest)
		return
	}

	// Fetch groups from DB
	groups, err := m.MdataStore.ListGroupsByCompany(companyID)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		m.logger.Println("[getGroupsHandler] error:", err)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (m *MetadataRouter) getGroupByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract "id" query param
	groupID := r.URL.Query().Get("id")
	if groupID == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch group from store
	group, err := m.MdataStore.GetGroupByID(groupID)
	if err != nil {
		http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
		m.logger.Println("[getGroupByIDHandler] error:", err)
		return
	}

	if group == nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (m *MetadataRouter) createDeviceHandler(w http.ResponseWriter, r *http.Request) {
	
}
