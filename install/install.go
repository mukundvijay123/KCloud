package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

func main() {
	credentialsFilename := "credentials.yaml"
	credentials, err := getCredentials(credentialsFilename)
	if err != nil {
		log.Fatal("Couldnt fetch credentials ", err)
	}
	//fmt.Println(credentials)

	//Connecting t the database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", credentials["host"], credentials["port"], credentials["user"], credentials["password"], credentials["dbname"])
	db1, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to postgres: ", err)
	}
	defer db1.Close()

	//Pinging database
	err = db1.Ping()
	if err != nil {
		log.Fatal("Ping to postgres failed: ", err)
	}

	//Connected to datbase
	fmt.Println("Succesfully connected to postgres")
	//Setup is done
	err = setupDB(db1)
	if err != nil {
		log.Fatal(err)
	}

	//Connecting to new database and making tables
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", credentials["host"], credentials["port"], credentials["user"], credentials["password"], "kcloud")
	db2, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error connecting to the database KCloud %v", err)
	}
	defer db2.Close()

	err = db2.Ping()
	if err != nil {
		log.Fatalf("Ping to KCloud DB failed : %v", err)
	}
	//connected to kcloud db
	fmt.Println("Connected to kcloud DB")
	err = createTables(db2)
	if err != nil {
		log.Fatalf("Creating tables failed :%v", err)
	}
	fmt.Println("Tables created installation successful")

}

// This function fetches credentals from the yaml file
func getCredentials(filename string) (map[string]interface{}, error) {

	//Creating map to return
	var result map[string]interface{}

	//Reads file and returns slice of bytes
	data, err := os.ReadFile(filename)
	if err != nil {
		return result, errors.Join(errors.New("error reading file for credentials"), err)
	}

	//Unmarshaling yaml
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return result, errors.Join(errors.New("error Unmarshalling crediantials yaml file"), err)
	}

	//fmt.Println(result)
	credentials := result["credentials"].(map[string]interface{}) //type assertion

	return credentials, nil

}

// setupDB creates the two databases (kcloud and kcloud_data).
func setupDB(db *sql.DB) error {

	// Defining the database names
	DbName := "kcloud"

	// Create the metadata database (kcloud)
	_, err := db.Exec(fmt.Sprintf(`CREATE DATABASE %s`, DbName))
	if err != nil {
		return fmt.Errorf("error creating database %s: %v", DbName, err)
	}

	return nil
}

func createTables(db *sql.DB) error {
	// Lambda function to generate the query for checking if a table exists
	generateQuery := func(tableName string) string {
		return fmt.Sprintf(
			"SELECT EXISTS (\n  SELECT 1\n  FROM   information_schema.tables\n  WHERE  table_schema = 'public'\n  AND    table_name = '%s'\n);",
			tableName,
		)
	}

	// Function to check if the table exists and create it if necessary
	checkAndCreateTable := func(tableName, createQuery string) error {
		// Generate the query to check if the table exists
		query := generateQuery(tableName)

		// Execute the query to check if the table exists
		var exists bool
		err := db.QueryRow(query).Scan(&exists)
		if err != nil {
			return fmt.Errorf("error checking if table %s exists: %v", tableName, err)
		}

		// If the table exists, return an error
		if exists {
			return fmt.Errorf("table %s already exists", tableName)
		}

		// If the table doesn't exist, execute the create table query
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("error creating table %s: %v", tableName, err)
		}

		log.Printf("Table %s created successfully", tableName)
		return nil
	}

	// Queries to create tables
	companyCreateQuery := `CREATE TABLE "company"(
		id SERIAL PRIMARY KEY,
		company_name VARCHAR(32) NOT NULL,
		username VARCHAR(32) NOT NULL UNIQUE,
		company_password VARCHAR(32) NOT NULL,
		no_of_grps INT NOT NULL,
		no_of_devices INT NOT NULL
	);`
	groupCreateQuery := `CREATE TABLE grp (
		id SERIAL PRIMARY KEY,
		group_name VARCHAR(32) ,
		no_of_devices INT NOT NULL,
		company_id INT,
		FOREIGN KEY (company_id) REFERENCES company(id)
			ON DELETE CASCADE 
	);`
	devicesCreateQuery := `CREATE TABLE devices (
		id SERIAL PRIMARY KEY,
		device_name VARCHAR(32) NOT NULL ,
		group_id INT NOT NULL,
		company_id INT NOT NULL,
		telemetry_data_schema JSONB NOT NULL,
		device_description VARCHAR(100),
		device_type VARCHAR(32) NOT NULL,
		longitude DOUBLE PRECISION,
		latitude DOUBLE PRECISION,
		FOREIGN KEY (group_id) REFERENCES grp(id) 
			ON DELETE CASCADE,
		FOREIGN KEY (company_id) REFERENCES company(id)
			ON DELETE CASCADE
		
		
	);`

	telemetryCreateQuery := `CREATE TABLE data(
		company_id INT NOT NULL,
		device_id INT NOT NULL,
		timestamp TIMESTAMPTZ NOT NULL,
		telemetry_data JSONB NOT NULL,
		PRIMARY KEY(company_id,device_id,timestamp)
	) PARTITION BY LIST (company_id);`

	// Check and create each table
	if err := checkAndCreateTable("company", companyCreateQuery); err != nil {
		return err
	}

	if err := checkAndCreateTable("group", groupCreateQuery); err != nil {
		return err
	}

	if err := checkAndCreateTable("devices", devicesCreateQuery); err != nil {
		return err
	}

	if err := checkAndCreateTable("data", telemetryCreateQuery); err != nil {
		return err
	}

	return nil
}
