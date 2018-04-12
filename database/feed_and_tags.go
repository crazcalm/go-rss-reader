package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

//UndeleteFeedTag -- Flips the delete flag off for a feed tag in the database
func UndeleteFeedTag(db *sql.DB, feedTagID int64) error {
	_, err := db.Exec("UPDATE feeds_and_tags SET deleted = 0 WHERE id = $1", feedTagID)
	return err
}

//GetFeedTagID -- gets the id for a feed tag
func GetFeedTagID(db *sql.DB, feedID, tagID int64) (int64, error) {
	var id int64
	row := db.QueryRow("SELECT id FROM feeds_and_tags WHERE feed_id = $1 AND tag_id = $2", feedID, tagID)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("Error happened while trying to get the feed tag id for feed id (%d) and tag id (%d): %s", feedID, tagID, err.Error())
	}
	return id, nil
}

//IsFeedTagDeleted -- Checks to see if the feed tag is currently marked as deleted
func IsFeedTagDeleted(db *sql.DB, feedTagID int64) bool {
	var result bool
	var deleted int64

	row := db.QueryRow("SELECT deleted FROM feeds_and_tags WHERE id = $1", feedTagID)
	err := row.Scan(&deleted)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Feed tag (%d) does not exist: %s", feedTagID, err.Error())
		} else {
			log.Fatalf("Error happened while trying check the value of the delete flag for feed tag (%d): %s", feedTagID, err.Error())
		}
	}

	if deleted == 1 {
		result = true
	} else {
		result = false
	}
	return result
}

//DeleteAllTagsFromFeed -- flips the delete flag for all tags associated with a feed
func DeleteAllTagsFromFeed(db *sql.DB, feedID int64) error {
	_, err := db.Exec("UPDATE feeds_and_tags SET deleted = 1 WHERE feed_id = $1", feedID)
	return err
}

//DeleteTagFromFeed -- Deletes a tag from a feed
func DeleteTagFromFeed(db *sql.DB, feedID, tagID int64) error {
	_, error := db.Exec("UPDATE feeds_and_tags SET deleted = 1 WHERE feed_id = $1 AND tag_id = $2", feedID, tagID)

	return error
}

//FeedHasTag -- Checks to see if a feed has a specific tag
func FeedHasTag(db *sql.DB, feedID, tagID int64) bool {
	var id int64
	var result bool
	query := "SELECT id FROM feeds_and_tags WHERE feed_id = $1 AND tag_id = $2 AND deleted != 1"

	row := db.QueryRow(query, feedID, tagID)
	err := row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error happened when trying to check if feed_id (%d) has tag_id (%d): %s", feedID, tagID, err.Error())
		}
	} else {
		result = true
	}
	return result
}

//AllActiveFeedTags -- Returns all of the active tags associated with a feed
func AllActiveFeedTags(db *sql.DB, feedID int64) map[int64]string {
	var result = make(map[int64]string)

	query := `SELECT tags.id, tags.name FROM feeds
	INNER JOIN feeds_and_tags ON feeds.id = feeds_and_tags.feed_id
	INNER JOIN tags ON tags.id = feeds_and_tags.tag_id
	WHERE feeds.id = $1`

	rows, err := db.Query(query, feedID)
	if err != nil {
		log.Fatalf("Error happened while trying to get all tags for feed id (%d): %s", feedID, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatalf("Error happened while scanning the rows for the all active feed tags function: %s", err.Error())
		}
		result[id] = name
	}
	return result
}

//FilterFeedTags -- Takes in a map of feed tags and compares them with the tags listed in the database for that feed.
//Returns all the tags for that feed that are listed as active in the database but where not passed in.
func FilterFeedTags(db *sql.DB, feedID int64, tags map[int64]string) map[int64]string {
	var result = make(map[int64]string)
	DBTags := AllActiveFeedTags(db, feedID)

	for dbKey, dbValue := range DBTags {
		found := false

		for feedKey, feedValue := range tags {
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

//AddTagToFeed -- Adds a Tag to a feed via the feeds_and_tags table
func AddTagToFeed(db *sql.DB, feedID, tagID int64) (int64, error) {
	var result int64
	stmt := "INSERT INTO feeds_and_tags (feed_id, tag_id) VALUES ($1, $2)"

	if FeedHasTag(db, feedID, tagID) {
		return result, fmt.Errorf("This feed already has that tag")
	}

	dbResult, err := db.Exec(stmt, feedID, tagID)
	if err != nil {
		log.Fatal(err)
	}

	result, err = dbResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
