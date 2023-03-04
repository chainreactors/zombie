//go:generate go run templates/templates_gen.go -t templates -o pkg/templates.go -need zombie
package main

import "github.com/chainreactors/zombie/cmd"

func main() {
	cmd.Zombie()
}
