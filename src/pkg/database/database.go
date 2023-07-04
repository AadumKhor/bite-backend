// Package database is the actual implementation of the DB
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AadumKhor/bitespeed-backend-task/src/pkg/utils"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	pingQuery = "SELECT 1"
)

type (
	// PGStore represents the DB group containing both read & write DB
	PGStore struct {
		master  *bun.DB
		replica *bun.DB
	}
)

var pgStore *PGStore

func newPGStore(config utils.Config) (*PGStore, error) {
	// Singleton pattern followed here
	if pgStore == nil {
		var readDB, writeDB *bun.DB

		// iterate over the database config and create connections
		/*
			NOTE: This is not optimal since connection count is not
			taken into consideration

			This is done only for this task
		*/
		for _, dbConfig := range config.Databases {
			dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
				dbConfig.User,
				dbConfig.Password,
				dbConfig.Host,
				dbConfig.Port,
				dbConfig.Name)

			sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dataSourceName)))
			db := bun.NewDB(sqldb, pgdialect.New())

			switch dbConfig.Type {
			case "read":
				readDB = db
			case "write":
				writeDB = db
			default:
				err := fmt.Sprintf("Unknown database type: %s", dbConfig.Type)
				return nil, errors.New(err)
			}
		}

		pgStore = &PGStore{
			master:  writeDB,
			replica: readDB,
		}
	}

	return pgStore, nil
}

// GetPGStore returns the store reference
func GetPGStore() *PGStore {
	return pgStore
}

// Connect is a wrapper function to establish connection with DB
func Connect(config utils.Config) error {
	store, err := newPGStore(config)
	if err != nil {
		return err
	}

	// Check master connection
	master := store.master
	if err := ping(master); err != nil {
		return err
	}

	// Check read replica connection
	replica := store.replica
	if err := ping(replica); err != nil {
		return err
	}

	return err
}

func ping(db *bun.DB) error {
	_, err := db.Exec(pingQuery)
	return err
}
