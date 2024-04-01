package pkg

import (
	"fmt"
	"strings"
)

type Statistor struct {
	Total   int
	Success int
	Cur     string
	Tasks   map[string]int
}

func (stat *Statistor) TaskString() string {
	var s strings.Builder
	for k, v := range stat.Tasks {
		s.WriteString(fmt.Sprintf("%s:%d ", k, v))
	}
	return s.String()
}
