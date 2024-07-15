module github.com/chainreactors/zombie

go 1.11

require (
	github.com/denisenkom/go-mssqldb v0.9.0
	//github.com/go-ldap/ldap/v3 v3.4.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gosnmp/gosnmp v1.32.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/lib/pq v1.9.0
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/sijms/go-ora/v2 v2.2.15
	golang.org/x/crypto v0.19.0
)

require (
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874
	github.com/chainreactors/files v0.0.0-20231123083421-cea5b4ad18a8
	github.com/chainreactors/fingers v0.0.0-20240702104653-a66e34aa41df
	github.com/chainreactors/logs v0.0.0-20240207121836-c946f072f81f
	github.com/chainreactors/neutron v0.0.0-20240712080924-c31f760d89d0
	github.com/chainreactors/parsers v0.0.0-20240708072709-07deeece7ce2
	github.com/chainreactors/utils v0.0.0-20240715080349-d2d0484c95ed
	github.com/chainreactors/words v0.4.1-0.20240126095632-02379f43c9f7
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/jessevdk/go-flags v1.5.0
	github.com/knadh/go-pop3 v0.3.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/streadway/amqp v1.1.0
	github.com/vbauerster/mpb/v8 v8.7.2
	go.mongodb.org/mongo-driver v1.12.0
	golang.org/x/net v0.21.0
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/lcvvvv/kscan/grdp v0.0.0-00010101000000-000000000000
	github.com/xinsnake/go-http-digest-auth-client v0.6.0
)

replace github.com/lcvvvv/kscan/grdp => ./external/github.com/lcvvvv/grdp

replace golang.org/x/crypto => github.com/golang/crypto v0.23.0

replace golang.org/x/text => golang.org/x/text v0.12.0
