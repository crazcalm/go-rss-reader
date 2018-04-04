DROP TABLE IF EXISTS feeds_and_tags;

CREATE TABLE feeds_and_tags (
	feed_id		INTEGER,
	tag_id		INTEGER,
	FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
	FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
