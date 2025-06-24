package main

import (
    _ "framework/cache/one"
    _ "framework/cache/two"
	"log"
	"framework/core/run"
)
func main() {
	log.Println("code generate")
	run.RunCase()
}
