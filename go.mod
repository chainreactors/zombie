module github.com/chainreactors/zombie

go 1.18

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
	github.com/chainreactors/fingers v0.0.0-20240304115656-fa8ca9fc375f
	github.com/chainreactors/logs v0.0.0-20240207121836-c946f072f81f
	github.com/chainreactors/neutron v0.0.0-20240417160347-cb9446e38283
	github.com/chainreactors/parsers v0.0.0-20240415080936-e3e484abe2f7
	github.com/chainreactors/utils v0.0.0-20240302165634-2b8494c9cfc3
	github.com/chainreactors/words v0.4.1-0.20240126095632-02379f43c9f7
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/jessevdk/go-flags v1.5.0
	github.com/knadh/go-pop3 v0.3.0
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/streadway/amqp v1.1.0
	github.com/tomatome/grdp v0.0.0-00010101000000-000000000000
	github.com/vbauerster/mpb/v8 v8.7.2
	go.mongodb.org/mongo-driver v1.12.0
	golang.org/x/net v0.21.0
	sigs.k8s.io/yaml v1.4.0
)

require github.com/xinsnake/go-http-digest-auth-client v0.6.0

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/emersion/go-message v0.15.0 // indirect
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-dedup/megophone v0.0.0-20170830025436-f01be21026f5 // indirect
	github.com/go-dedup/simhash v0.0.0-20170904020510-9ecaca7b509c // indirect
	github.com/go-dedup/text v0.0.0-20170907015346-8bb1b95e3cb7 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/huin/asn1ber v0.0.0-20120622192748-af09f62e6358 // indirect
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

//github.com/go-sql-driver/mysql => ./external/github.com/go-sql-driver/mysql
//github.com/hirochachacha/go-smb2 => ./external/github.com/hirochachacha/go-smb2
replace github.com/tomatome/grdp => ./external/github.com/tomatome/grdp
