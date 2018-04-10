package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

//AllActiveFeeds -- Returns all active feeds
func AllActiveFeeds(db *sql.DB) map[int64]string {
	var result = make(map[int64]string)

	rows, err := db.Query("SELECT id, uri FROM feeds WHERE deleted = 0")
	if err != nil {
		log.Fatalf("Error happened when trying to get all active feeds: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var url string
		err := rows.Scan(&id, &url)
		if err != nil {
			log.Fatalf("Error happened while scanning the rows for the all active feeds function: %s", err.Error())
		}
		result[id] = url
	}
	return result
}

//FilterFeeds -- Takes in a list of feeds and compares them with the feeds listed in the Database.
//Returns all the feeds that are listed as active in the database but where not in the list.
func FilterFeeds(db *sql.DB, feeds map[int64]string) map[int64]string {
	var result = make(map[int64]string)
	allFeeds := AllActiveFeeds(db)

	for dbKey, dbValue := range allFeeds {
		found := false

		for feedKey, feedValue := range feeds {
			if dbKey == feedKey && strings.EqualFold(dbValue, feedValue) {
				found = true
				break
			}
		}

		if !found {
			result[dbKey] = dbValue
		}
	}

	return result
}

//DeleteFeed -- Flips the delete flag for a feed in the database
func DeleteFeed(db *sql.DB, feedID int64) error {
	_, err := db.Exec("UPDATE feeds SET deleted = 1 WHERE id = $1", feedID)
	return err
}

//FeedURLExist -- Checks to see if a feed exists
func FeedURLExist(db *sql.DB, url string) bool {
	var id int64
	var result bool

	row := db.QueryRow("SELECT id FROM feeds WHERE uri = $1", url)
	err := row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error happened when trying to check if the feed (%s) exists: %s", url, err.Error())
		}
	} else {
		result = true
	}
	return result
}

//AddFeedURL -- Adds a feed url to the database
func AddFeedURL(db *sql.DB, url string) (int64, error) {
	var result int64
	feedStmt := "INSERT INTO feeds (uri) VALUES ($1)"

	if FeedURLExist(db, url) {
		return result, fmt.Errorf("Feed already exists")
	}

	dbResult, err := db.Exec(feedStmt, url)
	if err != nil {
		log.Fatal(err)
	}

	result, err = dbResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}