package database

import (
	"database/sql"
	"fmt"
	"log"
)

//GetTagID -- Given a tag name, returns the tag id
func GetTagID(db *sql.DB, name string) (int64, error) {
	var id int64
	row := db.QueryRow("SELECT id FROM tags WHERE name = $1", name)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("Error occured while trying to find the tag id for tag name (%s): %s", name, err.Error())
	}
	return id, nil
}

//TagExist -- Checks to see if a tag exists
func TagExist(db *sql.DB, tag string) bool {
	var id int64
	var result bool

	row := db.QueryRow("SELECT id FROM tags WHERE name = $1", tag)
	err := row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error happened when trying to check if the tag (%s) exists: %s", tag, err.Error())
		}
	} else {
		result = true
	}

	return result
}

//AddTag -- Adds a tag to the database
func AddTag(db *sql.DB, tag string) (int64, error) {
	var result int64
	tagStmt := "INSERT INTO tags (name) VALUES ($1)"

	if TagExist(db, tag) {
		return result, fmt.Errorf("Tag (%s) already exists", tag)
	}
	dbResult, err := db.Exec(tagStmt, tag)
	if err != nil {
		log.Fatalf("Error adding tag to database: %s", err.Error())
	}

	result, err = dbResult.LastInsertId()
	if err != nil {
		log.Fatalf("Error happened when trying to get the last inserted id: %s", err.Error())
	}

	return result, nil
}
