DROP TABLE IF EXISTS tags;

CREATE TABLE tags (
	id		INTEGER PRIMARY KEY,
	name	TEXT NOT NULL UNIQUE
);
