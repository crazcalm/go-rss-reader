package database

import (
	"database/sql"
	"log"
)

//AuthorExist -- Checks to see if an author exists
func AuthorExist(db *sql.DB, authorName, authorEmail string) (result bool) {
	var count int64
	row := db.QueryRow("SELECT COUNT(*) FROM authors WHERE name = $1 AND email = $2", authorName, authorEmail)
	err := row.Scan(&count)
	if err != nil {
		log.Fatalf("Error occurred while checking to see is an author exists (%s, %s): %s", authorName, authorEmail, err.Error())
	}

	if count != 0 {
		result = true
	}
	return
}

//GetAuthorByNameAndEmail -- Given a name and email, will return authorID
func GetAuthorByNameAndEmail(db *sql.DB, authorName, authorEmail string) (id int64, err error) {
	row := db.QueryRow("SELECT id FROM authors WHERE name = $1 AND email = $2", authorName, authorEmail)
	err = row.Scan(&id)
	if err != nil {
		return
	}
	return
}

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
