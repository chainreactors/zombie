package core

import (
	"github.com/chainreactors/parsers"
	"io/ioutil"
	"math/rand"
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
