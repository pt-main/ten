package ten

import (
	"strings"

	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

func (l *Language) config_parse(e *system.Engine, pn parsing.ParsedNode) error {
	_data := pn.Parsed[0]
	_data = _data[len(_CONFIG_START_END[0]) : len(_data)-len(_CONFIG_START_END[1])]
	p := parsing.Parser2{}
	res, err := p.Parse(_data)
	if err != nil {
		return err
	}
	for _, parsed := range res {
		if strings.TrimSpace(parsed.Metadata["__raw"].(string)) != "" {
			cmd := parsed.Metadata["command"].(string)
			args := parsed.Metadata["args"].(string)
			l.config[cmd] = l.replaceAll(args)
		}
	}
	return nil
}
