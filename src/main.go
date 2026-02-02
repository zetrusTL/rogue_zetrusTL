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

	r, err := render.NewRender(nil, game, 120, 35)
	if err != nil {
		panic(err)
	}

	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
