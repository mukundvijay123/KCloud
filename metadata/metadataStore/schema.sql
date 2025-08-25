-- COMPANY TABLE
CREATE TABLE company (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    company_password VARCHAR(255) NOT NULL,
    no_of_grps INT DEFAULT 0,
    no_of_devices INT DEFAULT 0,
    UNIQUE (company_name),
    UNIQUE (username)
);

-- GROUP TABLE
CREATE TABLE grp (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    grp_name VARCHAR(255) NOT NULL,
    no_of_devices INT DEFAULT 0,
    CONSTRAINT unique_grp_per_company UNIQUE (company_id, grp_name)
);

-- DEVICE TABLE
CREATE TABLE device (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grp_id UUID NOT NULL REFERENCES grp(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES company(id) ON DELETE CASCADE,
    device_name VARCHAR(255) NOT NULL,
    device_type VARCHAR(255),
    device_description TEXT,
    longitude DOUBLE PRECISION,
    latitude DOUBLE PRECISION,
    telemetry_data_schema JSONB DEFAULT '{}'::jsonb NOT NULL,
    CONSTRAINT unique_device_per_grp UNIQUE (grp_id, device_name)
);
