package migrations

import (
	"errors"
	"fmt"
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/pkg/db"
	"github.com/PerfectELK/go-import-fias/pkg/db/helpers"
)

func CreateDatabase() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbName := config.GetConfig("DB_NAME")
	switch dbInstance.GetDriverName() {
	case "PGSQL":
		rows, err := dbInstance.Query(fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", dbName))
		if err != nil {
			return err
		}
		rowsArr := helpers.Scan(rows)
		if len(rowsArr) == 0 {
			return dbInstance.Exec(fmt.Sprintf("CREATE DATABASE \"%s\";", dbName))
		}

		err = dbInstance.Use(dbName)
		if err != nil {
			return err
		}
		dbSchema := config.GetConfig("DB_SCHEMA")
		return dbInstance.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", dbSchema))
	case "MYSQL":
		err = dbInstance.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		if err != nil {
			return err
		}
		return dbInstance.Use(dbName)
	default:
		return errors.New("doesn't selected db driver")
	}

}
