DROP TABLE IF EXISTS feeds_and_tags;

CREATE TABLE feeds_and_tags (
	id			INTEGER PRIMARY KEY,
	feed_id		INTEGER,
	tag_id		INTEGER,
	deleted		INTEGER CHECK (deleted = 0 OR deleted = 1) DEFAULT 0,
	FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
