package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" //Sqlite3 driver

	"github.com/crazcalm/go-rss-reader"
)

const (
	driver            = "sqlite3"
	foreignKeySupport = "?_foreign_keys=1"
	//TestDB -- testing database
	TestDB = "file:test.db?_foreign_keys=1"
)

var (
	sqlFiles = [...]string{"sql/authors.sql", "sql/tags.sql", "sql/feeds.sql", "sql/episodes.sql", "sql/feeds_and_tags.sql"}
)

func createTables(db *sql.DB) error {
	for _, path := range sqlFiles {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(data))
		if err != nil {
			return err
		}
	}
	return nil
}

//Create -- Created the database
func Create(path string) (*sql.DB, error) {
	//Create the database file
	_, err := os.Create(path)
	if err != nil {
		return &sql.DB{}, err
	}

	//Initialize database and create tables
	db, err := Init(fmt.Sprintf("file:%s%s", path, foreignKeySupport), true)
	if err != nil {
		return db, nil
	}
	return db, nil
}

//Exist -- checks for the existance of a file
func Exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//Init -- Initializes the database. The reset param allows you to recreate the database.
func Init(dsn string, reset bool) (*sql.DB, error) {
	//Prep the connection to the database
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return db, err
	}

	//Test connection to the database
	err = db.Ping()
	if err != nil {
		return db, err
	}

	if reset {
		//Drop all the tables and create all the tables again
		err = createTables(db)
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

//AddFeedFileData -- Adds Feed File Data to the database
func AddFeedFileData(fileData []rss.FileData) error {
	var feedID int64
	var tagID int64

	db, err := Init(TestDB, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, fd := range fileData {

		//Need to check if the feed is already in the database...
		if !FeedURLExist(db, fd.URL) {
			//Add feed url the database
			feedID, err = AddFeedURL(db, fd.URL)
			if err != nil {
				log.Fatal(err)
			}
		} // If it does exist, I need to get the feed id

		for _, tag := range fd.Tags {
			//Need to check if the tag is already in the databse...
			if !TagExist(db, tag) {
				//Add tag to database
				tagID, err = AddTag(db, tag)
				if err != nil {
					log.Fatal(err)
				}
			} //  If it does exist, I need to get the tag id

			//Need to check if the feed and tag are in the feeds_and_tag database
			if !FeedHasTag(db, feedID, tagID) {
				_, err := AddTagToFeed(db, feedID, tagID)
				if err != nil {
					log.Fatal(err)
				}
			}
			//Need to add new tags to the database
			//Need to delete old tags that are no longer associated with the feed
			// TODO: I also need to do the same thing for feeds...

		}
	}

	return nil
}
