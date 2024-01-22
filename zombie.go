//go:generate go run templates/templates_gen.go -t templates -o pkg/templates.go -need zombie
package main

import (
	"github.com/chainreactors/zombie/cmd"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	cmd.Zombie()
}
