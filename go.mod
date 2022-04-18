module Zombie

go 1.18

require (
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/go-ldap/ldap/v3 v3.4.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gosnmp/gosnmp v1.32.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/huin/asn1ber v0.0.0-20120622192748-af09f62e6358
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/lib/pq v1.9.0
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/sijms/go-ora/v2 v2.2.15
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20220331220935-ae2d96664a29
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20211209120228-48547f28849e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)

replace (
	github.com/go-sql-driver/mysql => ./pkg/github.com/go-sql-driver/mysql
	github.com/hirochachacha/go-smb2 => ./pkg/github.com/hirochachacha/go-smb2
)
