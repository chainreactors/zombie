package ExecAble

import "database/sql"

func DoRowsMapper(rows *sql.Rows) ([]map[string]string, []string) {

	var result []map[string]string
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// 这个map用来存储一行数据，列名为map的key，map的value为列的值
		rowMap := make(map[string]string)
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col != nil {
				value = string(col)
				rowMap[columns[i]] = value
			}
		}
		result = append(result, rowMap)

	}
	return result, columns
}

type ExecAble interface {
	Query() bool
	GetInfo() bool
	Connect() bool
	SetQuery(string)
}
