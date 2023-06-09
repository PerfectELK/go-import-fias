package migrations

import (
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/shutdown"
	"fias_to_sql/migrations/mysql"
	"fias_to_sql/migrations/pgsql"
	"fias_to_sql/pkg/db"
)

func CreateTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbDriver := config.GetConfig("DB_DRIVER")

	fiasTableName := config.GetConfig("DB_OBJECTS_TABLE")
	fiasObjectTypesTableName := config.GetConfig("DB_OBJECT_TYPES_TABLE")
	fiasHierarchyTableName := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	fiasKladrTableName := config.GetConfig("DB_OBJECTS_KLADR_TABLE")

	_, tableCheck := dbInstance.Table(fiasTableName).Limit(1).Get()
	if tableCheck == nil &&
		shutdown.IsReboot &&
		config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT") == "original" {
		return nil
	}

	if tableCheck == nil {
		config.SetConfig("DB_ORIGINAL_OBJECTS_TABLE", config.GetConfig("DB_OBJECTS_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE", config.GetConfig("DB_OBJECT_TYPES_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE", config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE"))
		config.SetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE", config.GetConfig("DB_OBJECTS_KLADR_TABLE"))

		fiasTableName = config.GetConfig("DB_OBJECTS_TABLE") + "_temp"
		fiasObjectTypesTableName = config.GetConfig("DB_OBJECT_TYPES_TABLE") + "_temp"
		fiasHierarchyTableName = config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE") + "_temp"
		fiasKladrTableName = config.GetConfig("DB_OBJECTS_KLADR_TABLE") + "_temp"
		_, tempTableCheck := dbInstance.Table(fiasTableName).Limit(1).Get()

		config.SetConfig("DB_OBJECTS_TABLE", fiasTableName)
		config.SetConfig("DB_OBJECT_TYPES_TABLE", fiasObjectTypesTableName)
		config.SetConfig("DB_OBJECTS_HIERARCHY_TABLE", fiasHierarchyTableName)
		config.SetConfig("DB_OBJECTS_KLADR_TABLE", fiasKladrTableName)
		config.SetConfig("DB_TABLE_TYPES_FOR_IMPORT", "temp")
		if tempTableCheck == nil && !shutdown.IsReboot {
			return errors.New("fias tables and temp tables is exists")
		} else if tempTableCheck == nil {
			return nil
		}
	}

	return createFiasTables(
		dbDriver,
		fiasTableName,
		fiasObjectTypesTableName,
		fiasHierarchyTableName,
		fiasKladrTableName,
	)
}

func MigrateDataFromTempTables() error {
	if config.GetConfig("DB_TABLE_TYPES_FOR_IMPORT") != "temp" {
		return nil
	}

	dbDriver := config.GetConfig("DB_DRIVER")
	switch dbDriver {
	case "MYSQL":
		return mysql.MigrateFromTempTables()
	case "PGSQL":
		return pgsql.MigrateFromTempTables()
	default:
		return nil
	}
}

func createFiasTables(
	dbDriver string,
	fiasTableName string,
	fiasObjectTypesTableName string,
	fiasHierarchyTableName string,
	fiasKladrTableName string,
) error {
	switch dbDriver {
	case "MYSQL":
		err := mysql.ObjectsTableCreate(fiasTableName)
		if err != nil {
			return err
		}
		err = mysql.ObjectTypesTableCreate(fiasObjectTypesTableName)
		if err != nil {
			return err
		}
		err = mysql.HierarchyTableCreate(fiasHierarchyTableName)
		if err != nil {
			return err
		}
		return mysql.KladrTableCreate(fiasKladrTableName)
	case "PGSQL":
		err := pgsql.ObjectsTableCreate(fiasTableName)
		if err != nil {
			return err
		}
		err = pgsql.ObjectTypesTableCreate(fiasObjectTypesTableName)
		if err != nil {
			return err
		}
		err = pgsql.HierarchyTableCreate(fiasHierarchyTableName)
		if err != nil {
			return err
		}
		return pgsql.KladrTableCreate(fiasKladrTableName)
	default:
		return nil
	}
}
