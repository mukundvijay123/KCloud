package metadatastore

import (
	"database/sql"
	"fmt"

	types "github.com/mukundvijay123/KCloud/metadata"
)

// adding a  company to metadata store
func (mdb *MetadataDb) CreateCompany(c *types.Company) error {

	if !isValidName(c.CompanyName) || !isValidName(c.Username) {
		mdb.logger.Println("invalid company or username: ", c.CompanyName)
		return ErrInvalidName
	}

	if !isValidPasswd(c.CompanyPassword) {
		mdb.logger.Println("invalid password")
		return ErrInvalidPasswd
	}

	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("Error creating a transaction: ", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertCompanyQuery := `INSERT INTO company (company_name, username, company_password, no_of_grps, no_of_devices) 
                    VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = tx.QueryRow(insertCompanyQuery, c.CompanyName, c.Username, c.CompanyPassword, c.NoOfGrps, c.NoOfDevices).Scan(&c.ID)
	if err != nil {
		mdb.logger.Println("Error creating company: ", ErrDbErrorGeneric.Error(), err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	//If required communicate to storage engine here
	//Functionality to be added later
	if err = tx.Commit(); err != nil {
		mdb.logger.Println(ErrDbErrorGeneric)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}
	mdb.logger.Println("Company provisioned successfully")

	return nil
}

// Delete a  company entry from metadata store
func (mdb *MetadataDb) DeleteCompany(c *types.Company) error {
	if !isValidName(c.Username) {
		mdb.logger.Println("invalid company or username: ", c.CompanyName)
		return ErrInvalidName
	}

	if !isValidPasswd(c.CompanyPassword) {
		mdb.logger.Println("invalid password")
		return ErrInvalidPasswd
	}
	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("Error creating a transaction: ", err)
		return fmt.Errorf("%w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var dbPassword string
	queryCheck := `SELECT company_password FROM company WHERE username = $1`
	err = tx.QueryRow(queryCheck, c.Username).Scan(&dbPassword)

	if err == sql.ErrNoRows {
		mdb.logger.Println(ErrCompanyNoExist, c.Username)
		return ErrCompanyNoExist
	}
	if err != nil {
		mdb.logger.Println(ErrDbErrorGeneric, err)
		return ErrDbErrorGeneric
	}

	// simple password match (better: hash+compare, not plain text)
	if dbPassword != c.CompanyPassword {
		mdb.logger.Println("Invalid password for company:", c.Username)
		return ErrInvalidPasswd
	}

	deleteQuery := `DELETE FROM company WHERE username = $1`
	_, err = tx.Exec(deleteQuery, c.Username)
	if err != nil {
		mdb.logger.Println(ErrDbErrorGeneric, err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		mdb.logger.Println("Error committing delete transaction:", commitErr)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), commitErr)
	}

	mdb.logger.Println("Company deleted successfully:", c.Username)
	return nil
}

func (mdb *MetadataDb) UpdatePassword(c *types.Company, newPassword string) error {
	// Validate inputs
	if !isValidName(c.Username) {
		mdb.logger.Println("invalid username:", c.Username)
		return ErrInvalidName
	}
	if !isValidPasswd(c.CompanyPassword) {
		mdb.logger.Println("invalid current password")
		return ErrInvalidPasswd
	}
	if !isValidPasswd(newPassword) {
		mdb.logger.Println("invalid new password")
		return ErrInvalidPasswd
	}

	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("Error creating transaction:", err)
		return fmt.Errorf("%w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Verify old password
	var dbPassword string
	queryCheck := `SELECT company_password FROM company WHERE username = $1`
	err = tx.QueryRow(queryCheck, c.Username).Scan(&dbPassword)
	if err == sql.ErrNoRows {
		mdb.logger.Println("Company not found for username:", c.Username)
		return fmt.Errorf("company not found")
	}
	if err != nil {
		mdb.logger.Println("Error checking company:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	if dbPassword != c.CompanyPassword {
		mdb.logger.Println("Incorrect current password for username:", c.Username)
		return ErrInvalidPasswd
	}

	// Update password
	updateQuery := `UPDATE company SET company_password = $1 WHERE username = $2`
	_, err = tx.Exec(updateQuery, newPassword, c.Username)
	if err != nil {
		mdb.logger.Println("Error updating password:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	// Commit
	if err = tx.Commit(); err != nil {
		mdb.logger.Println("Error committing password update:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	// Update in-memory struct
	c.CompanyPassword = newPassword
	mdb.logger.Println("Password updated successfully for username:", c.Username)

	return nil
}
