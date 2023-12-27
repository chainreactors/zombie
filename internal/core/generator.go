package core

import (
	"errors"
	"github.com/chainreactors/files"
	"github.com/chainreactors/words"
	"github.com/chainreactors/zombie/pkg"
	"io/ioutil"
	"os"
	"strings"
)

func NewNullGenerator() *Generator {
	g := &Generator{
		C: make(chan string),
	}
	return g
}

func NewGeneratorWithInput(in []string) *Generator {
	g := &Generator{
		C: make(chan string),
	}

	g.Word = words.NewWorder(in)
	return g
}

func NewGeneratorWithFile(filename string) (*Generator, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	g := &Generator{
		C:        make(chan string),
		Filename: filename,
	}

	g.Word = words.NewWorderWithFile(f)
	return g, nil
}

func NewGeneratorWithWord(word string, params [][]string, keywords map[string][]string) (*Generator, error) {
	g := &Generator{}

	w, err := words.NewWorderWithDsl(word, params, keywords)
	if err != nil {
		return nil, err
	}

	g.Word = w
	return g, nil
}

type Generator struct {
	Filename     string
	File         *os.File
	WordString   string
	Word         *words.Worder
	RuleFilename string
	Rules        string
	Filter       string
	Fns          []func(string) string
	C            chan string
	cache        []string
	running      bool
}

func (g *Generator) initWord() {
	if g.Fns != nil {
		g.Word.Fns = g.Fns
	}

	if g.Rules != "" {
		g.Word.SetRules(g.Rules, g.Filter)
	}
}

func (g *Generator) Run() {
	g.running = true
	g.initWord()
	g.Word.Run()
	g.C = g.Word.C
}

func (g *Generator) RunAsSlice() []string {
	if !g.running {
		g.Run()
		g.cache = g.All()
	}
	return g.cache
}

func (g *Generator) SetFile(filename string) error {
	g.Filename = filename
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	g.File = f
	return nil
}

func (g *Generator) SetRuleFile(filename string) error {
	if files.IsExist(filename) {
		g.RuleFilename = filename
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		g.Rules = string(content)
	} else {
		if content, ok := pkg.Rules[filename]; ok {
			g.Rules = content
			return nil
		} else {
			return errors.New("not found file and not found preset rule," + filename)
		}
	}

	return nil
}

func (g *Generator) SetInternalRule(rulename string) error {
	if content, ok := pkg.Rules[rulename]; ok {
		g.Rules = content
		return nil
	}
	return errors.New("rule not found")
}

func (g *Generator) SetFilter(filter []string) {
	g.Filter = strings.Join(filter, "\n")
}

func (g *Generator) AddFunc(fun func(string) string) {
	g.Fns = append(g.Fns, fun)
}

func (g *Generator) AddFuncs(funs []func(string) string) {
	g.Fns = append(g.Fns, funs...)
}

func (g *Generator) All() []string {
	var l []string
	for i := range g.C {
		l = append(l, i)
	}
	return l
}
