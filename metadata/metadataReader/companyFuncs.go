package metadatareader

import (
	"database/sql"
	"fmt"

	types "github.com/mukundvijay123/KCloud/metadata"
)

// GetCompanyByID fetches a company by ID and nulls out the password
func (r *MetadataDBReader) GetCompanyByID(id string) (*types.Company, error) {
	row := r.dbConn.QueryRow(`
		SELECT id, company_name, username, no_of_grps, no_of_devices
		FROM company
		WHERE id=$1
	`, id)

	c := &types.Company{}
	err := row.Scan(&c.ID, &c.CompanyName, &c.Username, &c.NoOfGrps, &c.NoOfDevices)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Println("[GetCompanyByID] company not found:", id)
			return nil, nil
		}
		r.logger.Println("[GetCompanyByID] error querying company:", err)
		return nil, fmt.Errorf(ErrComanyNotFound.Error(), err)

	}

	// Null out password
	c.CompanyPassword = ""
	r.logger.Println("GetCompanyByID] found company")
	return c, nil
}

// GetCompanyByUsername fetches a company by username and nulls out the password
func (r *MetadataDBReader) GetCompanyByUsername(username string) (*types.Company, error) {
	row := r.dbConn.QueryRow(`
		SELECT id, company_name, username, no_of_grps, no_of_devices
		FROM company
		WHERE username=$1
	`, username)

	c := &types.Company{}
	err := row.Scan(&c.ID, &c.CompanyName, &c.Username, &c.NoOfGrps, &c.NoOfDevices)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Println("[GetCompanyByUsername] company not found:", username)
			return nil, nil
		}
		r.logger.Println("[GetCompanyByUsername] error querying company:", err)
		return nil, fmt.Errorf(ErrComanyNotFound.Error(), err)
	}

	// Null out password
	c.CompanyPassword = ""
	r.logger.Println("[GetCompanyByUsername] Company found")
	return c, nil
}

func (r *MetadataDBReader) ListCompanies() ([]*types.Company, error) {
	rows, err := r.dbConn.Query(`
		SELECT id, company_name, username, no_of_grps, no_of_devices
		FROM company
	`)
	if err != nil {
		r.logger.Println("[ListCompanies] error querying companies:", err)
		return nil, err
	}
	defer rows.Close()

	var companies []*types.Company
	for rows.Next() {
		c := &types.Company{}
		err := rows.Scan(&c.ID, &c.CompanyName, &c.Username, &c.NoOfGrps, &c.NoOfDevices)
		if err != nil {
			r.logger.Println("[ListCompanies] error scanning row:", err)
			continue
		}
		c.CompanyPassword = "" // null out password
		companies = append(companies, c)
	}

	if err = rows.Err(); err != nil {
		r.logger.Println("[ListCompanies] rows iteration error:", err)
		return nil, ErrComanyNotFound
	}

	return companies, nil
}

func (r *MetadataDBReader) VerifyCompany(username string, hashedPassword string) (bool, error) {
	var storedPassword string
	row := r.dbConn.QueryRow(`
		SELECT company_password
		FROM company
		WHERE username = $1
	`, username)

	err := row.Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Println("[VerifyCompany] username not found:", username)
			return false, nil
		}
		r.logger.Println("[VerifyCompany] error querying password:", err)
		return false, ErrDbErrorGeneric
	}

	// Compare the hashed password
	if storedPassword != hashedPassword {
		r.logger.Println("[VerifyCompany] password mismatch for username:", username)
		return false, nil
	}

	r.logger.Println("[VerifyCompany] credentials verified for username:", username)
	return true, nil
}
