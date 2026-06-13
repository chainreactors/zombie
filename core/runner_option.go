package core

import "github.com/chainreactors/zombie/pkg"

type RunnerOption struct {
	Threads         int
	Concurrency     int // 单 host 并发上限(0=不限),见 hostLimiter
	Timeout         int
	Top             int
	Mod             string // clusterbomb / pitchfork / sniper
	FirstOnly       bool
	NoUnAuth        bool
	NoCheckHoneyPot bool
	Strict          bool
	Raw             bool
	Quiet           bool
	// ManualDrain：SDK 调用方若自行消费 OutputCh，则置 true 以禁用内建
	// OutputHandler；CLI/默认为 false，由 OutputHandler 负责消费并打印结果。
	// 注意：OutputHandler 是无缓冲 OutputCh 的唯一内建读者，关掉它而又无人
	// drain 会让首个 Output() 永久阻塞——故默认必须启动。
	ManualDrain bool

	// ProxyDial 非 nil 时透传到每个 Task，使插件通过代理建立连接。
	ProxyDial pkg.DialFunc
}

var DefaultRunnerOption = &RunnerOption{
	Threads:     100,
	Concurrency: 8,
	Timeout:     5,
	Mod:         ModBomb,
	FirstOnly:   true,
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
