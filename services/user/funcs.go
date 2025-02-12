package user

import (
	"database/sql"
	"fmt"

	"github.com/mukundvijay123/KCloud/types"
)

// Function for endpoint /user/login for users using Web UI to login
func Login(c *types.Company, db *sql.DB) (bool, error) {

	//Query for Login
	LoginQuery := `SELECT id FROM company WHERE username=$1 AND company_password=$2`

	var userID int
	err := db.QueryRow(LoginQuery, c.Username, c.CompanyPassword).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("incorrect credentials")
		}

		return false, fmt.Errorf("failed to login")
	}

	return true, nil
}
