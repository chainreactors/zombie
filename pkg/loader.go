package pkg

import (
	"encoding/json"
	templates "github.com/chainreactors/neutron/templates_gogo"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/words/mask"
	"strings"
)

var (
	Rules       map[string]string              = make(map[string]string)
	Keywords    map[string][]string            = make(map[string][]string)
	TemplateMap map[string]*templates.Template = make(map[string]*templates.Template)
)

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

func LoadRules() error {
	var err error
	var data map[string]interface{}
	err = json.Unmarshal(LoadConfig("zombie_rule"), &data)
	if err != nil {
		return err
	}
	for k, v := range data {
		Rules[k] = v.(string)
	}
	return nil
}

func LoadTemplates() error {
	var err error
	var content []byte
	content = LoadConfig("zombie_template")
	var t []*templates.Template
	err = json.Unmarshal(content, &t)
	if err != nil {
		return err
	}
	for _, template := range t {
		for _, tag := range strings.Split(template.Info.Tags, ",") {
			if name, ok := Services.GetName(Service(tag)); ok {
				err = template.Compile(nil)
				if err != nil {
					iutils.Fatal("" + err.Error())
				}
				TemplateMap[name] = template
				break
			}
		}
	}
	return nil
}

func Load() error {
	var err error
	err = LoadKeyword()
	if err != nil {
		return err
	}

	err = LoadRules()
	if err != nil {
		return err
	}

	err = LoadTemplates()
	if err != nil {
		return err
	}
	return err
}
