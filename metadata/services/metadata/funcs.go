package metadata

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mukundvijay123/KCloud/metadata/types"
	"github.com/mukundvijay123/KCloud/metadata/utils"
)

// Create method inserts a new company record into the database
func ProvisionCompany(c *types.Company, db *sql.DB) error {

	//Verifying if incoming credentials are correct
	if !utils.IsValidName(c.CompanyName) {
		log.Printf("%s : CompanyName cannot contain space and should only have alphanumeric characters", c.CompanyName)
		return fmt.Errorf("companyName cannot contain  spaces and should have only alphanumeric characters")
	}
	if !utils.IsValidName(c.Username) {
		log.Printf("%s :Username cannot consist of spaces and should have only alphanumeric characters", c.Username)
		return fmt.Errorf("username cannot consist of spaces and should have only alphanumeric characters")
	}
	if !utils.IsNotEmptySring(c.CompanyPassword) {
		log.Printf("Invalid Password")
		return fmt.Errorf("password cannot consist of spaces and should have only alphanumeric characters")
	}

	//starting a database transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction")
	}
	defer tx.Rollback()

	//provisioning a comapny
	var existingCompanyID int
	checkQuery := `SELECT id FROM company WHERE username = $1`
	err = tx.QueryRow(checkQuery, c.Username).Scan(&existingCompanyID)

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
	err = tx.QueryRow(insertCompanyQuery, c.CompanyName, c.Username, c.CompanyPassword, c.NoOfGrps, c.NoOfDevices).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed to provision company")
	}

	//Create partitioned table for  telemetry data
	telemetryCreateQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS data_company_%d 
		PARTITION OF data
		FOR VALUES IN (%d);
	`, c.ID, c.ID)

	_, err = tx.Exec(telemetryCreateQuery)
	if err != nil {
		return fmt.Errorf("failed to create partition table for company %d: %v", c.ID, err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Company provisioned successfully with ID: %d", c.ID)

	return nil
}

//Below functions are for provisioning devices ,groups and companies

// Method is used to provision a new group in the database
func ProvisionGroup(c *types.Company, g *types.Grp, db *sql.DB) error {

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
	log.Printf("Group provisioned successfully with ID: %d", g.ID)

	return nil

}

func ProvisionDevice(g *types.Grp, c *types.Company, d *types.Device, db *sql.DB) error {
	//verify if device credentials are valid
	if !utils.IsValidName(d.DeviceName) {
		log.Printf("%s : DeviceName cannot contain spaces and should only contain alphanumeric characters", d.DeviceName)
		return fmt.Errorf("deviceName cannot contain spaces and should only contain alphanumeric characters")
	}
	if !utils.IsNotEmptySring(d.DeviceType) {
		log.Printf("Invalid Device type")
		return fmt.Errorf("invalid device type")
	}
	//checking if company  and group  id is   valid
	var existingCompanyID int
	companyCheckQuery := `SELECT id FROM company WHERE username = $1`
	err := db.QueryRow(companyCheckQuery, c.Username).Scan(&existingCompanyID)
	if err != nil {
		return fmt.Errorf("failed to find company with username %v", c.Username)
	}

	var existingGroupID int
	groupCheckQuery := `SELECT id FROM grp WHERE company_id = $1 AND group_name = $2`
	err = db.QueryRow(groupCheckQuery, existingCompanyID, g.GroupName).Scan(&existingGroupID)
	if err != nil {
		return fmt.Errorf("failed to find group %v in company %v,error:%v", g.GroupName, c.CompanyName, err)

	}
	//verify if device_name is unique per group
	uniqueDeviceQuery := `SELECT EXISTS( SELECT 1 
						FROM devices 
						WHERE company_id = $1
						AND group_id = $2
						AND device_name = $3
						);`
	var exists bool
	err = db.QueryRow(uniqueDeviceQuery, existingCompanyID, existingGroupID, d.DeviceName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if device with name %v exists in group %v,db error", d.DeviceName, g.GroupName)
	} else if exists {
		return fmt.Errorf("device with name %v exists in group %v", d.DeviceName, g.GroupName)
	}

	//verify json schema
	isValid, err := utils.IsValidTelemetrySchema(d.TelemetryDataSchema)
	if err != nil || !isValid {
		return fmt.Errorf("error while parsing telemetry schema")
	}
	//All checks for provisioning a device are completed
	//Provision device
	d.CompanyID = existingCompanyID
	d.GrpID = existingGroupID

	/*
		Insert code here to verify latitude and longitude values
	*/

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("coudnt begin database transaction : %v", err)
	}
	defer tx.Rollback()

	// Insert the new device into the 'devices' table
	insertDeviceQuery := `
	INSERT INTO devices (
		device_name, 
		group_id, 
		company_id, 
		telemetry_data_schema, 
		device_description, 
		device_type, 
		longitude, 
		latitude
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	var deviceID int
	err = tx.QueryRow(insertDeviceQuery, d.DeviceName, d.GrpID, d.CompanyID, d.TelemetryDataSchema, d.DeviceDescription, d.DeviceType, d.Longitude, d.Latitude).Scan(&deviceID)
	if err != nil {
		return fmt.Errorf("couldnt insert device: %v", err)
	}

	//Update the number of devices in the grp table
	updateCompanyQuery := `
	UPDATE company
	SET no_of_devices=no_of_devices+1
	WHERE id = $1
	`
	_, err = tx.Exec(updateCompanyQuery, d.CompanyID)
	if err != nil {
		return fmt.Errorf("error updating company")
	}

	updateGrpQuery := `
	UPDATE grp
	SET	no_of_devices=no_of_devices+1
	WHERE id = $1
	`
	_, err = tx.Exec(updateGrpQuery, d.GrpID)
	if err != nil {
		return fmt.Errorf("error updating grp")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error commiting transaction")
	}
	log.Printf("device with id %v provisioned", deviceID)
	return nil
}

