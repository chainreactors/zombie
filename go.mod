module Zombie

go 1.13

require (
	github.com/alouca/gologger v0.0.0-20120904114645-7d4b7291de9c // indirect
	github.com/alouca/gosnmp v0.0.0-20170620005048-04d83944c9ab
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/lib/pq v1.9.0
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
)

replace (
	github.com/go-sql-driver/mysql => ./pkg/github.com/go-sql-driver/mysql
	github.com/hirochachacha/go-smb2 => ./pkg/github.com/hirochachacha/go-smb2
)

