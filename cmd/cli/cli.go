package main

import (
	"github.com/pt-main/tap"
)

func cli() error {
	p := tap.NewParser("ten", `[?MA]╭───────[?RT] [?BE]Ten cli[?RT]
[?MA]│[?RT]    Ten - Template Engine.
[?MA]│[?RT]    Ten cli [[?YW]v`+Version+`[?RT]]
[?MA]│[?RT]    Only [?YW]humanmade[?RT], By [?YW]Pt[?RT].
[?MA]╰───────[?RT]`, []string{"h", "help"},
		tap.NewParserConfig("", "", "", "", "", ""))
	p.AddCommand("process", process_handler, `[?YW]Process command.[?RT]
Process file and print result if file out not declarated, else write result to file out.
Usage: 
    [?BK]ten process [file_in] <placeholder_values> <flags> <optional: <[--O/--OUT]:[file_out]> - file out>[?RT]`, []string{"file_in"},
		nil, false)
	p.AddCommand("docs", docs_handler, `[?YW]Show documentation command.[?RT]
Procwss template and show docstring from [?BK]template.config["documentation"][?RT].`,
		[]string{"file_in"}, nil, false)
	return p.Main()
}
