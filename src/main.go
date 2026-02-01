package main

import (
	"fmt"

	"example.com/rogue/data"
	"example.com/rogue/domain/implementation"
	"example.com/rogue/presentation/render"
)

func main() {
	repo, err := data.NewJsonRepository("data/files")
		if err != nil {
			panic(err)
		}

	game := implementation.NewGameImpl(repo)

	session, err := game.StartNewGame("Player")
	if err != nil {
		panic(err)
	}

	r, err := render.NewRender(session, game, 120, 35)
	if err != nil {
		panic(err)
	}

	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
