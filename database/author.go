package database

import (
	"database/sql"
	"log"
)

//AddAuthor -- adds a author to the database
func AddAuthor(db *sql.DB, name, email string) (int64, error) {
	stmt := "INSERT INTO authors (name, email) VALUES ($1, $2)"

	dbResult, err := db.Exec(stmt, name, email)
	if err != nil {
		log.Fatal(err)
	}

	result, err := dbResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

//GetAuthor -- returns the name and email of the author
func GetAuthor(db *sql.DB, authorID int64) (name string, email string, err error) {
	row := db.QueryRow("SELECT name, email FROM authors WHERE id = $1", authorID)
	err = row.Scan(&name, &email)
	if err != nil {
		return
	}
	return
}
