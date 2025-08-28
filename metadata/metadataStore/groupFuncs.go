package metadatastore

import (
	"fmt"

	types "github.com/mukundvijay123/KCloud/metadata"
)

func (mdb *MetadataDb) CreateGroup(g *types.Grp) error {
	if !isValidName(g.GroupName) {
		mdb.logger.Println("invalid group name:", g.GroupName)
		return ErrInvalidName
	}

	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("Error creating transaction:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert the group
	insertGroupQuery := `
		INSERT INTO grp (company_id, group_name, no_of_devices)
		VALUES ($1, $2, $3) RETURNING id
	`
	err = tx.QueryRow(insertGroupQuery, g.CompanyID, g.GroupName, g.NoOfDevices).Scan(&g.ID)
	if err != nil {
		mdb.logger.Println("Error inserting group:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	// Increment company's group count
	updateCompanyQuery := `
		UPDATE company SET no_of_grps = no_of_grps + 1 WHERE id = $1
	`
	_, err = tx.Exec(updateCompanyQuery, g.CompanyID)
	if err != nil {
		mdb.logger.Println("Error updating company group count:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	if err = tx.Commit(); err != nil {
		mdb.logger.Println("Error committing group creation:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	mdb.logger.Println("Group created successfully:", g.GroupName)
	return nil
}

func (mdb *MetadataDb) DeleteGroup(g *types.Grp) error {
	if !isValidName(g.GroupName) {
		mdb.logger.Println("invalid group name:", g.GroupName)
		return ErrInvalidName
	}

	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("Error creating transaction:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Delete group by ID + CompanyID
	deleteQuery := `DELETE FROM grp WHERE id = $1 AND company_id = $2`
	res, err := tx.Exec(deleteQuery, g.ID, g.CompanyID)
	if err != nil {
		mdb.logger.Println("Error deleting group:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		mdb.logger.Println("No group deleted, not found:", g.ID)
		return fmt.Errorf("group not found")
	}

	// Decrement company's group count
	updateCompanyQuery := `
		UPDATE company SET no_of_grps = no_of_grps - 1 WHERE id = $1
	`
	_, err = tx.Exec(updateCompanyQuery, g.CompanyID)
	if err != nil {
		mdb.logger.Println("Error updating company group count:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	if err = tx.Commit(); err != nil {
		mdb.logger.Println("Error committing group deletion:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	mdb.logger.Println("Group deleted successfully:", g.GroupName)
	return nil
}
