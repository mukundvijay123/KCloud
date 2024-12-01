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
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", credentials["host"], credentials["port"], credentials["user"], credentials["password"], "KCloud")
	db2, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("error connecting to the database KCloud %v", err)
	}
	defer db2.Close()

	err = db2.Ping()
	if err != nil {
		log.Fatalf("Ping to KCloud DB failed : %v", err)
	}

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

// Making the database and required tables within postgres
func setupDB(db *sql.DB) error {

	//Checking if KCloud Database Exists
	dbName := "postgres"
	QueryForCheckingifDBExists := fmt.Sprintf(`SELECT 1 FROM pg_database WHERE datname= '%s';`, dbName)
	fmt.Println(QueryForCheckingifDBExists)
	row, err := db.Query(QueryForCheckingifDBExists)
	if err != nil {
		return fmt.Errorf("error while checking for existence of %s : %v", dbName, err)
	}
	defer row.Close()

	rowCount := 0
	for row.Next() {
		rowCount++
	}
	if rowCount == 0 {
		QueryForCreatingDB := fmt.Sprintf(`CREATE DATABASE %s`, dbName)
		_, err = db.Exec(QueryForCreatingDB)

		if err != nil {
			return fmt.Errorf("error creating database %s : %v", dbName, err)
		}
	}

	return nil

}

func createTables(db *sql.DB) error {

}
