package scoirparser

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

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

func (d DBClient) CheckProcessedFile(fileName string) (bool, error) {
	var selectedFileName string
	query := "select file_name from processed_files where file_name = $1"
	row := d.DB.QueryRow(query, fileName)
	switch err := row.Scan(&selectedFileName); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		fmt.Println(selectedFileName)
		return true, nil
	default:
		return false, err
	}
}

func (d DBClient) InsertProcessedFile(fileName string) error {
	query := `insert into processed_files ( file_name, process_date)
			  values ($1, $2)`
	something, err := d.DB.Exec(
		query,
		fileName,
		time.Now(),
	)
	fmt.Println(something)
	if err != nil {
		return err
	}
	return nil
}
