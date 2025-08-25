package metadatareader

import (
	"database/sql"
	"encoding/json"

	types "github.com/mukundvijay123/KCloud/metadata"
)

// GetDeviceByID fetches a device by ID
func (r *MetadataDBReader) GetDeviceByID(id string) (*types.Device, error) {
	row := r.dbConn.QueryRow(`
		SELECT id, grp_id, company_id, device_name, device_type, device_description, longitude, latitude, telemetry_data_schema
		FROM device
		WHERE id=$1
	`, id)

	d := &types.Device{}
	var schemaJSON []byte
	err := row.Scan(&d.ID, &d.GrpID, &d.CompanyID, &d.DeviceName, &d.DeviceType, &d.DeviceDescription,
		&d.DeviceLocation.Longitude, &d.DeviceLocation.Latitude, &schemaJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Println("[GetDeviceByID] device not found:", id)
			return nil, nil
		}
		r.logger.Println("[GetDeviceByID] query error:", err)
		return nil, err
	}

	// Unmarshal JSON schema
	var schema types.TelemetrySchema
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		r.logger.Println("[GetDeviceByID] failed to unmarshal schema:", err)
		return nil, err
	}
	d.TelemetryDataSchema = schema

	return d, nil
}

// ListDevicesByGroup lists all devices for a given group
func (r *MetadataDBReader) ListDevicesByGroup(groupID string) ([]*types.Device, error) {
	rows, err := r.dbConn.Query(`
		SELECT id, grp_id, company_id, device_name, device_type, device_description, longitude, latitude, telemetry_data_schema
		FROM device
		WHERE grp_id=$1
	`, groupID)
	if err != nil {
		r.logger.Println("[ListDevicesByGroup] query error:", err)
		return nil, err
	}
	defer rows.Close()

	var devices []*types.Device
	for rows.Next() {
		d := &types.Device{}
		var schemaJSON []byte
		if err := rows.Scan(&d.ID, &d.GrpID, &d.CompanyID, &d.DeviceName, &d.DeviceType,
			&d.DeviceDescription, &d.DeviceLocation.Longitude, &d.DeviceLocation.Latitude, &schemaJSON); err != nil {
			r.logger.Println("[ListDevicesByGroup] row scan error:", err)
			continue
		}

		var schema types.TelemetrySchema
		if err := json.Unmarshal(schemaJSON, &schema); err != nil {
			r.logger.Println("[ListDevicesByGroup] failed to unmarshal schema for device:", d.DeviceName, err)
			continue
		}
		d.TelemetryDataSchema = schema

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		r.logger.Println("[ListDevicesByGroup] rows iteration error:", err)
		return nil, err
	}

	return devices, nil
}

// ListDevicesByCompany lists all devices for a given company
func (r *MetadataDBReader) ListDevicesByCompany(companyID string) ([]*types.Device, error) {
	rows, err := r.dbConn.Query(`
		SELECT id, grp_id, company_id, device_name, device_type, device_description, longitude, latitude, telemetry_data_schema
		FROM device
		WHERE company_id=$1
	`, companyID)
	if err != nil {
		r.logger.Println("[ListDevicesByCompany] query error:", err)
		return nil, err
	}
	defer rows.Close()

	var devices []*types.Device
	for rows.Next() {
		d := &types.Device{}
		var schemaJSON []byte
		if err := rows.Scan(&d.ID, &d.GrpID, &d.CompanyID, &d.DeviceName, &d.DeviceType,
			&d.DeviceDescription, &d.DeviceLocation.Longitude, &d.DeviceLocation.Latitude, &schemaJSON); err != nil {
			r.logger.Println("[ListDevicesByCompany] row scan error:", err)
			continue
		}

		var schema types.TelemetrySchema
		if err := json.Unmarshal(schemaJSON, &schema); err != nil {
			r.logger.Println("[ListDevicesByCompany] failed to unmarshal schema for device:", d.DeviceName, err)
			continue
		}
		d.TelemetryDataSchema = schema

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		r.logger.Println("[ListDevicesByCompany] rows iteration error:", err)
		return nil, err
	}

	return devices, nil
}
