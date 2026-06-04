package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pt-main/lc/parsing"
	"github.com/pt-main/lc/system"
)

func (l *Language) code(e *system.Engine, pn parsing.ParsedNode) error {
	_data := strings.TrimSpace(pn.Parsed[0])
	_data = _data[len("{{?") : len(_data)-len("}}")]
	if strings.TrimSpace(_data) == "" {
		return nil
	}
	p := parsing.Parser2{}
	res, err := p.Parse(_data)
	if err != nil {
		return err
	}
	write_to := ""
	err_text := ""
	raw := ""
	idx := 0
	var parsed parsing.ParsedNode
	for idx, parsed = range res {
		cmd := parsed.Metadata["command"].(string)
		args := parsed.Metadata["args"].(string)
		raw = parsed.Metadata["__raw"].(string)
		argsplit := strings.Split(args, " ")
		// MARK: writing
		if write_to != "" {
			raw := parsed.Metadata["__raw"].(string)
			if !strings.HasSuffix(raw, l.config["string_end"]) {
				l.templateSource[write_to] += "\n" + raw
			} else {
				l.templateSource[write_to] += "\n" + raw[:len(raw)-len(l.config["string_end"])]
				write_to = ""
			}
			// MARK: comment and logic
		} else if strings.HasPrefix(cmd, l.config["comment"]) {
		} else if cmd == "logic" {
			set := -1
			for idx, block := range argsplit {
				if block == "=" {
					set = idx
				}
			}
			if set == -1 {
				err_text = "has no '=' symbol or incorrect position"
				goto err_label
			}
			res, _, err := l.parseLogic(argsplit)
			if err != nil {
				return err
			}
			l.logic_scope[argsplit[set+1]] = res
			// MARK: set
		} else if cmd == "set" {
			if argsplit[1] != "=" {
				err_text = "has no '=' symbol or incorrect position"
				goto err_label
			}
			param := argsplit[0]
			if strings.HasPrefix(argsplit[2], l.config["string_start"]) {
				write_to = param
				l.templateSource[param] = strings.Join(argsplit[2:], " ")[len(l.config["string_start"]):]
			} else {
				l.templateSource[param] = strings.Join(argsplit[2:], " ")
			}
			// MARK: bool
		} else if cmd == "bool" {
			if len(argsplit) != 5 {
				err_text = "invalid length"
				goto err_label
			}
			arg1, ok1 := l.logic_scope[argsplit[0]]
			arg2, ok2 := l.logic_scope[argsplit[2]]
			if !ok1 || !ok2 {
				return errors.New("Can't find argument in logic scope")
			}
			if argsplit[3] != "=" {
				err_text = "has no '=' symbol or incorrect position"
				goto err_label
			}
			op := argsplit[1]
			if !slices.Contains([]string{"or", "and"}, op) {
				err_text = "invalid instruction"
				goto err_label
			}
			var res bool
			switch op {
			case "or":
				res = arg1 || arg2
			case "and":
				res = arg1 && arg2
			}
			l.logic_scope[argsplit[4]] = res
			// MARK: place
		} else if cmd == "place" {
			set := -1
			for idx, val := range argsplit {
				if val == "=" {
					set = idx
				}
			}
			if set == -1 {
				err_text = "has no '=' symbol or incorrect position"
				goto err_label
			}
			placeholder, err := l.placeholderEval(strings.Join(argsplit[:set], " "))
			if err != nil {
				return err
			}
			name := argsplit[set+1]
			l.templateSource[name] = l.processValue(placeholder)
			// MARK: printing
		} else if cmd == "print" {
			fmt.Print(l.replaceAll(strings.TrimSpace(args)))
		} else if cmd == "printv" {
			arg := ""
			if len(argsplit) > 1 {
				arg = argsplit[1]
			}
			if argsplit[0] == "logic" {
				if arg != "" {
					fmt.Print(l.logic_scope[arg])
				} else {
					fmt.Print(l.logic_scope)
				}
			} else if argsplit[0] == "ts" {
				if arg != "" {
					fmt.Print(l.templateSource[arg])
				} else {
					fmt.Print(l.templateSource)
				}
			} else {
				err_text = "invalid printv type"
				goto err_label
			}
			// MARK: compare
		} else if cmd == "compare" {
			arg1 := l.processValue(argsplit[0])
			arg2 := l.processValue(argsplit[2])
			if argsplit[3] != "=" {
				err_text = "has no '=' symbol or incorrect position"
				goto err_label
			}
			res := argsplit[4]
			switch argsplit[1] {
			case "==", "equals":
				l.logic_scope[res] = arg1 == arg2
			case ".>", "startswith":
				l.logic_scope[res] = strings.HasPrefix(arg1, arg2)
			case ">.", "endswith":
				l.logic_scope[res] = strings.HasSuffix(arg1, arg2)
			default:
				err_text = "invalid instruction"
				goto err_label
			}
			// MARK: add
		} else if cmd == "add" {
			set := -1
			for idx, val := range argsplit {
				if val == "=" {
					set = idx
				}
			}
			if set == -1 {
				goto err_label
			}
			res := argsplit[set+1]
			_args := argsplit[:set]
			for _, arg := range _args {
				l.templateSource[res] += l.processValue(arg)
			}
			// MARK: make
		} else if cmd == "make" {
			op := argsplit[0]
			arg := l.processValue(argsplit[1])
			if argsplit[2] != "=" {
				err_text = "has no '=' symbol or incorrect position"
				goto err_label
			}
			res := argsplit[3]
			if op == "upper" {
				l.templateSource[res] = strings.ToUpper(arg)
			} else if op == "lower" {
				l.templateSource[res] = strings.ToLower(arg)
			} else if strings.HasPrefix(op, "slice_e:") || strings.HasPrefix(op, "slice_s:") {
				_from := op[len("slice_?:"):]
				from, err := strconv.Atoi(_from)
				if err != nil {
					return err
				}
				runes := []rune(arg)
				length := len(runes)
				if from < 0 {
					from = length + from
				}
				if from < 0 {
					from = 0
				}
				if from > length {
					from = length
				}
				var _res string
				if strings.HasPrefix(op, "slice_s:") {
					_res = string(runes[from:])
				} else {
					_res = string(runes[:from])
				}
				l.templateSource[res] = _res
			} else {
				err_text = "invalid operation"
				goto err_label
			}
		} else if cmd != "" {
			goto err_label
		}
	}
	return nil
err_label:
	text := "Invalid code block"
	if err_text != "" {
		text += ", error: " + err_text
	}
	if raw != "" {
		text += " \n(at '" + strings.TrimSpace(raw) + ", block line " + strconv.Itoa(idx) + ")"
	}
	return errors.New(text)
}
