package main

import (
	"Aimy/engine"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Aimy")
	if err := ebiten.RunGame(engine.NewGame(640, 480)); err != nil {
		log.Fatal(err)
	}
}
