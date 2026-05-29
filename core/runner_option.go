package core

import "github.com/chainreactors/zombie/pkg"

type RunnerOption struct {
	Threads         int
	Timeout         int
	Top             int
	Mod             string // clusterbomb / pitchfork / sniper
	FirstOnly       bool
	NoUnAuth        bool
	NoCheckHoneyPot bool
	Strict          bool
	Raw             bool
	Quiet           bool

	// ProxyDial 非 nil 时透传到每个 Task，使插件通过代理建立连接。
	ProxyDial pkg.DialFunc
}

var DefaultRunnerOption = &RunnerOption{
	Threads:   100,
	Timeout:   5,
	Mod:       ModBomb,
	FirstOnly: true,
}

func NewDefaultRunnerOption() *RunnerOption {
	return DefaultRunnerOption.Clone()
}

func (o *RunnerOption) Clone() *RunnerOption {
	if o == nil {
		return NewDefaultRunnerOption()
	}
	clone := *o
	return &clone
}
