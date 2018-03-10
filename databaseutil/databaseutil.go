package databaseutil

import (
	"database/sql"
	// Notice that we’re loading the driver anonymously, The driver registers itself as being available to the database/sql package.
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var db *sql.DB

func Connect(dbUrl string) error {
	temp, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}

	db = temp

	// Quickly test if the connection to the database worked.
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

func SaveNewUser(
	displayName string,
	emailAddress string,
	password string) (int64, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	sqlStatement := `
		INSERT INTO users (display_name, email_address, password, creation_time) 
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err = db.QueryRow(sqlStatement, displayName, emailAddress, hash, time.Now().UTC()).Scan(&id)
	if err != nil {
		return -1, err
	}

	fmt.Sprintf("%d", id)
	return id, nil
}

func ValidateUser(
	emailAddress string,
	password string) (bool, error) {

	sqlStatement := `
	SELECT password FROM users WHERE email_address = $1
	`

	// TODO handle the scenario where there is nobody in the db
	var hashFromDatabase []byte
	err := db.QueryRow(sqlStatement, emailAddress).Scan(&hashFromDatabase)
	if err != nil {
		return false, err
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword(hashFromDatabase, []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}
