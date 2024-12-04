package provisioning

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mukundvijay123/KCloud/utils"
)

// Company struct represents the company table
type Company struct {
	ID              int    `json:"id"`
	CompanyName     string `json:"company_name"`
	Username        string `json:"username"`
	CompanyPassword string `json:"company_password"`
	NoOfGrps        int    `json:"no_of_grps"`
	NoOfDevices     int    `json:"no_of_devices"`
}

// Grp struct represents grp table
type Grp struct {
	ID          int    ` json:"id"`
	CompanyID   int    ` json:"company_id"`
	GroupName   string ` json:"group_name"`
	NoOfDevices int    ` json:"no_of_devices"`
}

// Device dtruct represents device table
type Device struct {
	ID                  int             `json:"id"`
	GrpID               int             `json:"grp_id"`
	CompanyID           int             `json:"company_id"`
	DeviceName          string          `json:"device_name"`
	DeviceDescription   string          `json:"device_description"`
	Longitude           float64         `json:"longitude"`
	Latitude            float64         `json:"latitude"`
	TelemetryDataSchema json.RawMessage `json:"telemetry_data_schema"`
}

// Create method inserts a new company record into the database
func (c *Company) ProvisionCompany(db *sql.DB) error {

	//Verifying if incoming credentials are correct
	if !utils.IsValidName(c.CompanyName) {
		log.Printf("%s : CompanyName cannot contain space and should only have alphanumeric characters", c.CompanyName)
		return fmt.Errorf("companyName cannot contain  spaces and should have only alphanumeric characters")
	}
	if !utils.IsValidName(c.Username) {
		log.Printf("%s :Username cannot consist of spaces and should have only alphanumeric characters", c.Username)
		return fmt.Errorf("username cannot consist of spaces and should have only alphanumeric characters")
	}
	if !utils.IsValidName(c.CompanyPassword) {
		log.Printf("Invalid Password")
		return fmt.Errorf("password cannot consist of spaces and should have only alphanumeric characters")
	}

	//
	var existingCompanyID int
	checkQuery := `SELECT id FROM company WHERE username = $1`
	err := db.QueryRow(checkQuery, c.Username).Scan(&existingCompanyID)

	if err != sql.ErrNoRows {
		if err != nil {
			return fmt.Errorf("failed to check username uniqueness: %v", err)
		}

		return fmt.Errorf("username %s is already taken", c.Username)
	}
	//Query to Insert company in companies table
	insertCompanyQuery := `INSERT INTO company (company_name, username, company_password, no_of_grps, no_of_devices) 
                    VALUES ($1, $2, $3, $4, $5) RETURNING id`
	c.NoOfGrps = 0
	c.NoOfDevices = 0
	err = db.QueryRow(insertCompanyQuery, c.CompanyName, c.Username, c.CompanyPassword, c.NoOfGrps, c.NoOfDevices).Scan(&c.ID)

	if err != nil {
		return fmt.Errorf("failed to provision company")
	}
	log.Printf("Company provisioned successfully with ID: %d", c.ID)

	return nil
}

// Method is used to provision a new group in the database
func (g *Grp) ProvisionGroup(c *Company, db *sql.DB) error {

	//verify if group credentials are valid
	if !utils.IsValidName(g.GroupName) {
		return fmt.Errorf("group name cannot conatain spaces , should obnly contain alphanumeric characters")
	}

	//verifying if company with given  usernname exists
	var existingCompanyID int
	companyCheckQuery := `SELECT id FROM company WHERE username = $1`
	err := db.QueryRow(companyCheckQuery, c.Username).Scan(&existingCompanyID)

	if err != nil {
		return fmt.Errorf("failed to find company with username %v", c.Username)
	}

	// Check if a group with the same name already exists for the company
	groupCheckQuery := `SELECT id FROM grp WHERE company_id = $1 AND group_name = $2`
	var existingGroupID int

	err = db.QueryRow(groupCheckQuery, existingCompanyID, g.GroupName).Scan(&existingGroupID)
	if err == nil {
		// Group already exists
		return fmt.Errorf("a group with the name '%v' already exists for company %v", g.GroupName, existingCompanyID)
	} else if err != sql.ErrNoRows {
		// Handle any other errors (e.g., database issues)
		return fmt.Errorf("failed to check for existing group: %v", err)
	}

	//All preconditions are met to provision a group
	//starting the database transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	//Adding Grp to grp table
	groupInsertQuery := `INSERT INTO grp (company_id, group_name, no_of_devices) 
		VALUES ($1, $2, $3) RETURNING id`

	g.CompanyID = existingCompanyID
	g.NoOfDevices = 0
	err = tx.QueryRow(groupInsertQuery, g.CompanyID, g.GroupName, g.NoOfDevices).Scan(&g.ID)
	if err != nil {
		return fmt.Errorf("error inserting group: %v", err)
	}

	//Updating no of grps in company table
	updateCompanyQuery := `UPDATE company SET no_of_grps = no_of_grps + 1 WHERE id = $1`
	_, err = tx.Exec(updateCompanyQuery, g.CompanyID)
	if err != nil {
		return fmt.Errorf("error updating company while provisioning group: %v", err)
	}

	//commiting final transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction for provisioning a group: %v", err)
	}

	return nil

}

func (d *Device) ProvisionDevice(g *Grp, c *Company) error {

	return nil
}
