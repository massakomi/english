package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Field struct {
	Field    string
	Type     string
	Nullable string
	Default  sql.NullString
}

var fieldsCache = make(map[string]map[string]Field)

func Fields(database *sqlx.DB, table string) map[string]Field {
	if fieldsCache[table] == nil {
		s := fmt.Sprintf(`SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = 'tester'
		AND table_name = '%v';`, table)
		data := GetFieldsBySql(database, s, func(rows *sql.Rows, p *Field) error {
			return rows.Scan(&p.Field, &p.Type, &p.Nullable, &p.Default)
		})
		fieldsCache[table] = data
	}
	return fieldsCache[table]
}

// GetFieldsBySql Общая механика выборки из таблицы
func GetFieldsBySql(database *sqlx.DB, sql string, scanner func(rows *sql.Rows, p *Field) error) map[string]Field {
	rows, err := database.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var items = make(map[string]Field)
	for rows.Next() {
		p := Field{}
		err := scanner(rows, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		items[p.Field] = p
	}
	return items
}
