package cli

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var db_path string
var url_file string
var log_level string

func init() {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		slog.Warn(fmt.Sprintf("Unable to get user's home directory: %v", err))
	}

	default_db_path := filepath.Join(home_dir, ".go-rss-reader", "feeds.db")
	default_url_file := filepath.Join(home_dir, ".go-rss-reader", "urls")

	flag.StringVar(&db_path, "db", default_db_path, "Path to sqlite database")
	flag.StringVar(&url_file, "url-path", default_url_file, "Path to url file")
	flag.StringVar(&log_level, "log-level", "Info", "Valid levels are 'Debug', 'Info', 'Warn', 'Error'")

}

type Config interface {
	CliParse() error
	SetDBPath(string) error
	GetDBPath() string
	SetUrlFile(string) error
	GetUrlFile() string
	DBExist() bool
	LogLevel() slog.Level
}

type MyConfig struct {
	db_path   string
	url_file  string
	db_exist  bool
	log_level string
}

func (c *MyConfig) LogLevel() slog.Level {
	switch strings.ToLower(c.log_level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Warn(fmt.Sprintf("%v is not a Valid log level, so the log level will be set to Info", c.log_level))
		return slog.LevelInfo
	}

}

func (c *MyConfig) DBExist() bool {
	return c.db_exist
}

func (c *MyConfig) CliParse() error {
	slog.Debug("Starting CliParse")

	flag.Parse()

	// Need to re-set this after parsing flags
	c.log_level = log_level

	c.SetUrlFile(url_file)
	c.SetDBPath(db_path)

	return nil
}

func (c *MyConfig) SetDBPath(db_path string) error {
	file_info, err := os.Stat(db_path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Warn(fmt.Sprintf("Database file not found -> %v\n", err))

			dir, _ := filepath.Split(db_path)

			file_info, err = os.Stat(dir)
			if err != nil {
				return fmt.Errorf("Issue validating db path and directory: %w", err)
			}

			if file_info.IsDir() == true {
				slog.Info(fmt.Sprintf("Couldn't find DB file, but we have validated the directory: %s", dir))
				// The directory exists, so, when the time comes, we will create a new database
				c.db_path = db_path
				c.db_exist = false
				return nil
			}
		}

		// The DB path file exists, but I am still getting an error
		return fmt.Errorf("Unable to valid DB file path: %w", err)
	}

	if file_info.IsDir() == true {
		return fmt.Errorf("%s is a directory and not a path to file", db_path)
	}

	//We got a path to an existing file
	c.db_path = db_path
	c.db_exist = true
	return nil
}

func (c *MyConfig) GetDBPath() string {
	return c.db_path
}

func (c *MyConfig) SetUrlFile(url_file string) error {

	if len(url_file) == 0 {
		return errors.New("Must provide path to url file")
	}

	if _, err := os.Stat(url_file); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("The urls file was not found: %w", err)
		}
	}

	c.url_file = url_file

	return nil
}

func (c *MyConfig) GetUrlFile() string {
	return c.url_file
}

func NewConfig() Config {
	return &MyConfig{db_path: db_path, url_file: url_file, log_level: log_level}
}

func NewConfigWithValues(db_path, url_file, log_level string) Config {
	return &MyConfig{db_path: db_path, url_file: url_file, log_level: log_level}
}
