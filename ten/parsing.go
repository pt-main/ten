package ten

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func (l *Language) parseLogic(data []string) (bool, []string, error) {
	arg := ""
	not_enable := false
	has_enable := false
	var_enable := false
	if len(data) >= 1 {
		if data[0] == "not" {
			not_enable = true
			data = data[1:]
		}
		if data[0] == "logic" {
			var_enable = true
			data = data[1:]
		}
		if data[0] == "has" && !var_enable {
			has_enable = true
			data = data[1:]
		}
		arg = strings.TrimSpace(data[0])
		data = data[1:]
	} else {
		return false, nil, fmt.Errorf("Invalid placeholder length: %d (min is 1)", len(data))
	}
	var cond bool
	if has_enable {
		keys := make([]string, 0, len(l.templateSource))
		for k := range l.templateSource {
			keys = append(keys, k)
		}
		cond = slices.Contains(keys, arg)
	} else if var_enable {
		var ok bool
		cond, ok = l.logic_scope[arg]
		if !ok {
			return false, nil, errors.New("Can't find arg " + arg + " in logic scope.")
		}
	} else {
		cond = slices.Contains(l.flags, arg)
	}
	if not_enable {
		cond = !cond
	}
	return cond, data, nil
}

func (l *Language) replaceAll(data string) string {
	replaces := [][]string{{"\\s", " "}, {"\\n", "\n"}}
	for _, replace := range replaces {
		data = strings.ReplaceAll(data, replace[0], replace[1])
	}
	return data
}

func (l *Language) processValue(what string) string {
	pre := ""
	post := ""
	data := ""
	_what := what
	pre_plus := l.config["pre_plus"]
	post_plus := l.config["post_plus"]
	raw := l.config["raw"]
	if strings.HasPrefix(strings.ToLower(what), raw) {
		data = l.replaceAll(what[len(raw):])
	} else {
		what = strings.ReplaceAll(what, pre_plus, "")
		what = strings.ReplaceAll(what, post_plus, "")
		keys := make([]string, 0, len(l.templateSource))
		for k := range l.templateSource {
			keys = append(keys, k)
		}
		if slices.Contains(keys, what) {
			if strings.Contains(strings.ToLower(_what), pre_plus) {
				for range strings.Count(_what, pre_plus) {
					pre += " "
				}
			}
			if strings.Contains(strings.ToLower(_what), post_plus) {
				for range strings.Count(_what, post_plus) {
					post += " "
				}
			}
			data = l.templateSource[what]
		}
	}
	return pre + data + post
}
