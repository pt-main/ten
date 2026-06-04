package ten

import (
	"github.com/dlclark/regexp2"
	"github.com/pt-main/lc"
	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

// Version of the template engine.
const Version = "0.9.1"

// NewTemplateEngine creates and initializes a new system.Engine instance
// configured with the lexer rules for config, code, placeholder, and raw blocks.
// The engine is ready to process template text after registering commands.
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

// NewTemplate creates a new Language instance and an associated system.Engine.
// The engine is pre‑registered with commands for "placeholder", "code", "config",
// and "raw". Returns the Language instance, the Engine, or an error if language creation fails.
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

// TemplateReplace processes the given template string using the specified flags,
// placeholders, and configuration. It returns the final rendered string or an error.
//
// Parameters:
//   - template: raw template text containing special blocks:
//     ((?CONFIG ... )) – configuration overrides,
//     {{?CODE ... }}   – code/logic instructions,
//     [[? ... ]]       – placeholders.
//   - flags: list of string flags that can be tested in conditions.
//   - placeholders: initial key‑value map for template source variables.
//   - config: configuration overrides for the language behaviour (e.g. comment markers,
//     string delimiters). May be nil to use defaults.
//
// Returns:
//   - The fully processed template as a string.
//   - Non‑nil error if parsing, logic evaluation, or processing fails.
//
// Example:
//
//	result, err := TemplateReplace(
//	    "Hello [[? name if has name else World]]!",
//	    []string{},
//	    map[string]string{"name": "Alice"},
//	    nil,
//	)
//	// result == "Hello Alice!"
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
