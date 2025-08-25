package metadata

import "github.com/google/uuid"

type Company struct {
	ID              uuid.UUID `json:"id"` //UUID for company
	CompanyName     string    `json:"company_name"`
	Username        string    `json:"username"` //has to be unique
	CompanyPassword string    `json:"company_password"`
	NoOfGrps        int       `json:"no_of_grps"`
	NoOfDevices     int       `json:"no_of_devices"`
}

type Grp struct {
	ID          uuid.UUID ` json:"id"`
	CompanyID   uuid.UUID ` json:"company_id"`
	GroupName   string    ` json:"group_name"`
	NoOfDevices int       ` json:"no_of_devices"`
}

type Device struct {
	ID                  uuid.UUID       `json:"id"`
	GrpID               uuid.UUID       `json:"grp_id"`
	CompanyID           uuid.UUID       `json:"company_id"`
	DeviceName          string          `json:"device_name"`
	DeviceType          string          `json:"device_type"`
	DeviceDescription   string          `json:"device_description"`
	DeviceLocation      Location        `json:"device_location"`       //Can be null
	TelemetryDataSchema TelemetrySchema `json:"telemetry_data_schema"` //Non nested json schema with colname:type mapping
}
type TelemetrySchema map[string]string

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
