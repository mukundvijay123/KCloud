package metadata

type MetadataStore interface {
	CreateCompany(c *Company) error //Creates a company
	DeleteCompany(c *Company) error //Deletes a company
	UpdatePassword(c *Company, newPassword string) error
	CreateGroup(g *Grp) error                                    //Creates a group of sensors in a company
	DeleteGroup(g *Grp) error                                    //Deletes a group of sensors within a company
	CreateDevice(d *Device) error                                //Create a device entry
	DeleteDevice(d *Device) error                                //Deletes a device entry
	UpdateDeviceLocation(d *Device, l *Location) error           //Update Device Location
	UpdateDeviceSchema(d *Device, schema *TelemetrySchema) error //Updates Device Schema
}
