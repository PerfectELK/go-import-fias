package mysql

import (
	"database/sql"
	"fmt"
	"github.com/PerfectELK/go-import-fias/internal/config"
	"github.com/PerfectELK/go-import-fias/pkg/db/helpers"
	"github.com/PerfectELK/go-import-fias/pkg/db/interfaces"
	"github.com/PerfectELK/go-import-fias/pkg/db/types"
	_ "github.com/go-sql-driver/mysql"
)

type Processor struct {
	db          *sql.DB
	isConnected bool
	table       string
	sel         []string
	where       [][]string
	limit       int
}

func (m *Processor) Connect(dbName ...string) error {
	connectStr := config.GetConfig("DB_USER") + ":" + config.GetConfig("DB_PASSWORD") + "@tcp(" + config.GetConfig("DB_HOST") + ":" + config.GetConfig("DB_PORT") + ")/"
	if len(dbName) > 0 {
		connectStr += dbName[0]
	}
	connectStr += "?multiStatements=true"
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return err
	}
	m.db = db
	m.isConnected = true
	return nil
}

func (m *Processor) Disconnect() error {
	m.isConnected = false
	return m.db.Close()
}

func (m *Processor) Exec(q string) error {
	rows, err := m.db.Query(q)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (m *Processor) Insert(table string, mm map[string]string) error {
	queryStr := "INSERT INTO " + table
	var keys []string
	var values []string
	for key, elem := range mm {
		if elem == "" {
			continue
		}
		keys = append(keys, key)
		values = append(values, elem)
	}

	keysStr := ""
	valuesStr := ""
	for key, _ := range keys {
		afterStr := ""
		if key != len(keys)-1 {
			afterStr += ", "
		}
		keysStr += keys[key] + afterStr
		valuesStr += "\"" + values[key] + "\"" + afterStr
	}

	queryStr += " (" + keysStr + ") VALUES (" + valuesStr + ");"
	return m.Exec(queryStr)
}

func (m *Processor) InsertList(table string, keys []types.Key, values [][]string) error {
	queryStr := "INSERT INTO " + table

	keysStr := ""
	valuesStr := ""
	for i, val := range keys {
		afterStr := ""
		if i != len(keys)-1 {
			afterStr += ", "
		}
		keysStr += val.Name + afterStr
	}
	queryStr += " (" + keysStr + ") "

	queryCount := 0
	for i, vals := range values {
		queryCount++
		valuesStr += "( "
		for key, val := range vals {
			afterStr := ""
			if key != len(vals)-1 {
				afterStr += ", "
			}
			valuesStr += "\"" + helpers.MysqlRealEscapeString(val) + "\"" + afterStr
		}
		closeStr := ") "
		if i != len(values)-1 && queryCount < 4000 {
			closeStr += ", "
		}
		valuesStr += closeStr
		if queryCount >= 4000 {
			q := queryStr + "VALUES " + valuesStr + ";"
			err := m.Exec(q)
			valuesStr = ""
			if err != nil {
				fmt.Println(err)
				return err
			}
			queryCount = 0
		}
	}

	if valuesStr != "" {
		return m.Exec(queryStr + "VALUES " + valuesStr + ";")
	}
	return nil
}

func (m *Processor) IsConnected() bool {
	return m.isConnected
}

func (m *Processor) Where(q [][]string) interfaces.DbProcessor {
	m.where = q
	return m
}

func (m *Processor) Table(t string) interfaces.DbProcessor {
	m.table = t
	return m
}

func (m *Processor) Select(s []string) interfaces.DbProcessor {
	m.sel = s
	return m
}

func (m *Processor) Limit(l int) interfaces.DbProcessor {
	m.limit = l
	return m
}

func (m *Processor) Get() (*sql.Rows, error) {
	return nil, nil
}

func (m *Processor) Use(q string) error {
	m.db.Close()
	m.isConnected = false
	return m.Connect(q)
}

func (m *Processor) Query(q string) (*sql.Rows, error) {
	return m.db.Query(q)
}

func (m *Processor) GetDriverName() string {
	return "MYSQL"
}
