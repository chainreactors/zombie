package Server

import (
	"Zombie/src/Utils"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func TomcatConnect(User string, Password string, info Utils.IpInfo) (err error, result bool) {
	auth := User + ":" + Password

	auth = base64.StdEncoding.EncodeToString([]byte(auth))

	auth = "Basic " + auth
	var url string
	//var url2 string

	if info.SSL {
		url = fmt.Sprintf("https://%s:%d/manager/html", info.Ip, info.Port)
		//url2 = fmt.Sprintf("https://%s:%d/manager/html",info.Ip, info.Port)
	} else {
		url = fmt.Sprintf("http://%s:%d/manager/html", info.Ip, info.Port)
		//url2 = fmt.Sprintf("http://%s:%d/manager/html",info.Ip, info.Port)
	}

	var client = &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: false,
			DialContext: (&net.Dialer{
				Timeout: 2 * time.Second,
			}).DialContext,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := http.NewRequest("GET", url, nil)
	if err == nil {
		res.Header.Add("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36")
		res.Header.Add("Authorization", auth)
		res.Header.Add("Accept", "*/*")
		res.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
		res.Header.Add("Connection", "close")
		resp, err := client.Do(res)

		if err == nil {
			_, _ = io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == 200 {
				result = true
				return nil, result
			}
		}

	}

	return err, result
}
