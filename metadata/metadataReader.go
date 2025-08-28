package metadata

//Readonly interface to get matadata

// MetadataReader defines read-only operations for metadata
type MetadataReader interface {
	// Companies
	GetCompanyByID(id string) (*Company, error)
	GetCompanyByUsername(username string) (*Company, error)
	ListCompanies() ([]*Company, error)
	VerifyCompany(username string, hashedPassword string) (bool, error)

	// Groups
	GetGroupByID(id string) (*Grp, error)
	ListGroupsByCompany(companyID string) ([]*Grp, error)

	// Devices
	GetDeviceByID(id string) (*Device, error)
	ListDevicesByGroup(groupID string) ([]*Device, error)
	ListDevicesByCompany(companyID string) ([]*Device, error)
}
