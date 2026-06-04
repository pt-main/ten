package ten

import (
	"errors"
	"strconv"
	"strings"

	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

func (l *Language) placeholderIfGenerator(data []string) (string, error) {
	else_phd := ""
	var cond bool
	var _data []string
	if data[1] != "if" {
		goto invalid
	} else {
		cond, _data, _ = l.parseLogic(data[2:])
	}
	if len(_data) != 0 {
		if len(_data) == 2 {
			if _data[0] != "else" {
				goto invalid
			}
			else_phd = _data[1]
		} else {
			goto invalid
		}
	}
	if cond {
		return data[0], nil
	} else {
		return else_phd, nil
	}
invalid:
	return "", errors.New("Invalid placeholder")
}

func (l *Language) placeholderIf(data []string) (string, error) {
	if l.Config["ALPHA"] != "" {
		result := ""
		_res := data
		for len(_res) > 1 {
			start := 0
			for idx, val := range _res {
				if val == "if" {
					start = idx - 1
				}
			}
			_res = _res[start:]
			res, err := l.placeholderIfGenerator(_res)
			if err != nil {
				return "", err
			}
			_res = strings.Split(res, " ")
		}
		result = _res[0]
		return result, nil
	} else {
		res, err := l.placeholderIfGenerator(data)
		if err != nil {
			return "", err
		}
		return res, nil
	}
}

func (l *Language) placeholderEval(_data string) (string, error) {
	data := strings.Split(strings.TrimSpace(_data), " ")
	if len(strings.TrimSpace(_data)) == 0 {
	} else if len(data) == 1 {
		return data[0], nil
	} else {
		switch data[1] {
		case "if":
			return l.placeholderIf(data)
		case "for":
			iters, err := strconv.Atoi(data[2])
			if err != nil {
				return "", err
			}
			result := ""
			for i := 0; i < iters; i++ {
				result += l.processValue(data[0])
			}
			return l.Config["raw"] + result, nil
		default:
			goto invalid
		}
	}
	return "", errors.New("Empty")
invalid:
	return "", errors.New("Invalid placeholder")
}

func (l *Language) placeholder(e *system.Engine, pn parsing.ParsedNode) error {
	_data := pn.Parsed[0]
	_data = _data[len(_PLACEHOLDER_START_END[0]) : len(_data)-len(_PLACEHOLDER_START_END[1])]
	str, err := l.placeholderEval(_data)
	if err != nil {
		return err
	}
	return e.Generator.AddString(l.processValue(str), "main")
}
