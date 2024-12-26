package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mukundvijay123/KCloud/types"
	"github.com/mukundvijay123/KCloud/utils"
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

func (h *Handler) RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc("/user/login", h.UserLogin).Methods(http.MethodPost)
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	//Limit the size of the request
	r.Body = http.MaxBytesReader(w, r.Body, types.MaxUserRequestSize)

	//Parsing the form
	if err := parseFormData(w, r); err != nil {
		return
	}

	// Define compulsory fields for each action
	requiredFields := []string{"company_username", "company_password"}
	if err := checkRequiredFields(w, r, requiredFields); err != nil {
		return
	}

	// Extract form values
	companyUsername := r.FormValue("company_username")
	companyPassword := r.FormValue("company_password")

	//Checking if the strings are valid
	if !utils.IsValidName(companyUsername) && !utils.IsNotEmptySring(companyPassword) {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	//Initialising a company user
	UserToBeLoggedIn := types.Company{
		Username:        companyUsername,
		CompanyPassword: companyPassword,
	}

	//Useing UILogin functiion to determine if login was successful
	loginSuccess, err := UILogin(&UserToBeLoggedIn, h.db)
	if err != nil {
		if err.Error() == "failed to login" {
			fmt.Println(err)
			http.Error(w, "Error logging in", http.StatusInternalServerError)
			w.Write([]byte(LoginUnsuccessfulMessage))
			return
		} else {
			http.Error(w, "Incorrect Credentials", http.StatusAccepted)
			w.Write([]byte(LoginUnsuccessfulMessage))
			return
		}
	} else if loginSuccess {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(LoginSuccessMessage))
	} else {
		http.Error(w, "Incorrect Credentials", http.StatusAccepted)
		w.Write([]byte(LoginUnsuccessfulMessage))
		return
	}

	// The control flow will never reach here

}

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
