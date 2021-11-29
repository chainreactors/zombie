module Zombie

go 1.16

require (
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gosnmp/gosnmp v1.32.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/lib/pq v1.9.0
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/sijms/go-ora/v2 v2.2.15
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	github.com/huin/asn1ber v0.0.0-20120622192748-af09f62e6358
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
)

replace (
	github.com/go-sql-driver/mysql => ./pkg/github.com/go-sql-driver/mysql
	github.com/hirochachacha/go-smb2 => ./pkg/github.com/hirochachacha/go-smb2
)
