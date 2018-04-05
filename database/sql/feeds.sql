DROP TABLE IF EXISTS feeds;

CREATE TABLE feeds (
	id			INTEGER PRIMARY KEY,
	uri			TEXT NOT NULL UNIQUE,
	title		TEXT,
	raw_data	TEXT,
	deleted		INTEGER CHECK (deleted = 0 OR deleted = 1) DEFAULT 0,
	author_id	INTEGER,
	FOREIGN KEY (author_id) REFERENCES authors(id) ON DELETE SET NULL
);
