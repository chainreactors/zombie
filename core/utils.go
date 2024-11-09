package core

import (
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"math/rand"
	"net"
	"net/url"
	"strings"
)

func LoadGogoFile(filename string) ([]*Target, error) {
	gd, err := parsers.ParseGogoData(filename)
	if err != nil {
		return nil, err
	}
	ptargets := gd.ToZombie()
	var targets []*Target
	for _, t := range ptargets {
		targets = append(targets, &Target{
			IP:      t.IP,
			Port:    t.Port,
			Service: t.Service,
			Scheme:  t.Scheme,
			Param:   t.Param,
		})
	}
	return targets, nil
}

func loadFileToSlice(filename string) ([]string, error) {
	var ss []string
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ss = strings.Split(strings.TrimSpace(string(content)), "\n")

	// 统一windows与linux的回车换行差异
	for i, word := range ss {
		ss[i] = strings.TrimSpace(word)
	}

	return ss, nil
}

func parseAuthPair(auth string) (string, string) {
	pair := strings.Split(auth, "::")
	switch len(pair) {
	case 1:
		return auth, ""
	case 2:
		return pair[0], pair[1]
	default:
		return pair[0], auth[strings.Index(auth, "::")+2:]
	}
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func ParseUrl(u string) (*Target, bool) {
	var t *Target
	parsed, err := url.Parse(u)
	if err != nil && parsed != nil {
		if ip, port, err := net.SplitHostPort(u); err == nil {
			t = &Target{
				IP:   ip,
				Port: port,
			}
		} else {
			return nil, false
		}

	} else if parsed == nil {
		return nil, false
	}

	if parsed.Host == "" {
		if utils.IsIp(u) {
			return &Target{
				IP: u,
			}, true
		}
		return nil, false
	}

	t = &Target{
		IP: parsed.Hostname(),
	}
	if parsed.Port() != "" {
		t.Port = parsed.Port()
	}

	if parsed.Scheme != "" {
		t.Service = parsed.Scheme
		if t.Port == "" {
			t.Port = pkg.Services.DefaultPort(t.Service)
		}
		t.Scheme = parsed.Scheme
	} else if t.Port != "" {
		t.Service = pkg.GetDefault(t.Port)
	}
	if parsed.User != nil {
		if parsed.User.Username() != "" {
			t.Username = parsed.User.Username()
		}
		if pwd, _ := parsed.User.Password(); pwd != "" {
			t.Password, _ = parsed.User.Password()
		}
	}

	return t, true
}

func SimpleParseUrl(u string) *Target {
	result := strings.Split(u, ":")
	if len(result) == 1 {
		return &Target{
			IP: result[0],
		}
	} else {
		return &Target{
			IP:   result[0],
			Port: result[1],
		}
	}
}

func transformChan(ipch chan *utils.IP) chan string {
	ch := make(chan string)
	go func() {
		for i := range ipch {
			ch <- i.String()
		}
		close(ch)
	}()
	return ch
}
