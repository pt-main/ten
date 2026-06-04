package main

import (
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/pt-main/lc"
	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

func NewTemplateEngine() *system.Engine {
	p := parsing.NewLexer([]parsing.LexerRule{
		{
			Type:    "config",
			Pattern: regexp2.MustCompile(`\(\(\?\s*([\s\S]*?)\s*\)\)`, 0),
		},
		{
			Type:    "code",
			Pattern: regexp2.MustCompile(`\{\{\?\s*([\s\S]*?)\s*\}\}`, 0),
		},
		{
			Type:    "placeholder",
			Pattern: regexp2.MustCompile(`\[\[\?\s\n*(.*?)\n*\s\]\]`, 0),
		},
		{
			Type:    "raw",
			Pattern: regexp2.MustCompile(`(?s)(.+?)(?=\[\[\?|\{\{\?|\(\(\?|$)`, 0),
		},
	})
	e := lc.NewEngine(system.StringResType, []string{"main"}, true, p)
	return e
}

func NewTemplate(
	flags []string,
	placeholders map[string]string,
	config map[string]string,
) (*Language, *system.Engine, error) {
	lang, err1 := NewLanguage(flags, placeholders, config)
	if err1 != nil {
		return nil, nil, err1
	}
	engine := NewTemplateEngine()
	engine.NewCommand("placeholder", lang.placeholder, "")
	engine.NewCommand("code", lang.code, "")
	engine.NewCommand("config", lang.config_parse, "")
	engine.NewCommand("raw", func(e *system.Engine, pn parsing.ParsedNode) error {
		return e.Generator.AddString(pn.Metadata["__raw"].(string), "main")
	}, "")
	return lang, engine, nil
}

func TemplateReplace(
	template string,
	flags []string,
	placeholders map[string]string,
	config map[string]string,
) (string, error) {
	_, engine, err1 := NewTemplate(flags, placeholders, config)
	if err1 != nil {
		return "", err1
	}
	err2 := engine.Process(template)
	if err2 != nil {
		return "", err2
	}
	res, _ := engine.Generator.GetStringRes("")
	return res, nil
}

func test() {
	res, _ := TemplateReplace(`((?
	pre_plus /
	raw >
)){{?
	set star = (⭐)
	// logics
	compare name == >Pt = isPt
	compare name == >Tim = isTim
	bool isPt or isTim = isCreator
	// placing
	place >evening if evening else >day = bl1
	place name if has name else >User = bl2
	place /star if vip = bl3
	place >\s( if var isCreator = bl4
	place >hello,\screator if var isCreator = bl5
	place >? if var isTim = bl6
	place >! if var isPt = bl7 
	place >) if var isCreator = bl8
	place >\n if has message = bl9
	// summing
	add >Good\s bl1 >,\s bl2 bl3 >! bl4 bl5 bl6 bl7 bl8 bl9 message = blocks
	// debug
	printv ts
	print \n
	printv logic
	print \n
}}[? sep for 40 ]
[? blocks ]
[? sep for 40 ]`,
		[]string{"evening", "vip"}, map[string]string{"name": "Pt", "message": "Here is test message", "sep": "="}, nil)
	fmt.Println(res) // ...\nGood evening, Tim (⭐)! (hello, creator!)\ntest message\n...
}

func main() {
	cli()
}
