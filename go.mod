module Zombie

go 1.12

require (
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/go-ldap/ldap/v3 v3.4.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gosnmp/gosnmp v1.32.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5 // indirect
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/lib/pq v1.9.0
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/sijms/go-ora/v2 v2.2.15
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
)

require (
	github.com/chainreactors/grdp v0.0.0-20220523130340-cb6039da6de3
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
)

replace (
	github.com/go-sql-driver/mysql => ./pkg/github.com/go-sql-driver/mysql
	github.com/hirochachacha/go-smb2 => ./pkg/github.com/hirochachacha/go-smb2
//github.com/icodeface/grdp => github.com/chainreactors/grdp
)
