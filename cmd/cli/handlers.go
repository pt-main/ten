package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pt-main/tap"
	"github.com/pt-main/ten/ten"
)

func process_handler(parser *tap.Parser, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Invalid argument length: %v (must be 1)", len(args))
	}
	_file, err := os.ReadFile(args[0])
	if err != nil {
		return errors.New("Readng error: " + err.Error())
	}
	file := string(_file)
	keys := make([]string, 0, len(parser.Flags))
	for k := range parser.Flags {
		keys = append(keys, k)
	}
	flags := []string{"__inCli"}
	placeholders := map[string]string{}
	systems := map[string]string{}
	for _, key := range keys {
		if parser.Flags[key] == "" && strings.ToLower(key) == key {
			flags = append(flags, key)
		} else if strings.ToLower(key) == key {
			placeholders[key] = parser.Flags[key]
		} else {
			systems[key] = parser.Flags[key]
		}
	}
	out := ""
	val1, ok1 := systems["O"]
	val2, ok2 := systems["OUT"]
	if ok1 || ok2 {
		flags = append(flags, "__hasOut")
	}
	result, err := ten.TemplateReplace(file, flags, placeholders, nil)
	if err != nil {
		return err
	}
	if ok1 {
		out = val1
	} else if ok2 {
		out = val2
	} else {
		print(result)
		return nil
	}
	err = os.WriteFile(out, []byte(result), 0644)
	if err != nil {
		fmt.Println("Writing error:", err)
	}
	return nil
}

func docs_handler(parser *tap.Parser, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Invalid argument length: %v (must be 1)", len(args))
	}
	_file, err := os.ReadFile(args[0])
	if err != nil {
		return errors.New("Readng error: " + err.Error())
	}
	file := string(_file)
	lang, engine, err := ten.NewTemplate(nil, nil, nil)
	if err != nil {
		return err
	}
	err = engine.Process(file)
	if err != nil {
		return err
	}
	println(lang.config["documentation"])
	return nil
}
