package storage

import (
	"database/sql"
	"fmt"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	"github.com/midir99/sqload"
)

//go:embed db.sql
var dbScript string

var Q = sqload.MustLoadFromString[struct {
	InitDb              string `query:"InitDb"`
	FileInsertOrReplace string `query:"FileInsertOrReplace"`
	FileDeleteOne       string `query:"FileDeleteOne"`
	FileReadOne         string `query:"FileReadOne"`
	FileReadMany        string `query:"FileReadMany"`
}](dbScript)

type StoreConfig struct {
	Init bool
	URI  string
}

func ConnectDatabase(cfg StoreConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.URI)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}

	if cfg.Init {
		if _, err = db.Exec(Q.InitDb); err != nil {
			return db, err
		}
	}
	return db, err
}

type FileStore struct {
	db *sql.DB
}

func NewFileStore(conn *sql.DB) *FileStore {
	return &FileStore{db: conn}
}

func (store FileStore) createManyInTransaction(tx *sql.Tx, files []*File) error {
	stmt, err := tx.Prepare(Q.FileInsertOrReplace)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, file := range files {
		if err = stmt.QueryRow(file.Name, file.Path).Scan(&file.Id); err != nil {
			return err
		}
	}

	return nil
}

func (store FileStore) CreateMany(files []*File) error {
	db := store.db
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = store.createManyInTransaction(tx, files)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (store FileStore) Create(file *File) error {
	db := store.db
	err := db.QueryRow(Q.FileInsertOrReplace, file.Name, file.Path).Scan(&file.Id)
	return err
}

func (store FileStore) Delete(id int) error {
	db := store.db
	_, err := db.Exec(Q.FileDeleteOne, id)
	return err
}

func (store FileStore) Read(id int) (File, error) {
	db := store.db
	var file File
	err := db.QueryRow(
		Q.FileReadOne, id,
	).Scan(&file.Id, &file.Name, &file.Path)
	if err == sql.ErrNoRows {
		return file, fmt.Errorf("no files with id '%d' found", id)
	}
	return file, err

}

func (store FileStore) ReadMany(page int, pageSize int) ([]File, error) {
	db := store.db
	var (
		files  []File
		limit  = pageSize
		offset = 0
	)

	if page > 1 {
		offset = page * pageSize
	}

	rows, err := db.Query(
		Q.FileReadMany,
		limit, offset,
	)
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		var file File
		if err = rows.Scan(&file.Id, &file.Name, &file.Path); err != nil {
			return files, err
		}
		files = append(files, file)
	}

	return files, rows.Err()
}
