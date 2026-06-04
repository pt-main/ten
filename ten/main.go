package ten

import (
	"github.com/dlclark/regexp2"
	"github.com/pt-main/lc"
	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

const Version = "0.8.2"

// Create engine for trmplate engine
func NewTemplateEngine() *system.Engine {
	p := parsing.NewLexer([]parsing.LexerRule{
		{
			Type:    "config",
			Pattern: regexp2.MustCompile(_CONFIG_BLOCK_SYNTAX, 0),
		},
		{
			Type:    "code",
			Pattern: regexp2.MustCompile(_CODE_BLOCK_SYNTAX, 0),
		},
		{
			Type:    "placeholder",
			Pattern: regexp2.MustCompile(_PLACEHOLDER_SYNTAX, 0),
		},
		{
			Type:    "raw",
			Pattern: regexp2.MustCompile(_RAW_SYNTAX, 0),
		},
	})
	e := lc.NewEngine(system.StringResType, []string{"main"}, true, p)
	return e
}

// New Template engine.
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
