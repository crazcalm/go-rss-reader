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
