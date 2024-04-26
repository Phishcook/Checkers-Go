package main

import (
	"log"

	"github.com/Phishcook/Checkers-Go/checkers"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game, err := checkers.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Phunky Checkers")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
