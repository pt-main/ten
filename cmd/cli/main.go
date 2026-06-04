package main

import "github.com/pt-main/tap/color"

func main() {
	err := cli()
	if err != nil {
		color.PrintlnColored("[?RD]%s[?RT]", err.Error())
	}
}
