package plugin

import "github.com/chainreactors/zombie/pkg"

type Plugin interface {
	Name() string
	Open(task *pkg.Task) (pkg.Session, error)
	Unauth(task *pkg.Task) (pkg.Session, error)
}
