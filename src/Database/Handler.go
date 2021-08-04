package Database

import (
	"Zombie/src/Utils"
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
)

func GetMoreInfo(spcon *gosnmp.GoSNMP) *SwitchInfo {
	Res := &SwitchInfo{}

	var OidMap = map[string]string{
		"SystemInfo":  ".1.3.6.1.2.1.1.1.0",
		"Time":        ".1.3.6.1.2.1.1.3.0",
		"Concat":      ".1.3.6.1.2.1.1.4.0",
		"MachineName": ".1.3.6.1.2.1.1.5.0",
		"Location":    ".1.3.6.1.2.1.1.6.0",
		"MemorySize":  ".1.3.6.1.2.1.25.2.2.0",
		"SsCpuUser":   ".1.3.6.1.4.1.2021.11.9.0",
		"SsCpuSystem": ".1.3.6.1.4.1.2021.11.10.0",
		"SsCpuIdle":   ".1.3.6.1.4.1.2021.11.11.0",
	}
	v := reflect.ValueOf(Res).Elem()
	for name, oid := range OidMap {
		GetRes, err := spcon.Get([]string{oid})
		if err != nil {
			continue
		}
		variable := GetRes.Variables[0]
		if variable.Value != nil {
			switch variable.Type {
			case gosnmp.OctetString:
				bytes := variable.Value.([]byte)
				svalue := string(bytes)
				if v.FieldByName(name).Type() != reflect.TypeOf(svalue) {
					continue
				}
				v.FieldByName(name).Set(reflect.ValueOf(svalue))

			default:
				svalue := gosnmp.ToBigInt(variable.Value)
				s2int := svalue.Int64()
				if v.FieldByName(name).Type() != reflect.TypeOf(svalue) {
					continue
				}
				v.FieldByName(name).Set(reflect.ValueOf(s2int))

			}
		}

	}

	return Res
}

func HandleinetCidrRouteEntry(spcon *gosnmp.GoSNMP) (*CiderRoute, error) {
	var OidMap = map[string]string{
		"GeneralOid": ".1.3.6.1.2.1.4.24.7.1.12",
		"BackOid":    ".1.3.6.1.2.1.4.24.7.1.13",
	}

	result := CiderRoute{}

	for _, oid := range OidMap {
		IPAlive, err := spcon.BulkWalkAll(oid)
		if err != nil {
			continue
		}
		for _, alive := range IPAlive {
			if alive.Name == oid {
				continue
			}
			OidName := alive.Name
			HandledOidName := strings.Replace(OidName, oid+".", "", 1)
			Cider := HandleCiderFromRoute(HandledOidName)
			if strings.HasPrefix(Cider, "127.0.0.1") || strings.HasSuffix(Cider, "/0") {
				continue
			}
			result.Cidr = append(result.Cidr, Cider)
			GateWay := GetIPFromOid(HandledOidName)
			result.GateWay = append(result.GateWay, GateWay)

		}

	}
	result.Cidr = Utils.RemoveDuplicateElement(result.Cidr)
	result.GateWay = Utils.RemoveDuplicateElement(result.GateWay)

	return &result, nil

}

func HandleCiderFromRoute(HandledName string) string {
	var cidr string
	OidList := strings.Split(HandledName, ".")

	for i := 2; i <= 6; i++ {
		cidr += OidList[i]
		if i == 5 {
			cidr += "/"
		} else {
			cidr += "."
		}
	}
	return cidr[:len(cidr)-1]
}

func HandleipNetToMediaEntry(spcon *gosnmp.GoSNMP) (*[]string, error) {
	var result []string
	oid := ".1.3.6.1.2.1.4.22.1.3"

	IPAlive, err := spcon.BulkWalkAll(oid)

	if err != nil {
		return nil, err
	}

	for _, alive := range IPAlive {
		ip := alive.Value

		result = append(result, ip.(string))
	}

	return &result, nil
}

// IPNetMedia获取的和Physical互为替代品

func HandleipNetToPhysical(spcon *gosnmp.GoSNMP) (*[]string, error) {
	oid := ".1.3.6.1.2.1.4.35.1.5"

	var result []string

	IpList, err := spcon.BulkWalkAll(oid)

	if err != nil {
		return nil, err
	}

	for _, info := range IpList {
		name := info.Name
		ip := GetIPFromOid(name)

		if ip == "" {
			continue
		}
		result = append(result, ip)
	}
	return &result, nil
}

func HandleIpSubmask(spcon *gosnmp.GoSNMP) (*IPSubRoute, error) {

	oid := ".1.3.6.1.2.1.4.20.1.3"

	SubMask, err := spcon.BulkWalkAll(oid)

	res := IPSubRoute{}

	var AddList []string

	if err != nil {
		return nil, err
	}

	for _, subinfo := range SubMask {
		name := subinfo.Name
		value := subinfo.Value
		ip := GetIPFromOid(name)

		if ip == "" {
			continue
		}

		res.IP = append(res.IP, ip)

		switch value.(type) {
		case string:
			submask, err2 := ipMaskToInt(value.(string))
			if err2 != nil {
				log.Println(err2)
				continue
			}
			StrSubMask := strconv.Itoa(submask)

			IpSub := ip + "/" + StrSubMask

			_, FinIp, _ := net.ParseCIDR(IpSub)

			CurIP := FinIp.String()

			AddList = append(AddList, CurIP)
		case []uint8:
			fmt.Println(spcon.Target + "is different")
		}

	}

	res.Cidr = AddList

	return &res, err

}

func GetIPFromOid(OidName string) string {
	var ip string
	oidList := strings.Split(OidName, ".")
	for i := 4; i >= 1; i-- {
		ip += oidList[len(oidList)-i] + "."
	}
	ip = ip[:len(ip)-1]
	if ip != "127.0.0.1" && ip != "0.0.0.0" {
		return ip
	}
	return ""
}

func ipMaskToInt(netmask string) (int, error) {
	ipSplitArr := strings.Split(netmask, ".")
	if len(ipSplitArr) != 4 {
		return 0, fmt.Errorf("netmask:%v is not valid, pattern should like: 255.255.255.0", netmask)
	}
	ipv4MaskArr := make([]byte, 4)
	for i, value := range ipSplitArr {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("ipMaskToInt call strconv.Atoi error:[%v] string value is: [%s]", err, value)
		}
		if intValue > 255 {
			return 0, fmt.Errorf("netmask cannot greater than 255, current value is: [%s]", value)
		}
		ipv4MaskArr[i] = byte(intValue)
	}

	ones, _ := net.IPv4Mask(ipv4MaskArr[0], ipv4MaskArr[1], ipv4MaskArr[2], ipv4MaskArr[3]).Size()
	return ones, nil
}
