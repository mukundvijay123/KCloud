package metadatastore

import "github.com/mukundvijay123/KCloud/metadata"

// Forwarding methods so *MetadataDb satisfies metadata.MetadataStore

func (mdb *MetadataDb) GetCompanyByID(id string) (*metadata.Company, error) {
	return mdb.MetadataDbReader.GetCompanyByID(id)
}

func (mdb *MetadataDb) GetCompanyByUsername(username string) (*metadata.Company, error) {
	return mdb.MetadataDbReader.GetCompanyByUsername(username)
}

func (mdb *MetadataDb) ListCompanies() ([]*metadata.Company, error) {
	return mdb.MetadataDbReader.ListCompanies()
}

func (mdb *MetadataDb) VerifyCompany(username string, hashedPassword string) (bool, error) {
	return mdb.MetadataDbReader.VerifyCompany(username, hashedPassword)
}

func (mdb *MetadataDb) GetGroupByID(id string) (*metadata.Grp, error) {
	return mdb.MetadataDbReader.GetGroupByID(id)
}

func (mdb *MetadataDb) ListGroupsByCompany(companyID string) ([]*metadata.Grp, error) {
	return mdb.MetadataDbReader.ListGroupsByCompany(companyID)
}

func (mdb *MetadataDb) GetDeviceByID(id string) (*metadata.Device, error) {
	return mdb.MetadataDbReader.GetDeviceByID(id)
}

func (mdb *MetadataDb) ListDevicesByGroup(groupID string) ([]*metadata.Device, error) {
	return mdb.MetadataDbReader.ListDevicesByGroup(groupID)
}

func (mdb *MetadataDb) ListDevicesByCompany(companyID string) ([]*metadata.Device, error) {
	return mdb.MetadataDbReader.ListDevicesByCompany(companyID)
}
