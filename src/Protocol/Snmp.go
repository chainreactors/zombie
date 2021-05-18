package Protocol

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/alouca/gosnmp"
	"log"
)

func SnmpConnect(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {

	s, err := gosnmp.NewGoSNMP(info.Ip, Password, gosnmp.Version2c, int64(Utils.Timeout))
	if err != nil {
		log.Fatal(err)

	}
	resp, err := s.Get(".1.3.6.1.2.1.1.1.0")

	if err != nil {
		result.Result = false
	} else {
		result.Result = true
		fmt.Println(resp)
	}

	return err, result
}
