package core

import (
	"bytes"
	"github.com/chainreactors/words/rule"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"strings"
)

var (
	dictCache     = make(map[string][]string)
	wordlistCache = make(map[string][]string)
	ruleCache     = make(map[string][]rule.Expression)
)

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

func loadFileAndCombine(filename []string) (string, error) {
	var bs bytes.Buffer
	for _, f := range filename {
		if data, ok := pkg.Rules[f]; ok {
			bs.WriteString(strings.TrimSpace(data))
			bs.WriteString("\n")
		} else {
			content, err := ioutil.ReadFile(f)
			if err != nil {
				return "", err
			}
			bs.Write(bytes.TrimSpace(content))
			bs.WriteString("\n")
		}
	}
	return bs.String(), nil
}

func loadFileWithCache(filename string) ([]string, error) {
	if dict, ok := dictCache[filename]; ok {
		return dict, nil
	}
	dict, err := loadFileToSlice(filename)
	if err != nil {
		return nil, err
	}
	dictCache[filename] = dict
	return dict, nil
}

func loadDictionaries(filenames []string) ([][]string, error) {
	dicts := make([][]string, len(filenames))
	for i, name := range filenames {
		dict, err := loadFileWithCache(name)
		if err != nil {
			return nil, err
		}
		dicts[i] = dict
	}
	return dicts, nil
}
