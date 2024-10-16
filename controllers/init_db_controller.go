
package controllers

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("mysql", "root:NewPasswordHere@tcp(localhost:3306)/learning_platform")
    if err != nil {
        log.Fatal(err)
    }
    // Ping database to verify that it's working
    if err := DB.Ping(); err != nil {
        log.Fatal(err)
    }
    log.Println("Connected to the database")
}
