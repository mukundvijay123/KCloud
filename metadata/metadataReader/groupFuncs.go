package metadatareader

import (
	"database/sql"

	types "github.com/mukundvijay123/KCloud/metadata"
)

func (r *MetadataDBReader) GetGroupByID(id string) (*types.Grp, error) {
	row := r.dbConn.QueryRow(`
		SELECT id, company_id, grp_name, no_of_devices
		FROM grp
		WHERE id=$1
	`, id)

	g := &types.Grp{}
	err := row.Scan(&g.ID, &g.CompanyID, &g.GroupName, &g.NoOfDevices)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Println("[GetGroupByID] group not found:", id)
			return nil, nil
		}
		r.logger.Println("[GetGroupByID] query error:", err)
		return nil, err
	}
	return g, nil
}

// ListGroupsByCompany lists all groups for a given company
func (r *MetadataDBReader) ListGroupsByCompany(companyID string) ([]*types.Grp, error) {
	rows, err := r.dbConn.Query(`
		SELECT id, company_id, grp_name, no_of_devices
		FROM grp
		WHERE company_id=$1
	`, companyID)
	if err != nil {
		r.logger.Println("[ListGroupsByCompany] query error:", err)
		return nil, err
	}
	defer rows.Close()

	var groups []*types.Grp
	for rows.Next() {
		g := &types.Grp{}
		if err := rows.Scan(&g.ID, &g.CompanyID, &g.GroupName, &g.NoOfDevices); err != nil {
			r.logger.Println("[ListGroupsByCompany] row scan error:", err)
			continue
		}
		groups = append(groups, g)
	}

	if err := rows.Err(); err != nil {
		r.logger.Println("[ListGroupsByCompany] rows iteration error:", err)
		return nil, err
	}

	return groups, nil
}
