package main

import (
	"Aimy/engine"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Printf("Hello")
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Aimy")
	if err := ebiten.RunGame(engine.NewGame(640, 480)); err != nil {
		log.Fatal(err)
	}
}
