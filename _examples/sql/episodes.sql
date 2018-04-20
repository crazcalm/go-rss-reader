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
