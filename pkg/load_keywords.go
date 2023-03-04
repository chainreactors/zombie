package pkg

import (
	"encoding/json"
	"github.com/chainreactors/parsers/iutils"
	"github.com/chainreactors/words/mask"
)

func LoadKeyword() error {
	// load mask
	var err error
	var keywords map[string]interface{}
	err = json.Unmarshal(LoadConfig("zombie"), &keywords)
	if err != nil {
		return err
	}

	for k, v := range keywords {
		t := make([]string, len(v.([]interface{})))
		for i, vv := range v.([]interface{}) {
			t[i] = iutils.ToString(vv)
		}
		mask.SpecialWords[k] = t
	}
	return nil
}
