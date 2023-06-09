package pgsql

import (
	"fias_to_sql/internal/config"
	"fias_to_sql/pkg/db"
	"fmt"
)

func ObjectsTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	err = dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + tableName + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INTEGER NOT NULL DEFAULT 0," +
			"object_guid VARCHAR(100) NOT NULL DEFAULT ''," +
			"type_name VARCHAR(100) NOT NULL DEFAULT ''," +
			"level INT NOT NULL DEFAULT 0," +
			"name VARCHAR(255) NOT NULL DEFAULT '');",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_name_index ON " + dbSchema + "." + tableName + " (name);" +
			" CREATE INDEX " + tableName + "_object_guid_index ON " + dbSchema + "." + tableName + " (object_guid);" +
			" CREATE INDEX " + tableName + "_object_id_index ON " + dbSchema + "." + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_type_name_index ON " + dbSchema + "." + tableName + " (type_name);",
	)
}

func ObjectTypesTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	return dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + tableName + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"level INT NOT NULL DEFAULT 0," +
			"short_name VARCHAR(255) NOT NULL DEFAULT ''," +
			"name VARCHAR(255) NOT NULL DEFAULT '');",
	)
}

func HierarchyTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	err = dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + tableName + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INT NOT NULL DEFAULT 0," +
			"parent_object_id INT NOT NULL DEFAULT 0);",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_object_id_index ON " + dbSchema + "." + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_parent_object_id_index ON " + dbSchema + "." + tableName + " (parent_object_id);",
	)
}

func KladrTableCreate(tableName string) error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}
	dbSchema := config.GetConfig("DB_SCHEMA")
	err = dbInstance.Exec(
		"CREATE TABLE " + dbSchema + "." + tableName + " (" +
			"id BIGSERIAL PRIMARY KEY," +
			"object_id INT NOT NULL DEFAULT 0," +
			"kladr_id VARCHAR(50) NOT NULL DEFAULT '');",
	)
	if err != nil {
		return err
	}

	return dbInstance.Exec(
		"CREATE INDEX " + tableName + "_object_id_index ON " + dbSchema + "." + tableName + " (object_id);" +
			" CREATE INDEX " + tableName + "_kladr_id_index ON " + dbSchema + "." + tableName + " (kladr_id);",
	)
}

func MigrateFromTempTables() error {
	err := dropOldTables()
	if err != nil {
		return err
	}
	err = renameTables()
	if err != nil {
		return err
	}
	return renameIndexes()
}

func dropOldTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbSchema := config.GetConfig("DB_SCHEMA")

	originalObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE"))
	originalObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE"))
	originalHierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE"))
	originalFiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE"))

	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalObjectTypesTableName + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalHierarchyObjectsTable + ";")
	if err != nil {
		return err
	}
	err = dbInstance.Exec("DROP TABLE IF EXISTS " + originalFiasKladrTableName + ";")
	return err
}

func renameTables() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	dbSchema := config.GetConfig("DB_SCHEMA")

	originalObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE"))
	originalObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECT_TYPES_TABLE"))
	originalHierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE"))
	originalFiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE"))

	tempObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECTS_TABLE"))
	tempObjectTypesTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECT_TYPES_TABLE"))
	tempHierarchyObjectsTable := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE"))
	tempFiasKladrTableName := fmt.Sprintf("%s.%s", dbSchema, config.GetConfig("DB_OBJECTS_KLADR_TABLE"))

	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempObjectsTable, originalObjectsTable))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempObjectTypesTableName, originalObjectTypesTableName))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempHierarchyObjectsTable, originalHierarchyObjectsTable))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s RENAME TO %s", tempFiasKladrTableName, originalFiasKladrTableName))
	return err
}

func renameIndexes() error {
	dbInstance, err := db.GetDbInstance()
	if err != nil {
		return err
	}

	originalObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_TABLE")
	originalHierarchyObjectsTable := config.GetConfig("DB_ORIGINAL_OBJECTS_HIERARCHY_TABLE")
	originalFiasKladrTableName := config.GetConfig("DB_ORIGINAL_OBJECTS_KLADR_TABLE")

	tempObjectsTable := config.GetConfig("DB_OBJECTS_TABLE")
	tempHierarchyObjectsTable := config.GetConfig("DB_OBJECTS_HIERARCHY_TABLE")
	tempFiasKladrTableName := config.GetConfig("DB_OBJECTS_KLADR_TABLE")

	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_name_index RENAME TO %s_name_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_guid_index RENAME TO %s_object_guid_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_id_index RENAME TO %s_object_id_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_type_name_index RENAME TO %s_type_name_index",
		tempObjectsTable,
		originalObjectsTable,
	))
	if err != nil {
		return err
	}

	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_id_index RENAME TO %s_object_id_index",
		tempHierarchyObjectsTable,
		originalHierarchyObjectsTable,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_parent_object_id_index RENAME TO %s_parent_object_id_index",
		tempHierarchyObjectsTable,
		originalHierarchyObjectsTable,
	))
	if err != nil {
		return err
	}

	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_object_id_index RENAME TO %s_object_id_index",
		tempFiasKladrTableName,
		originalFiasKladrTableName,
	))
	if err != nil {
		return err
	}
	err = dbInstance.Exec(fmt.Sprintf(
		"ALTER INDEX %s_kladr_id_index RENAME TO %s_kladr_id_index",
		tempFiasKladrTableName,
		originalFiasKladrTableName,
	))
	return err
}
