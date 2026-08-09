package action

import (
	"github.com/chainreactors/logs"
	"github.com/chainreactors/zombie/pkg"
)

type AuditAction struct {
	Patterns []string
	Limit    int
}

func NewAuditAction() *AuditAction {
	return &AuditAction{
		Patterns: pkg.DefaultAuditPatterns,
		Limit:    100,
	}
}

func (a *AuditAction) Name() string { return "audit" }

func (a *AuditAction) Run(session pkg.Session, task *pkg.Task) (*pkg.ActionResult, error) {
	auditable, ok := session.(pkg.AuditableSession)
	if !ok {
		return nil, nil
	}

	data, err := auditable.Audit(a.Patterns, a.Limit)
	if err != nil {
		logs.Log.Debugf("[audit] %s: %v", task.URI(), err)
		return nil, nil
	}

	result := &pkg.ActionResult{
		Loot: make(map[string][]byte),
	}
	for location, samples := range data {
		result.Loot[location] = []byte(samples)
	}
	return result, nil
}
