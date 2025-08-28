package metadatastore

import (
	"encoding/json"
	"fmt"

	types "github.com/mukundvijay123/KCloud/metadata"
)

func (mdb *MetadataDb) CreateDevice(d *types.Device) error {
	if !isValidName(d.DeviceName) {
		mdb.logger.Println("[CreateDevice] invalid device name:", d.DeviceName)
		return ErrInvalidName
	}

	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("[CreateDevice] failed to begin transaction:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			mdb.logger.Println("[CreateDevice] transaction rolled back due to error:", err)
		}
	}()

	insertDeviceQuery := `
		INSERT INTO device (grp_id, company_id, device_name, device_type, device_description, longitude, latitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err = tx.QueryRow(
		insertDeviceQuery,
		d.GrpID,
		d.CompanyID,
		d.DeviceName,
		d.DeviceType,
		d.DeviceDescription,
		d.DeviceLocation.Longitude,
		d.DeviceLocation.Latitude,
	).Scan(&d.ID)
	if err != nil {
		mdb.logger.Println("[CreateDevice] failed to insert device:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}
	mdb.logger.Println("[CreateDevice] device inserted with ID:", d.ID)

	_, err = tx.Exec(`UPDATE grp SET no_of_devices = no_of_devices + 1 WHERE id=$1`, d.GrpID)
	if err != nil {
		mdb.logger.Println("[CreateDevice] failed to update grp count:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	_, err = tx.Exec(`UPDATE company SET no_of_devices = no_of_devices + 1 WHERE id=$1`, d.CompanyID)
	if err != nil {
		mdb.logger.Println("[CreateDevice] failed to update company count:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	if err = tx.Commit(); err != nil {
		mdb.logger.Println("[CreateDevice] failed to commit transaction:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	mdb.logger.Println("[CreateDevice] device created successfully:", d.DeviceName)
	return nil
}

func (mdb *MetadataDb) DeleteDevice(d *types.Device) error {
	tx, err := mdb.dbConn.Begin()
	if err != nil {
		mdb.logger.Println("[DeleteDevice] failed to begin transaction:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			mdb.logger.Println("[DeleteDevice] transaction rolled back due to error:", err)
		}
	}()

	deleteDeviceQuery := `DELETE FROM device WHERE id=$1`
	res, err := tx.Exec(deleteDeviceQuery, d.ID)
	if err != nil {
		mdb.logger.Println("[DeleteDevice] failed to delete device:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		mdb.logger.Println("[DeleteDevice] device not found with ID:", d.ID)
		return fmt.Errorf(ErrDeviceNotExist.Error(), d.ID)
	}
	mdb.logger.Println("[DeleteDevice] device deleted with ID:", d.ID)

	_, err = tx.Exec(`UPDATE grp SET no_of_devices = no_of_devices - 1 WHERE id=$1`, d.GrpID)
	if err != nil {
		mdb.logger.Println("[DeleteDevice] failed to decrement grp count:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	_, err = tx.Exec(`UPDATE company SET no_of_devices = no_of_devices - 1 WHERE id=$1`, d.CompanyID)
	if err != nil {
		mdb.logger.Println("[DeleteDevice] failed to decrement company count:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	if err = tx.Commit(); err != nil {
		mdb.logger.Println("[DeleteDevice] failed to commit transaction:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	mdb.logger.Println("[DeleteDevice] device deleted successfully:", d.DeviceName)
	return nil
}

func (mdb *MetadataDb) UpdateDeviceLocation(d *types.Device, l *types.Location) error {
	query := `UPDATE device SET longitude=$1, latitude=$2 WHERE id=$3`
	res, err := mdb.dbConn.Exec(query, l.Longitude, l.Latitude, d.ID)
	if err != nil {
		mdb.logger.Println("[UpdateDeviceLocation] failed to update location:", err)
		return fmt.Errorf(ErrDbErrorGeneric.Error(), err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		mdb.logger.Println("[UpdateDeviceLocation] device not found with ID:", d.ID)
		return fmt.Errorf(ErrDeviceNotExist.Error(), d.ID)
	}

	mdb.logger.Println("[UpdateDeviceLocation] device location updated for:", d.DeviceName)
	return nil
}

// UpdateDeviceSchema updates the telemetry schema JSON field of a device
func (mdb *MetadataDb) UpdateDeviceSchema(d *types.Device, schema *types.TelemetrySchema) (err error) {
	// Validate schema
	validTypes := map[string]bool{
		"int":    true,
		"float":  true,
		"bool":   true,
		"string": true,
	}

	for field, typ := range *schema {
		if len(field) > 32 {
			mdb.logger.Println("[UpdateDeviceSchema] invalid field name length:", field)
			return fmt.Errorf("field name '%s' exceeds 32 characters", field)
		}
		if !validTypes[typ] {
			mdb.logger.Println("[UpdateDeviceSchema] invalid type for field:", field, "type:", typ)
			return fmt.Errorf("invalid type '%s' for field '%s'", typ, field)
		}
	}

	// Marshal schema into JSON
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		mdb.logger.Println("[UpdateDeviceSchema] failed to marshal schema:", err)
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Update DB
	query := `UPDATE device SET telemetry_data_schema = $1 WHERE id = $2`
	res, err := mdb.dbConn.Exec(query, schemaJSON, d.ID)
	if err != nil {
		mdb.logger.Println("[UpdateDeviceSchema] failed to update schema in DB:", err)
		return fmt.Errorf("db error: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		mdb.logger.Println("[UpdateDeviceSchema] device not found with ID:", d.ID)
		return fmt.Errorf("device with id %s does not exist", d.ID)
	}

	// Update local struct copy
	d.TelemetryDataSchema = *schema

	mdb.logger.Println("[UpdateDeviceSchema] schema updated successfully for device:", d.DeviceName)
	return nil
}
