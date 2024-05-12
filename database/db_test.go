package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3" //Sqlite3 driver
)

type connectorErr struct{}

func (c connectorErr) Exec(query string, args ...any) (sql.Result, error) {
	return nil, errors.New("Exec Error")
}

func (c connectorErr) Ping() error {
	return errors.New("Ping Error")
}

func (c connectorErr) QueryRow(query string, args ...any) *sql.Row {
	return nil
}

func (c connectorErr) Close() error {
	return errors.New("Connector failed to close")
}

type driverOpenErr struct{}

func (d driverOpenErr) Open(drivername, dsn string) (Connector, error) {
	return nil, errors.New("Open Error")
}

type driverPingErr struct{}

func (d driverPingErr) Open(name, dsn string) (Connector, error) {
	return connectorErr{}, nil
}

func TestCreateTablesV2(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name string
		conn Connector
		err  error
	}{
		{"Exec error", connectorErr{}, errors.New("Creating tables failed with the following error: Exec Error")},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := createTablesV2(tc.conn)
			// TODO: Figure out success test cases
			if err == nil {
				t.Fatalf("Expected err to be nil, but got %s", err.Error())
			}
			if err != nil {
				if err.Error() != tc.err.Error() {
					t.Fatalf("Got %q, but expected %q", err.Error(), tc.err.Error())
				}
			}
		})

	}

}

func TestInitV2HappyPath(t *testing.T) {
	t.Parallel()

	driverName := "sqlite3"

	tcs := []struct {
		name  string
		dsn   string
		reset bool
	}{
		{"In Memory DB", "file:memory.db?_foreign_keys=1&mode=memory", false},
		{"In Memory DB Create Tables", "file:memory.db?_foreign_keys=1&mode=memory", true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := InitV2(driverName, tc.dsn, false, func(a, b string) (Connector, error) {
				conn, err := sql.Open(a, b)
				if err != nil {
					return conn, err
				}
				return conn, err
			})
			if err != nil {
				t.Fatalf("Exepected nil, but got %q", err.Error())
			}

		})
	}

}

func TestInitV2(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name       string
		driverName string
		dsn        string
		reset      bool
		driver     Driver
		err        error
	}{
		{"Open Error", "driverName", "dsn", false, driverOpenErr{}, errors.New("Unable to Open Database: Open Error")},
		{"Ping Error", "driverName", "dsn", false, driverPingErr{}, errors.New("Unable to ping the Database: Ping Error")},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := InitV2(tc.driverName, tc.dsn, tc.reset, tc.driver.Open)

			if err == nil {
				t.Fatalf("Expected err to be nil, but got %s", err.Error())
			}
			if err != nil {
				if err.Error() != tc.err.Error() {
					t.Fatalf("Got %q, but expected %q", err.Error(), tc.err.Error())
				}
			}
		})
	}
}

func TestCreateDBDsn(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		path     string
		inMemory bool
		expected string
	}{
		{"File Path string", "testing.db", false, "file:testing.db?_foreign_keys=1"},
		{"In Memory string", "memory.db", true, "file:memory.db?_foreign_keys=1&mode=memory"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := CreateDBDsn(tc.path, tc.inMemory)
			if strings.Compare(result, tc.expected) != 0 {
				t.Fatalf("Expectd %q, but got %q", tc.expected, result)
			}
		})
	}
}

func createTestDB(file string) *sql.DB {
	testDB := fmt.Sprintf("file:%s%s", file, foreignKeySupport)

	db, err := Init(testDB, true)
	if err != nil {
		log.Fatalf("Error when trying to create the database (%s): %s", file, err.Error())
	}
	return db
}

func TestExist(t *testing.T) {
	tests := []struct {
		File     string
		Expected bool
	}{
		{"db.go", true},
		{"DoesNotExist", false},
	}

	for _, test := range tests {
		result := Exist(test.File)

		if result != test.Expected {
			t.Errorf("For file %s, expected existence was %t, but got %t", test.File, test.Expected, result)
		}
	}
}

func TestCreate(t *testing.T) {
	file := "./testing/create_test_file.db"

	//Need to create the test db file
	_, err := os.Create(file)
	if err != nil {
		t.Errorf("Unexpected error when create the database: %s", err.Error())
	}

	if !Exist(file) {
		t.Errorf("File: %s does not exist", file)
	}

	err = os.Remove(file)
	if err != nil {
		t.Errorf("Error while removing file (%s): %s", file, err.Error())
	}
}

func TestInit(t *testing.T) {
	file := "./testing/init_test_file.db"
	dbPath := fmt.Sprintf("file:%s?_foreign_keys=1", file)

	tests := []struct {
		File  string
		Reset bool
	}{
		{dbPath, false},
		{dbPath, true},
	}

	for _, test := range tests {

		_, err := Init(test.File, test.Reset)

		if err != nil {
			t.Errorf("Unexpected err: %s", err.Error())
		}
	}
}
