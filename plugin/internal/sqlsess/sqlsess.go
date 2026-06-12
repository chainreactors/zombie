package sqlsess

import (
	"database/sql"
	"fmt"
)

type Session struct {
	DB      *sql.DB
	SvcName string
}

func (s *Session) Service() string      { return s.SvcName }
func (s *Session) Close() error         { return s.DB.Close() }
func (s *Session) Raw() interface{}     { return s.DB }

func (s *Session) Query(query string, args ...any) ([][]string, error) {
	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	ncol := len(cols)

	var result [][]string
	result = append(result, cols)

	vals := make([]sql.NullString, ncol)
	ptrs := make([]interface{}, ncol)
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			continue
		}
		row := make([]string, ncol)
		for i, v := range vals {
			if v.Valid {
				row[i] = v.String
			}
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (s *Session) Databases() ([]string, error) {
	var query string
	switch s.SvcName {
	case "mysql":
		query = "SHOW DATABASES"
	case "postgresql":
		query = "SELECT datname FROM pg_database WHERE datistemplate = false"
	case "mssql":
		query = "SELECT name FROM sys.databases"
	case "oracle":
		query = "SELECT DISTINCT owner FROM all_tables"
	default:
		return nil, fmt.Errorf("unsupported service: %s", s.SvcName)
	}

	rows, err := s.Query(query)
	if err != nil {
		return nil, err
	}

	var dbs []string
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) > 0 {
			dbs = append(dbs, row[0])
		}
	}
	return dbs, nil
}
