package ExecAble

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"os"
	"strings"
	"time"
)

type SnmpService struct {
	Utils.IpInfo
	Password string
	Cidr     []string
	GateWay  []string
	Input    string
	SwitchInfo
	SnmpCon *gosnmp.GoSNMP
}

type CiderRoute struct {
	Cidr    []string
	GateWay []string
}

type IPSubRoute struct {
	Cidr []string
	IP   []string
}

type SwitchInfo struct {
	SystemInfo     string   `json:"SystemInfo"`
	Time           int64    `json:"Time"`
	Concat         string   `json:"Concat"`
	MachineName    string   `json:"MachineName"`
	Location       string   `json:"Location"`
	MemorySize     int64    `json:"MemorySize"`
	SsCpuUser      int64    `json:"SsCpuUser"`
	SsCpuSystem    int64    `json:"SsCpuSystem"`
	SsCpuIdle      int64    `json:"SsCpuIdle"`
	InterfaceSlice []string `json:"InterfaceSlice"`
}

func SnmpConnect(Password string, info Utils.IpInfo) (err error, result bool, db *gosnmp.GoSNMP) {

	g := &gosnmp.GoSNMP{
		Target:             info.Ip,
		Port:               uint16(info.Port),
		Community:          Password,
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(Utils.Timeout/2) * time.Second,
		MaxOids:            gosnmp.MaxOids,
		Retries:            3,
		ExponentialTimeout: true,
	}
	err = g.Connect()

	if err != nil {
		//log.Println("Connect() err: %v", err)
		return err, false, nil
	}
	GetRes, err := g.Get([]string{".1.3.6.1.2.1.1.1.0"})

	if err != nil {
		result = false
	} else {

		variable := GetRes.Variables[0]
		if variable.Value != nil {
			result = true
		}

	}

	return err, result, g
}

func SnmpConnectTest(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	err, res, g := SnmpConnect(Password, info)

	if err == nil {
		result.Result = res
		_ = g.Conn.Close()
	}

	return err, result

}

func (s *SnmpService) Query() bool {
	defer s.SnmpCon.Conn.Close()

	if strings.HasPrefix(s.Input, "Walk") {
		input := strings.Replace(s.Input, "Walk", "", 1)
		GetRes, err := s.SnmpCon.BulkWalkAll(input)
		if err != nil {
			return false
		}
		for _, alive := range GetRes {
			fmt.Println(alive.Name)
			if alive.Value != nil {
				switch alive.Type {
				case gosnmp.OctetString:
					bytes := alive.Value.([]byte)
					svalue := string(bytes)
					fmt.Println(svalue)

				default:
					svalue := gosnmp.ToBigInt(alive.Value)
					s2int := svalue.Int64()
					fmt.Println(s2int)

				}
			}
		}

	} else {
		GetRes, err := s.SnmpCon.Get([]string{s.Input})
		if err != nil {
			return false
		}
		variable := GetRes.Variables[0]
		if variable.Value != nil {
			switch variable.Type {
			case gosnmp.OctetString:
				bytes := variable.Value.([]byte)
				svalue := string(bytes)
				fmt.Println(svalue)

			default:
				svalue := gosnmp.ToBigInt(variable.Value)
				s2int := svalue.Int64()
				fmt.Println(s2int)

			}
		}
	}

	return true
}

func (s *SnmpService) SetQuery(query string) {
	s.Input = query
}

func (s *SnmpService) Connect() bool {

	_, _, sn := SnmpConnect(s.Password, s.IpInfo)
	if sn != nil {
		s.SnmpCon = sn
		return true
	}
	return false
}

func (s *SnmpService) GetInfo() bool {
	cidr, _ := HandleinetCidrRouteEntry(s.SnmpCon)
	submask, _ := HandleIpSubmask(s.SnmpCon)
	var FinCidrSlice []string
	var FinIPSlice []string

	if submask != nil && len(submask.Cidr) != 0 {
		FinCidrSlice = append(FinCidrSlice, submask.Cidr...)
		if len(submask.IP) != 0 {
			FinIPSlice = append(FinIPSlice, submask.IP...)
		}
	}

	if cidr != nil && len(cidr.Cidr) != 0 {
		//fmt.Println(ip + " has CidrRoute")
		FinCidrSlice = append(FinCidrSlice, cidr.Cidr...)
		if len(submask.IP) != 0 {
			FinIPSlice = append(FinIPSlice, cidr.GateWay...)
		}
	}

	FinIPSlice = Utils.RemoveDuplicateElement(FinIPSlice)
	FinCidrSlice = Utils.RemoveDuplicateElement(FinCidrSlice)

	f, err1 := os.Create("./res/" + s.Ip + "Cidr.txt")
	if err1 != nil {
		panic(err1)
	}
	for _, resip := range FinCidrSlice {
		f.WriteString(resip + "\n")
	}

	f2, err2 := os.Create("./res/" + s.Ip + "AliveIP.txt")
	if err2 != nil {
		panic(err2)
	}
	for _, sub := range FinIPSlice {
		f2.WriteString(sub + "\n")
	}
	f.Close()
	f2.Close()

	s.Cidr = FinCidrSlice
	s.GateWay = FinIPSlice

	if Utils.More {
		s.SwitchInfo = *GetMoreInfo(s.SnmpCon)
	}
	return true
}
