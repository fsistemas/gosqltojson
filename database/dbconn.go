package database

import (
	"database/sql"
	"fmt"
	"franciscoperez.dev/gosqltojson/core"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

type DBConn interface {
	Close() error
	Query(rawQuery string, parameters map[string]interface{}) ([]map[string]interface{}, error)
}

type dBConnImpl struct {
	sqlDB *sql.DB
	db    *gorm.DB
}

func createDialector(connectionType string, sqlDB *sql.DB) (gorm.Dialector, error) {
	if connectionType == "mysql" {
		return mysql.New(mysql.Config{
			Conn: sqlDB,
		}), nil
	}

	if connectionType == "postgres" {
		return postgres.New(postgres.Config{
			Conn: sqlDB,
		}), nil
	}

	return nil, fmt.Errorf("invalid connectionType: %s", connectionType)
}

func NewDBConn(configFile *core.ConfigFile, connectionName string) (DBConn, error) {
	rawConnectionString := configFile.GetConnection(connectionName)

	parts := strings.Split(rawConnectionString, "+")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid connection name or connection string '%s'", rawConnectionString)
	}

	connectionType := parts[0]
	connectionString := parts[1]

	var sqlDB *sql.DB
	var db *gorm.DB
	var err error

	if strings.HasPrefix(connectionType, "sqlite") {
		db, err = gorm.Open(sqlite.Dialector{DSN: connectionString}, &gorm.Config{})

		if err != nil {
			return nil, err
		}
	} else if connectionType == "mysql" || connectionType == "postgres" {
		sqlDB, err = sql.Open(connectionType, connectionString)

		if err != nil {
			return nil, err
		}

		dialector, err := createDialector(connectionType, sqlDB)

		if err != nil {
			return nil, err
		}

		db, err = gorm.Open(dialector)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid connection type: %s. valid types: mysql, postgres, sqlite", connectionType)
	}

	var dbConnection DBConn = dBConnImpl{
		sqlDB,
		db,
	}

	return dbConnection, nil
}

func (coon dBConnImpl) Close() error {
	if coon.sqlDB != nil {
		return coon.sqlDB.Close()
	}

	return nil
}

func (coon dBConnImpl) Query(queryRaw string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	var dbQuery *gorm.DB

	if len(parameters) > 0 {
		dbQuery = coon.db.Raw(queryRaw, parameters)
	} else {
		dbQuery = coon.db.Raw(queryRaw)
	}

	rows, err := dbQuery.Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return parseRowsToMaps(rows)
}

func parseRowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	values := make([]*string, len(columns))
	scanArgs := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})

		for i, colum := range columns {
			currentValue := *values[i]

			rowMap[colum] = currentValue
		}

		results = append(results, rowMap)
	}

	return results, nil
}
