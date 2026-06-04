package main

import "github.com/pt-main/tap/color"

const Version = "1.0.5"

func main() {
	err := cli()
	if err != nil {
		color.PrintlnColored("[?RD]%s[?RT]", err.Error())
	}
}
