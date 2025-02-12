package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	authentication "github.com/mukundvijay123/KCloud/middleware/auth"
	"github.com/mukundvijay123/KCloud/services/metadata"
	"github.com/mukundvijay123/KCloud/types"
	"github.com/mukundvijay123/KCloud/utils"
	"golang.org/x/crypto/bcrypt"
)

const LoginSuccessMessage = "Login successful"
const LoginUnsuccessfulMessage = "Login unsuccessful"

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

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/company/login", h.CompanyLogin).Methods(http.MethodPost)
	router.HandleFunc("/company/create", h.CompanyCreateHandler).Methods(http.MethodPost)
}

func (h *Handler) CompanyCreateHandler(w http.ResponseWriter, r *http.Request) {
	//Limit Request Body Size
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxMetadataRequestSize)

	//ParseFormData
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}

	//Enforcing cumposory fields
	requiredFields := []string{"company_name", "company_username", "company_password"}
	if err := utils.CheckRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	//extract values
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")
	companyName := r.FormValue("company_name")

	companyHashedPassword, err := bcrypt.GenerateFromPassword([]byte(companyPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Creating company object for processing
	companyToBeProcessed := types.Company{
		CompanyName:     companyName,
		Username:        companyUsername,
		CompanyPassword: string(companyHashedPassword),
	}

	// Initiating process
	if err := metadata.ProvisionCompany(&companyToBeProcessed, h.db); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Company operation successful"))

}

func (h *Handler) CompanyLogin(w http.ResponseWriter, r *http.Request) {
	//Put a max sze cap on the body
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxUserRequestSize)

	//Parse form Data
	if err := utils.ParseFormData(w, r); err != nil {
		return
	}

	//Extracting Details
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")

	companyToBeProcessed := types.Company{
		Username:        companyUsername,
		CompanyPassword: companyPassword,
	}

	if result, err := Login(&companyToBeProcessed, h.db); !result || err != nil {
		fmt.Println(err, result)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		token, err := authentication.GenenerateJWT(companyUsername)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
