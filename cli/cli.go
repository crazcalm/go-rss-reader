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

func init (){
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
	DbPath() string
	UrlFile() string
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
	
	//flag.StringVar(&c.db_path, "db", default_db_path, "Path to sqlite database")
	//flag.StringVar(&c.url_file, "url-path", default_url_file, "Path to url file")
	//log_level := flag.String("log-level", "Info", "Valid levels are 'Debug', 'Info', 'Warn', 'Error'")

	flag.Parse()

	// Need to re-set this after parsing flags
	c.log_level = log_level
	
	//slog.Info("testing slog attr out", "log_level", *log_level)
	
	if c.url_file == "" {
		return errors.New("Must provide path to url file")
	}

	if _, err := os.Stat(c.url_file); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("The urls file was not found: %w", err)
		}
	}

	if _, err := os.Stat(c.db_path); err != nil {
		if os.IsNotExist(err) {
			//TODO figure out logging
			slog.Warn(fmt.Sprintf("Database file not found -> %v\n", err))
		}

		if err != nil {
			dir, _ := filepath.Split(c.db_path)

			if _, err = os.Stat(dir); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("Either the DB file or the directory the file is placed in must exist: %w\n", err)
				}
			}
		}
	} else {
		// Set db_exist to true
		c.db_exist = true
	}

	return nil
}

func (c *MyConfig) DbPath() string {
	return c.db_path
}

func (c *MyConfig) UrlFile() string {
	return c.url_file
}

func NewConfig() Config {
	return &MyConfig{db_path: db_path, url_file: url_file, log_level: log_level}
}
