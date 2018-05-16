package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

//GetEpisodeHeaderData -- Gets the infomation needed to create the episode header
func GetEpisodeHeaderData(db *sql.DB, feedID, episodeID int64) (feedTitle, episodeTitle, author, episodeLink, dateString, mediaContent string, err error) {
	//Getting the feed Title
	feedTitle, err = GetFeedTitle(db, feedID)
	if err != nil {
		return
	}

	//Getting the episodeLink, episodeTitle and dateString
	var date time.Time
	episodeLink, episodeTitle, date, _, _, err = GetEpisode(db, episodeID)
	if err != nil {
		return
	}
	dateString = date.Format(time.RFC1123)

	//Getting Media Content, if exists
	if EpisodeHasMediaContent(db, episodeID) {
		mediaContent, err = GetEpisodeMediaContent(db, episodeID)
		if err != nil {
			return
		}
	}

	//Getting an author, if exists
	var authorName string
	var authorEmail string
	if EpisodeHasAuthor(db, episodeID) {
		authorName, authorEmail, err = GetEpisodeAuthor(db, episodeID)
		if err != nil {
			return
		}

		if !strings.EqualFold(authorName, "") && !strings.EqualFold(authorEmail, "") {
			author = fmt.Sprintf("%s (%s)", authorName, authorEmail)
		} else {
			if !strings.EqualFold(authorName, "") {
				author = authorName
			} else if !strings.EqualFold(authorEmail, "") {
				author = authorEmail
			}
		}

	} else if FeedHasAuthor(db, feedID) {
		authorName, authorEmail, err = GetFeedAuthor(db, feedID)
		if err != nil {
			return
		}

		if !strings.EqualFold(authorName, "") && !strings.EqualFold(authorEmail, "") {
			author = fmt.Sprintf("%s (%s)", authorName, authorEmail)
		} else {
			if !strings.EqualFold(authorName, "") {
				author = authorName
			} else if !strings.EqualFold(authorEmail, "") {
				author = authorEmail
			}
		}
	}

	return
}

//GetEpisodeIDByFeedIDAndTitle -- Gets the episodes using feed id and episode title
func GetEpisodeIDByFeedIDAndTitle(db *sql.DB, feedID int64, episodeTitle string) (id int64, err error) {
	stmt := "SELECT id FROM episodes WHERE feed_id = $1 AND title = $2"
	row := db.QueryRow(stmt, feedID, episodeTitle)
	err = row.Scan(&id)
	if err != nil {
		var title string
		stmt2 := "SELECT id, title FROM episodes WHERE feed_id = $1"
		rows, err := db.Query(stmt2, feedID)
		if err != nil {
			return id, err
		}

		for rows.Next() {
			err = rows.Scan(&id, &title)
			if err != nil {
				return id, err
			}

			if strings.Contains(title, episodeTitle) {
				return id, err
			}

		}
	}
	return
}

//GetFeedEpisodeIDs -- return ???
func GetFeedEpisodeIDs(db *sql.DB, feedID int64) (ids []int64, err error) {
	rows, err := db.Query("SELECT id FROM episodes WHERE feed_id = $1", feedID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}

//GetFeedEpisodeSeenRatio -- return the seen over total episode seen for a feed.
func GetFeedEpisodeSeenRatio(db *sql.DB, feedID int64) (seen, total int64, err error) {
	row := db.QueryRow("SELECT COUNT(*) FROM episodes WHERE feed_id = $1", feedID)
	err = row.Scan(&total)
	if err != nil {
		return
	}
	row = db.QueryRow("SELECT COUNT(*) FROM episodes WHERE feed_id = $1 AND seen = 1", feedID)
	err = row.Scan(&seen)
	if err != nil {
		return
	}
	return
}