//below functions are for de-provisioning a device ,group or a company

// method below is to delete a company
func DeleteCompany(c *types.Company, db *sql.DB) error {
	deleteCompanyQuery := `DELETE FROM company WHERE username = $1`
	result, err := db.Exec(deleteCompanyQuery, c.Username)
	if err != nil {
		return fmt.Errorf("error deleting company with username %v:%v", deleteCompanyQuery, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected for company with username %v: %v", c.Username, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no company found with username %v", c.Username)
	}
	log.Printf("company with username %v deleted", c.Username)
	return nil
}

// Method below to delete a group
func DeleteGroup(c *types.Company, g *types.Grp, db *sql.DB) error {
	//starting db transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting a database transaction")
	}
	defer tx.Rollback()
	//fetching grp_id
	fetchGrpIdQuery := `
	SELECT id,no_of_devices
	FROM grp
	WHERE company_id = (
		SELECT id
		FROM company
		WHERE username = $1
	)
	AND group_name = $2;`
	err = db.QueryRow(fetchGrpIdQuery, c.Username, g.GroupName).Scan(&g.ID, &g.NoOfDevices)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no group found for the given username and group name")
		} else {
			return fmt.Errorf("error executing query: %v", err)
		}
	}

	//updating number of devices and grps in company table
	updateCompanyQuery := `UPDATE company SET no_of_devices = no_of_devices - $1,
	no_of_grps = no_of_grps -1 
	WHERE username = $2`
	result, err := tx.Exec(updateCompanyQuery, g.NoOfDevices, c.Username)
	if err != nil {
		return fmt.Errorf("error updating devices in the database : %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while updating company owning  company_id %v: %v", g.CompanyID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no company found %v", g.CompanyID)
	}

	//Deleting group
	deleteGroupQuery := `DELETE FROM grp WHERE id = $1`
	result, err = tx.Exec(deleteGroupQuery, g.ID)
	if err != nil {
		return fmt.Errorf("error while deleting group with id %v, err:%v", g.ID, err)
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while deleting group with group_id %v: %v", g.ID, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no group found %v", g.ID)
	}

	//Commiing transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error while commiting transaction , %v", err)
	}
	log.Printf("group with id %v deleted", g.ID)

	return nil

}

// Method below to delete a device
func DeleteDevice(c *types.Company, g *types.Grp, d *types.Device, db *sql.DB) error {
	//start a database transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting database trainsaction")
	}
	defer tx.Rollback()

	//update no_of_devices in grp and company table
	updateCompanyQuery := `UPDATE company SET no_of_devices = no_of_devices -1 WHERE username = $1`
	result, err := tx.Exec(updateCompanyQuery, c.Username)
	if err != nil {
		return fmt.Errorf("error while updating comapny table , %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while checking rows affected in company table,%v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("comapny for device doesnt exist, err ")
	}

	//fetch grp id from gp name
	fetchGrpIdQuery := `
	SELECT id
	FROM grp
	WHERE company_id = (
		SELECT id
		FROM company
		WHERE username = $1
	)
	AND group_name = $2;`
	err = db.QueryRow(fetchGrpIdQuery, c.Username, g.GroupName).Scan(&g.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no group found for the given username and group name")
		} else {
			return fmt.Errorf("error executing query: %v", err)
		}
	}

	updateGroupQuery := `UPDATE grp SET no_of_devices = no_of_devices -1 WHERE id = $1`
	result, err = tx.Exec(updateGroupQuery, g.ID)
	if err != nil {
		return fmt.Errorf("error while updating no_of_devices in grp table")
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while checking rows affected in company table,%v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("grp for device does not exist : %v", err)
	}

	//delete the device from device table
	deleteDeviceQuery := `DELETE FROM devices WHERE group_id = $1 AND device_name=$2 RETURNING id`
	err = tx.QueryRow(deleteDeviceQuery, g.ID, d.DeviceName).Scan(&d.ID)
	if err != nil {
		return fmt.Errorf("error while deleting form devices table:%v", err)
	}

	//commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error commiting transaction")
	}

	log.Printf("device with ID %v deleted", d.ID)

	return nil
}
