package sqlsess

import (
	"fmt"
	"strings"
)

type dialect struct {
	schemaCol      string
	columnsTable   string
	excludeSchemas []string
	textTypes      []string
	sampleSQL      func(schema, table, col string, limit int) string
}

var dialects = map[string]*dialect{
	"mysql": {
		schemaCol:      "TABLE_SCHEMA",
		columnsTable:   "INFORMATION_SCHEMA.COLUMNS",
		excludeSchemas: []string{"mysql", "information_schema", "performance_schema", "sys"},
		textTypes:      []string{"varchar", "text", "char", "mediumtext", "longtext", "tinytext"},

		sampleSQL: func(schema, table, col string, limit int) string {
			return fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s` IS NOT NULL AND `%s` != '' LIMIT %d",
				col, schema, table, col, col, limit)
		},
	},
	"mssql": {
		schemaCol:      "TABLE_SCHEMA",
		columnsTable:   "INFORMATION_SCHEMA.COLUMNS",
		excludeSchemas: []string{"sys", "INFORMATION_SCHEMA"},
		textTypes:      []string{"varchar", "nvarchar", "text", "ntext", "char", "nchar"},

		sampleSQL: func(schema, table, col string, limit int) string {
			return fmt.Sprintf("SELECT TOP %d [%s] FROM [%s].[%s] WHERE [%s] IS NOT NULL AND [%s] != ''",
				limit, col, schema, table, col, col)
		},
	},
	"postgre": {
		schemaCol:      "table_schema",
		columnsTable:   "information_schema.columns",
		excludeSchemas: []string{"pg_catalog", "information_schema"},
		textTypes:      []string{"character varying", "text", "character", "name"},

		sampleSQL: func(schema, table, col string, limit int) string {
			return fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE "%s" IS NOT NULL AND "%s" != '' LIMIT %d`,
				col, schema, table, col, col, limit)
		},
	},
	"oracle": {
		schemaCol:      "OWNER",
		columnsTable:   "ALL_TAB_COLUMNS",
		excludeSchemas: []string{"SYS", "SYSTEM", "CTXSYS", "MDSYS", "OLAPSYS", "XDB", "WMSYS", "ORDDATA", "ORDSYS"},
		textTypes:      []string{"VARCHAR2", "NVARCHAR2", "CHAR", "NCHAR", "CLOB", "NCLOB"},

		sampleSQL: func(schema, table, col string, limit int) string {
			return fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE "%s" IS NOT NULL AND ROWNUM <= %d`,
				col, schema, table, col, limit)
		},
	},
}

func (s *Session) Audit(patterns []string, limit int) (map[string]string, error) {
	d, ok := dialects[s.SvcName]
	if !ok {
		return nil, fmt.Errorf("audit: unsupported SQL dialect %q", s.SvcName)
	}
	if limit <= 0 {
		limit = 100
	}

	query := buildDiscoverySQL(d, patterns)
	rows, err := s.Query(query)
	if err != nil {
		return nil, fmt.Errorf("audit discover: %w", err)
	}

	results := make(map[string]string)
	// rows[0] is the header row from sqlsess.Query
	for i, row := range rows {
		if i == 0 || len(row) < 3 {
			continue
		}
		schema, table, col := row[0], row[1], row[2]
		location := fmt.Sprintf("%s.%s.%s", schema, table, col)

		sampleRows, err := s.Query(d.sampleSQL(schema, table, col, limit))
		if err != nil {
			continue
		}
		var values []string
		for j, sr := range sampleRows {
			if j == 0 {
				continue
			}
			if len(sr) > 0 && sr[0] != "" {
				values = append(values, sr[0])
			}
		}
		if len(values) > 0 {
			results[location] = strings.Join(values, "\n")
		}
	}
	return results, nil
}

func buildDiscoverySQL(d *dialect, patterns []string) string {
	var likes []string
	for _, p := range patterns {
		likes = append(likes, fmt.Sprintf("LOWER(COLUMN_NAME) LIKE '%%%s%%'", strings.ReplaceAll(p, "'", "''")))
	}

	var excludes []string
	for _, s := range d.excludeSchemas {
		excludes = append(excludes, fmt.Sprintf("'%s'", s))
	}

	var types []string
	for _, t := range d.textTypes {
		types = append(types, fmt.Sprintf("'%s'", t))
	}

	return fmt.Sprintf(
		"SELECT %s, TABLE_NAME, COLUMN_NAME FROM %s WHERE (%s) AND %s NOT IN (%s) AND DATA_TYPE IN (%s)",
		d.schemaCol, d.columnsTable,
		strings.Join(likes, " OR "),
		d.schemaCol, strings.Join(excludes, ","),
		strings.Join(types, ","),
	)
}
