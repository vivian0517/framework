package main

import (
	"framework/core"
	"framework/core/embed"
	"os"
)

func main() {
	args := os.Args[1]
	switch args {
	case "book":
		embed.Book()
	case "run":
		core.Core()
	}
}
