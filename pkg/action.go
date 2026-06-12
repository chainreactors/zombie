package pkg

import (
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/parsers"
)

type ActionResult struct {
	Extracteds parsers.Extracteds
	Vulns      common.Vulns
	Loot       map[string][]byte
}
type Action interface {
	Name() string
	Run(session Session, task *Task) (*ActionResult, error)
}
