package pkg

var DefaultStatistor Statistor

type Statistor struct {
	Count int
	Cur   string
}
