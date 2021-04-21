package Web

import (
	"Zombie/src/Utils"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func EsConnect(User string, Password string, info Utils.IpInfo) (err error, result Utils.BruteRes) {
	auth := User + ":" + Password

	auth = base64.StdEncoding.EncodeToString([]byte(auth))

	auth = "Basic " + auth
	var url string

	if info.SSL {
		url = fmt.Sprintf("https://%s:%d/_cat/", info.Ip, info.Port)

	} else {
		url = fmt.Sprintf("http://%s:%d/_cat/", info.Ip, info.Port)

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
		res.Header.Add("Accept", "*/*")
		res.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
		res.Header.Add("Connection", "close")
		resp, err := client.Do(res)

		if err == nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if strings.Contains(string(body), "/_cat/master") {
				//AnnonTest := fmt.Sprintf("Elastic:%s unauthorized", url)
				result.Result = true
				return err, result
			}
		}

		res2, err := http.NewRequest("GET", url, nil)
		if err == nil {
			res2.Header.Add("Authorization", auth)
			res2.Header.Add("User-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Safari/537.36")
			res2.Header.Add("Accept", "*/*")
			res2.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
			res2.Header.Add("Connection", "close")
			resp2, err := client.Do(res2)
			if err == nil {
				defer resp2.Body.Close()
				body2, _ := ioutil.ReadAll(resp2.Body)
				if strings.Contains(string(body2), "/_cat/master") {
					//AnnonTest := fmt.Sprintf("Elastic:%s unauthorized", url)
					result.Result = true
				}
			}
		}

	}

	return err, result
}
