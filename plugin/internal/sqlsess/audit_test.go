package sqlsess

import "testing"

func TestPostgreSQLAuditDialectName(t *testing.T) {
	if _, ok := dialects["postgresql"]; !ok {
		t.Fatal("postgresql session name has no audit dialect")
	}
}
