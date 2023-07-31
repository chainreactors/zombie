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
	github.com/chainreactors/files v0.2.5-0.20230731172931-8d0a64e79d98
	github.com/chainreactors/ipcs v0.0.13
	github.com/chainreactors/logs v0.6.1
	github.com/chainreactors/parsers v0.3.0
	github.com/chainreactors/words v0.4.1-0.20230731130315-158c047fd378
	github.com/jessevdk/go-flags v1.5.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	go.mongodb.org/mongo-driver v1.12.0 // indirect
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/go-sql-driver/mysql => ./external/github.com/go-sql-driver/mysql
	github.com/hirochachacha/go-smb2 => ./external/github.com/hirochachacha/go-smb2
//github.com/icodeface/grdp => github.com/chainreactors/grdp
)
