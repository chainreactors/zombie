package sqlsess

import "database/sql"

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

