package scoirparser

import (
	"database/sql"
	"fmt"
	"time"

	// psql import
	_ "github.com/lib/pq"
)

//DBClient simple db struct for export/receiver funcs
type DBClient struct {
	DB *sql.DB
}

// NewDBClient returns a client for the Adsub database
func NewDBClient(user, pass, host, database string, port int) (*DBClient, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port, database,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &DBClient{DB: db}, nil
}

// CheckProcessedFile does a look up for the file and returns a bool if it exists
func (d DBClient) CheckProcessedFile(fileName string) (bool, error) {
	var selectedFileName string
	query := "select file_name from processed_files where file_name = $1"
	row := d.DB.QueryRow(query, fileName)
	switch err := row.Scan(&selectedFileName); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// InsertProcessedFile inserts the processed files name for a later lookup
func (d DBClient) InsertProcessedFile(fileName string) error {
	query := `insert into processed_files ( file_name, process_date)
			  values ($1, $2)`
	_, err := d.DB.Exec(
		query,
		fileName,
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
