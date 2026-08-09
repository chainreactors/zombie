package plugin

import "github.com/chainreactors/zombie/pkg"

type Plugin interface {
	Open(*pkg.Task) (pkg.Session, error)
}

type UnauthPlugin interface {
	Unauth(*pkg.Task) (pkg.Session, error)
}

type Service = pkg.Service
