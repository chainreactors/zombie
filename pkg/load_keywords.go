package pkg

import (
	"encoding/json"
	"github.com/chainreactors/parsers/iutils"
	"github.com/chainreactors/words/mask"
)

var Keywords map[string][]string = make(map[string][]string)

func LoadKeyword() error {
	// load mask
	var err error
	var commonKeyword map[string]interface{}
	err = json.Unmarshal(LoadConfig("zombie_common"), &commonKeyword)
	if err != nil {
		return err
	}

	for k, v := range commonKeyword {
		t := make([]string, len(v.([]interface{})))
		for i, vv := range v.([]interface{}) {
			t[i] = iutils.ToString(vv)
		}
		Keywords[k] = t
	}

	var defaultKeyword map[string]interface{}
	err = json.Unmarshal(LoadConfig("zombie_default"), &defaultKeyword)
	if err != nil {
		return err
	}
	for k, v := range defaultKeyword {
		var tmplist []string
		for _, i := range v.([]interface{}) {
			if i == "{{common_pwd}}" {
				tmplist = append(tmplist, Keywords["common_pwd"]...)
			} else if i == "{{blank}}" {
				tmplist = append(tmplist, "")
			} else {
				tmplist = append(tmplist, iutils.ToString(i))
			}
		}
		Keywords[k] = tmplist
	}
	mask.SpecialWords = Keywords
	return nil
}
