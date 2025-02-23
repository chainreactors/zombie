module github.com/chainreactors/zombie

go 1.11

require (
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874
	github.com/chainreactors/files v0.0.0-20240716182835-7884ee1e77f0
	github.com/chainreactors/fingers v1.0.1
	github.com/chainreactors/logs v0.0.0-20241030063019-8ca66a3ee307
	github.com/chainreactors/neutron v0.0.0-20250219105559-912bdcebda9a
	github.com/chainreactors/parsers v0.0.0-20240708072709-07deeece7ce2
	github.com/chainreactors/utils v0.0.0-20250109082818-178eed97b7ab
	github.com/chainreactors/words v0.0.0-20241002061906-25d8893158d9
	github.com/denisenkom/go-mssqldb v0.9.0
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gosnmp/gosnmp v1.32.0
	github.com/hirochachacha/go-smb2 v1.0.10
	github.com/jessevdk/go-flags v1.6.1
	github.com/jlaffaye/ftp v0.0.0-20201112195030-9aae4d151126
	github.com/knadh/go-pop3 v0.3.0
	github.com/lib/pq v1.9.0
	github.com/mitchellh/go-vnc v0.0.0-20150629162542-723ed9867aed
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/sijms/go-ora/v2 v2.2.15
	github.com/streadway/amqp v1.1.0
	github.com/vbauerster/mpb/v8 v8.7.2
	go.mongodb.org/mongo-driver v1.12.0
	golang.org/x/crypto v0.19.0
	golang.org/x/net v0.21.0
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lcvvvv/kscan/grdp v0.0.0-00010101000000-000000000000
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/xinsnake/go-http-digest-auth-client v0.6.0
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/lcvvvv/kscan/grdp => ./external/github.com/lcvvvv/grdp
	golang.org/x/crypto => github.com/golang/crypto v0.23.0
	golang.org/x/text => golang.org/x/text v0.12.0
)
