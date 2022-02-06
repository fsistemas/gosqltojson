package database

import (
	"database/sql"
	"flag"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"franciscoperez.dev/gosqltojson/core"
	"franciscoperez.dev/gosqltojson/parameter"
)

func RunQuery(runConfig core.RunConfig) ([]map[string]string, error) {
	conectionConfig, err := core.NewConfigFile(runConfig.ConfigFile)

	if err != nil {
		return nil, err
	}

	rawConectionString := conectionConfig.GetConection(runConfig.ConectionName)
	queryRaw := conectionConfig.GetQuery(runConfig.QueryName)

	parts := strings.Split(rawConectionString, "+")

	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid conection string")
	}

	conectionType := parts[0]
	conectionString := parts[1]

	var sqlDB *sql.DB
	var db *gorm.DB

	if conectionType == "mysql" {
		sqlDB, err = sql.Open(conectionType, conectionString)

		if err != nil {
			return nil, err
		}

		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

		if err != nil {
			return nil, err
		}
	} else if conectionType == "postgres" {
		sqlDB, err = sql.Open(conectionType, conectionString)

		if err != nil {
			return nil, err
		}

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(conectionType, "sqlite") {
		db, err = gorm.Open(sqlite.Open(conectionString), &gorm.Config{})

		if err != nil {
			return nil, err
		}
	} else {
		return nil, (fmt.Errorf("invalid conection type: %s. valid types: mysql, postgres, sqlite", conectionType))
	}

	if sqlDB != nil {
		defer sqlDB.Close()
	}

	queryParams, err := parameter.ParseQueryParams(queryRaw, flag.Args())
	if err != nil {
		return nil, err
	}

	var dbQuery *gorm.DB

	if len(queryParams) > 0 {
		dbQuery = db.Raw(queryRaw, queryParams)
	} else {
		dbQuery = db.Raw(queryRaw)
	}

	rows, err := dbQuery.Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return parseRowsToMaps(rows)
}

func parseRowsToMaps(rows *sql.Rows) ([]map[string]string, error) {
	var results []map[string]string

	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	scan_args := make([]interface{}, len(columns))

	for i := range values {
		scan_args[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scan_args...)

		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]string)

		for i, colum := range columns {
			rowMap[colum] = fmt.Sprintf("%s", values[i])
		}

		results = append(results, rowMap)
	}

	return results, nil
}
