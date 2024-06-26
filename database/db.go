package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" //Sqlite3 driver

	"github.com/crazcalm/go-rss-reader/file"
)

func createTables(db *sql.DB) error {
	for _, sql := range sqlFiles {
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateDBDns(path string) string {
	return fmt.Sprintf("file:%s?_foreign_keys=1", path)
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

	return db, err
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
	var err error
	
	//Prep the connection to the database
	DB, err = sql.Open(driver, dsn)
	if err != nil {
		return DB, fmt.Errorf("Unable to Open Database: %w", err)
	}

	//Test connection to the database
	err = DB.Ping()
	if err != nil {
		return DB, err
	}

	if reset {
		//Drop all the tables and create all the tables again
		err = createTables(DB)
		if err != nil {
			return DB, err
		}
	}

	return DB, nil
}


//AddFeedFileData -- Adds Feed File Data to the database
func AddFeedFileData(db *sql.DB, fileData []file.Data) (map[int64]file.Data, error) {
	var feedID int64
	var tagID int64
	var err error
	var feedsFromFile = make(map[int64]string)
	var feedIDAndFileData = make(map[int64]file.Data)

	//This section ensures that all feeds are in the database
	for _, fd := range fileData {

		//Need to check if the feed is already in the database...
		if FeedURLExist(db, fd.URL) {
			feedID, err = GetFeedID(db, fd.URL)
			if err != nil {
				log.Fatal(err)
			}

			//Need to make sure that the feed is not deleted
			if IsFeedDeleted(db, feedID) {
				err = UndeleteFeed(db, feedID)
				if err != nil {
					log.Fatal(err)
				}
			}

		} else {
			//If it does not exist
			//Add feed url the database
			feedID, err = AddFeedURL(db, fd.URL)
			if err != nil {
				log.Fatal(err)
			}

			feedsFromFile[feedID] = fd.URL
		}

		//Used to keep track of all feeds
		feedIDAndFileData[feedID] = fd

		//This section ensures all the tags are in the database
		var tagsPerFeed = make(map[int64]string) // Container for all the tags associated with this feed
		var feedTagID int64
		for _, tag := range fd.Tags {
			//Need to check if the tag is already in the databse...
			if TagExist(db, tag) {
				tagID, err = GetTagID(db, tag)
				if err != nil {
					log.Fatal(err)
				}

			} else {
				//If the tag does not exist
				//Add tag to database
				tagID, err = AddTag(db, tag)
				if err != nil {
					log.Fatal(err)
				}
			}

			//Add tag to the list
			tagsPerFeed[tagID] = tag

			//Need to check if the feed and tag are in the feeds_and_tag database
			if FeedHasTag(db, feedID, tagID) {
				feedTagID, err = GetFeedTagID(db, feedID, tagID)
				if err != nil {
					log.Fatal(err)
				}
				if IsFeedTagDeleted(db, feedTagID) {
					err = UndeleteFeedTag(db, feedTagID)
					if err != nil {
						log.Fatal(err)
					}
				}

			} else {
				//Need to add feed tag to the database
				_, err = AddTagToFeed(db, feedID, tagID)
				if err != nil {
					log.Fatal(err)
				}
			}

			//This section ensures that the tags associated with the feeds are accurate.
			//We do this by comparing the tags we have from the file with the active
			//tags in the database. We will then delete all of the feed tags that are active
			//in the database that are not represented in the file.
			filteredTags := FilterFeedTags(db, feedID, tagsPerFeed)
			for dbTagID := range filteredTags {
				err = DeleteTagFromFeed(db, feedID, dbTagID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	//This section compares the feeds from the file and compares them
	//with the active feeds in the database. If a feed is marked as active
	//(deleted = 0), but is in the list of feeds from file, then that feed
	//is deleted.
	filteredFeeds := FilterFeeds(db, feedsFromFile)
	for fdID := range filteredFeeds {
		err := DeleteFeed(db, fdID)
		if err != nil {
			log.Fatal(err)
		}
	}
	return feedIDAndFileData, nil
}
