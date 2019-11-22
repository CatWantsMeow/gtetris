package main

import (
	"flag"

	"github.com/CatWantsMeow/gtetris/game"
)

func main() {
	debug := flag.Bool("debug", false, "Run game in debug mode.")
	fast := flag.Bool("fast", false, "Speeds up game.")
	flag.Parse()

	game.Run(*debug, *fast)
}
