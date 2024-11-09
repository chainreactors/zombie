package pkg

import (
	"github.com/chainreactors/fingers/fingers"
	"github.com/chainreactors/fingers/resources"
	templates "github.com/chainreactors/neutron/templates"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/iutils"
	"github.com/chainreactors/words/mask"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	Rules         map[string]string              = make(map[string]string)
	Keywords      map[string][]string            = make(map[string][]string)
	TemplateMap   map[string]*templates.Template = make(map[string]*templates.Template)
	FingersEngine *fingers.FingersEngine
)

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

	err = LoadFingers()
	if err != nil {
		return err
	}
	return err
}

func LoadKeyword() error {
	// load mask
	var err error
	var commonKeyword map[string]interface{}
	err = yaml.Unmarshal(LoadConfig("zombie_common"), &commonKeyword)
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
	err = yaml.Unmarshal(LoadConfig("zombie_default"), &defaultKeyword)
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
	err = yaml.Unmarshal(LoadConfig("zombie_rule"), &data)
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
	err = yaml.Unmarshal(content, &t)
	if err != nil {
		return err
	}
	for _, template := range t {
		if template.Info.Zombie == "" {
			continue
		}
		Services.Register(&Service{Name: template.Info.Zombie, Source: NeutronSource})
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
	err = yaml.Unmarshal(LoadConfig("port"), &ports)
	if err != nil {
		return err
	}

	utils.PrePort = utils.NewPortPreset(ports)
	return nil
}

func LoadFingers() error {
	resources.FingersHTTPData = LoadConfig("http")
	resources.FingersSocketData = LoadConfig("socket")
	engine, err := fingers.NewFingersEngine()
	if err != nil {
		return err
	}
	FingersEngine = engine
	return nil
}
