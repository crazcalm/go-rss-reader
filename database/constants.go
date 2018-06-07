package database

import (
	"database/sql"
	"fmt"
	"path/filepath"
)

const (
	driver            = "sqlite3"
	foreignKeySupport = "?_foreign_keys=1"

	//AuthorsTable -- author table sql
	AuthorsTable = `
	DROP TABLE IF EXISTS authors;
	
	CREATE TABLE authors (
		id		INTEGER PRIMARY KEY,
		name    TEXT,
		email	TEXT
	);
	`

	//EpisodesTable -- episodes table sql
	EpisodesTable = `
	DROP TABLE IF EXISTS episodes;
	
	CREATE TABLE episodes (
		id				INTEGER PRIMARY KEY,
		uri				TEXT,
		title			TEXT,
		date			TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		media_content	TEXT,
		raw_data		TEXT,
		seen			INTEGER CHECK (seen = 0 OR seen = 1) DEFAULT 0,
		deleted			INTEGER CHECK (deleted = 0 OR deleted = 1) DEFAULT 0,
		author_id		INTEGER,
		feed_id			INTEGER NOT NULL,
		FOREIGN KEY(author_id) REFERENCES authors(id) ON DELETE SET NULL,
		FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
	);
	`

	//FeedsTable -- feeds table sql
	FeedsTable = `
	DROP TABLE IF EXISTS feeds;
	
	CREATE TABLE feeds (
		id			INTEGER PRIMARY KEY,
		uri			TEXT NOT NULL UNIQUE,
		title		TEXT DEFAULT "",
		raw_data	TEXT DEFAULT "",
		deleted		INTEGER CHECK (deleted = 0 OR deleted = 1) DEFAULT 0,
		author_id	INTEGER,
		FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE SET NULL
	);
	`

	//TagsTable -- tags table sql
	TagsTable = `
	DROP TABLE IF EXISTS tags;
	
	CREATE TABLE tags (
		id		INTEGER PRIMARY KEY,
		name	TEXT NOT NULL UNIQUE
	);
	`

	//FeedsAndTagsTable -- feeds_and_tags table sql
	FeedsAndTagsTable = `
	DROP TABLE IF EXISTS feeds_and_tags;
	
	CREATE TABLE feeds_and_tags (
		id			INTEGER PRIMARY KEY,
		feed_id		INTEGER,
		tag_id		INTEGER,
		deleted		INTEGER CHECK (deleted = 0 OR deleted = 1) DEFAULT 0,
		FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	`
)

var (
	sqlFiles = [...]string{AuthorsTable, TagsTable, FeedsTable, EpisodesTable, FeedsAndTagsTable}
	//TestDB -- testing database
	TestDB = fmt.Sprintf("file:%s?_foreign_keys=1", DBPath)
	//DB -- A global reference to the DB
	DB *sql.DB
	//DBPath -- path to test database
	DBPath = filepath.Join(".go-rss-reader", "feeds.db")
)
