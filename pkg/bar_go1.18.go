//go:build nobar && !go1.18
// +build nobar,!go1.18

package pkg

func InitBar() {

}

func NewBar(u string, total int, stat *Statistor) *Bar {
	return &Bar{}
}

type Bar struct {
}

func (bar *Bar) Done() {
}

func (bar *Bar) Close() {
}
