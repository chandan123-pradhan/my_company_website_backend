package controllers

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection pool
var DB *sql.DB

// InitDB initializes the database connection pool.
// 
// This function establishes a connection to a MySQL database using the
// provided DSN (Data Source Name) and checks if the connection is alive
// by sending a ping to the database. If the connection fails or if the
// database is unreachable, the function will log the error and terminate
// the application.
//
// This function should be called during the application startup to ensure
// the database is ready for use.
func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:NewPasswordHere@tcp(localhost:3306)/learning_platform")
	if err != nil {
		log.Fatal(err) // Log and terminate if there is an error opening the database
	}
	// Ping the database to verify that the connection is working
	if err := DB.Ping(); err != nil {
		log.Fatal(err) // Log and terminate if the ping fails
	}
}
