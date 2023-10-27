module github.com/chainreactors/zombie

go 1.12

require (
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/go-ldap/ldap/v3 v3.4.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gosnmp/gosnmp v1.32.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/lib/pq v1.9.0
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/sijms/go-ora/v2 v2.2.15
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
)

require (
	github.com/chainreactors/files v0.0.0-20230731174853-acee21c8c45a
	github.com/chainreactors/logs v0.0.0-20231027080134-7a11bb413460
	github.com/chainreactors/parsers v0.0.0-20231027075644-5d67f3a109d8
	github.com/chainreactors/utils v0.0.0-20230928082055-ef5e1ae47f3e
	github.com/chainreactors/words v0.4.1-0.20231027073512-0ccf7e0f0e32
	github.com/jessevdk/go-flags v1.5.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/tomatome/grdp v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.12.0
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/go-sql-driver/mysql => ./external/github.com/go-sql-driver/mysql
	github.com/hirochachacha/go-smb2 => ./external/github.com/hirochachacha/go-smb2
	github.com/tomatome/grdp => ./external/github.com/tomatome/grdp
)
