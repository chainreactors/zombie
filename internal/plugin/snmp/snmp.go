package snmp

import (
	"github.com/chainreactors/zombie/pkg"
	"github.com/gosnmp/gosnmp"
	"time"
)

type SnmpPlugin struct {
	*pkg.Task
	Input string
	conn  *gosnmp.GoSNMP
}

func (s *SnmpPlugin) Unauth() (bool, error) {
	conn := &gosnmp.GoSNMP{
		Target:             s.IP,
		Port:               s.UintPort(),
		Community:          "",
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(s.Timeout) * time.Second,
		MaxOids:            gosnmp.MaxOids,
		Retries:            3,
		ExponentialTimeout: true,
	}
	err := conn.Connect()
	if err != nil {
		return false, err
	}
	s.conn = conn
	return true, nil
}

//type CiderRoute struct {
//	Cidr    []string
//	GateWay []string
//}
//
//type IPSubRoute struct {
//	Cidr []string
//	IP   []string
//}
//
//type SwitchInfo struct {
//	SystemInfo     string   `json:"SystemInfo"`
//	Time           int64    `json:"Time"`
//	Concat         string   `json:"Concat"`
//	MachineName    string   `json:"MachineName"`
//	Location       string   `json:"Location"`
//	MemorySize     int64    `json:"MemorySize"`
//	SsCpuUser      int64    `json:"SsCpuUser"`
//	SsCpuSystem    int64    `json:"SsCpuSystem"`
//	SsCpuIdle      int64    `json:"SsCpuIdle"`
//	InterfaceSlice []string `json:"InterfaceSlice"`
//}

//func (s *SnmpPlugin) Query() bool {
//	defer s.conn.Conn.Close()
//
//	if strings.HasPrefix(s.Input, "Walk") {
//		input := strings.Replace(s.Input, "Walk", "", 1)
//		GetRes, err := s.conn.BulkWalkAll(input)
//		if err != nil {
//			return false
//		}
//		for _, alive := range GetRes {
//			fmt.Println(alive.Name)
//			if alive.Value != nil {
//				switch alive.Type {
//				case gosnmp.OctetString:
//					bytes := alive.Value.([]byte)
//					svalue := string(bytes)
//					fmt.Println(svalue)
//
//				default:
//					svalue := gosnmp.ToBigInt(alive.Value)
//					s2int := svalue.Int64()
//					fmt.Println(s2int)
//
//				}
//			}
//		}
//
//	} else {
//		GetRes, err := s.conn.Get([]string{s.Input})
//		if err != nil {
//			return false
//		}
//		variable := GetRes.Variables[0]
//		if variable.Value != nil {
//			switch variable.Type {
//			case gosnmp.OctetString:
//				bytes := variable.Value.([]byte)
//				svalue := string(bytes)
//				fmt.Println(svalue)
//
//			default:
//				svalue := gosnmp.ToBigInt(variable.Value)
//				s2int := svalue.Int64()
//				fmt.Println(s2int)
//
//			}
//		}
//	}
//
//	return true
//}

func (s *SnmpPlugin) SetQuery(query string) {
	s.Input = query
}

func (s *SnmpPlugin) Login() error {
	conn := &gosnmp.GoSNMP{
		Target:             s.IP,
		Port:               s.UintPort(),
		Community:          s.Password,
		Version:            gosnmp.Version2c,
		Timeout:            time.Duration(s.Timeout) * time.Second,
		MaxOids:            gosnmp.MaxOids,
		Retries:            3,
		ExponentialTimeout: true,
	}
	err := conn.Connect()
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *SnmpPlugin) Name() string {
	return s.Service
}

func (s *SnmpPlugin) GetResult() *pkg.Result {
	return &pkg.Result{Task: s.Task, OK: true}
}

func (s *SnmpPlugin) Close() error {
	if s.conn != nil {
		return s.conn.Conn.Close()
	}
	return nil
}
