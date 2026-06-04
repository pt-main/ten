package main

import (
	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

func (l *Language) config_parse(e *system.Engine, pn parsing.ParsedNode) error {
	_data := pn.Parsed[0]
	_data = _data[len("((?") : len(_data)-len("))")]
	p := parsing.Parser2{}
	res, err := p.Parse(_data)
	if err != nil {
		return err
	}
	for _, parsed := range res {
		cmd := parsed.Metadata["command"].(string)
		args := parsed.Metadata["args"].(string)
		l.config[cmd] = l.replaceAll(args)
	}
	return nil
}
