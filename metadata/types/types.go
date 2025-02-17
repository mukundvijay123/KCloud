package types

import "encoding/json"

const MaxMetadataRequestSize = 1024 * 10
const MaxUserRequestSize = 1024 * 10

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
	DeviceType          string          `json:"device_type"`
	DeviceDescription   string          `json:"device_description"`
	Longitude           float64         `json:"longitude"`
	Latitude            float64         `json:"latitude"`
	TelemetryDataSchema json.RawMessage `json:"telemetry_data_schema"`
}
