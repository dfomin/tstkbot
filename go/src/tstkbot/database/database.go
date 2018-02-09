package database

import (
	"database/sql"
	"fmt"
	"log"
	"tstkbot/private"

	_ "github.com/lib/pq"
)

// DatabaseController represents controller for database
type DatabaseController struct {
	DataBase *sql.DB
}

// InitDatabase represents database initialization
func Init() *DatabaseController {
	databaseInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=verify-full", private.DatabaseUser, private.DatabasePassword, private.DatabaseName)

	database, err := sql.Open("postgres", databaseInfo)
	checkError("Open db failed", err)

	err = database.Ping()
	checkError("Ping db failed", err)

	return &DatabaseController{DataBase: database}
}

// JudgePhrases retrieves all accepted just phrases from database.
func (c *DatabaseController) JudgePhrases() []string {
	query := "SELECT phrase FROM " + private.DatabaseName + ".judge"
	rows, _ := c.DataBase.Query(query)
	defer rows.Close()
	var phrases []string
	for rows.Next() {
		var phrase string
		if err := rows.Scan(&phrase); err != nil {
			log.Fatal(err)
		}
		phrases = append(phrases, phrase)
	}

	return phrases
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}
