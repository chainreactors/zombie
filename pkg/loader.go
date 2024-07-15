package pkg

import (
	"encoding/json"
	templates "github.com/chainreactors/neutron/templates"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils"
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
		if template.Info.Zombie == "" {
			continue
		}
		Services[template.Info.Zombie] = &Service{Name: template.Info.Name, Source: "template"}
		err := template.Compile(nil)
		if err != nil {
			return err
		}
		TemplateMap[template.Info.Zombie] = template

		// load gogo_finger-zombie-service map
		if template.Info.Zombie != "" {
			for _, tag := range template.GetTags() {
				parsers.ZombieMap[strings.ToLower(tag)] = template.Info.Zombie
			}
			for _, finger := range template.Fingers {
				parsers.ZombieMap[finger] = template.Info.Zombie
			}
		}
	}
	parsers.RegisterZombieServiceAlias()
	return nil
}

func LoadPorts() error {
	var ports []*utils.PortConfig
	var err error
	err = json.Unmarshal(LoadConfig("port"), &ports)
	if err != nil {
		return err
	}

	utils.PrePort = utils.NewPortPreset(ports)
	return nil
}

func Load() error {
	var err error
	err = LoadPorts()
	if err != nil {
		return err
	}

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
