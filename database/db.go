package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3" //Sqlite3 driver
)

const (
	driver            = "sqlite3"
	foreignKeySupport = "?_foreign_keys=1"
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
