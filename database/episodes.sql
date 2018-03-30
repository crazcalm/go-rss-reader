DROP TABLE IF EXISTS episodes;

CREATE TABLE episodes (
	id				INTEGER PRIMARY KEY,
	uri				TEXT,
	title			TEXT,
	date			TIMESTAMP,
	media_content	TEXT,
	raw_data		TEXT,
	deleted			INTEGER CHECK (deleted = 0 OR deleted = 1),
	author_id		INTEGER,
	feed_id			INTEGER,
	FOREIGN KEY(author_id) REFERENCES authors(id) ON DELETE SET NULL,
	FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);
