package utils

import (
	"Go-lab/config"
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type DbUtils struct {
	DB     *sql.DB
	config *config.DBConfig
	mtx    sync.Mutex
}

func NewDbUtils(config *config.DBConfig) *DbUtils {
	res := &DbUtils{
		DB:     open(config),
		config: config,
	}
	return res
}

func (dbUtils *DbUtils) Init() error {
	dbUtils.mtx.Lock()
	defer dbUtils.mtx.Unlock()

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbUtils.createTables(ctx)
	if err != nil {
		return err
	}

	err = dbUtils.populateTables(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (dbUtils *DbUtils) WithTransaction(txFunc func(*sql.Tx) error) error {
	tx, err := dbUtils.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(tx)
	return err
}

func (dbUtils *DbUtils) Close() {
	log.Println("Closing the database...")

	dbUtils.mtx.Lock()
	defer dbUtils.mtx.Unlock()

	if dbUtils.DB == nil {
		return
	}

	err := dbUtils.DB.Close()
	if err != nil {
		log.Printf("Couldn't close the database :: %v", err)
	}

	dbUtils.DB = nil

	log.Println("Closed the database.")
}

func open(config *config.DBConfig) *sql.DB {
	log.Println("Opening the database...")

	// Open the database
	res, err := sql.Open(config.Driver, config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	res.SetConnMaxLifetime(time.Minute * 3)
	res.SetMaxOpenConns(10)
	res.SetMaxIdleConns(10)

	log.Println("Opened the database.")

	return res
}

//goland:noinspection SqlNoDataSourceInspection
func (dbUtils *DbUtils) createTables(ctx context.Context) error {
	log.Println("Creating the table...")

	err := dbUtils.WithTransaction(func(tx *sql.Tx) error {
		var sqlString string
		var err error

		sqlString = "DROP TABLE IF EXISTS [client_entity];"
		_, err = dbUtils.DB.ExecContext(ctx, sqlString)
		if err != nil {
			return err
		}

		sqlString = `CREATE TABLE IF NOT EXISTS [client_entity] (
			[id] INTEGER PRIMARY KEY AUTOINCREMENT,
			[account_no] TEXT NOT NULL UNIQUE,
			[account_name] TEXT NOT NULL,
			[created_at] DATETIME DEFAULT CURRENT_TIMESTAMP);`

		_, err = dbUtils.DB.ExecContext(ctx, sqlString)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	log.Println("Created the table.")

	return nil
}

//goland:noinspection SqlResolve,SqlNoDataSourceInspection
func (dbUtils *DbUtils) populateTables(ctx context.Context) error {
	data := [][]string{
		{"ABC001", "ABC Shoes"},
		{"XYZ001", "XYZ Trading As XYZ Enterprises"},
	}

	query := `INSERT INTO client_entity (account_no, account_name) VALUES (?, ?);`

	for _, row := range data {
		result, err := dbUtils.DB.ExecContext(ctx, query, row[0], row[1])
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			log.Println("No rows were inserted!")
		} else {
			log.Printf("rowsAffected=%d", affected)
		}
	}

	return nil
}
