package json_adapters

import (
	"fmt"

	"example.com/rogue/domain/entity"
)

type GameStateWrapper struct {
	Status            entity.GameStatus `json:"Status"`
	BalanceDifficulty int               `json:"BalanceDifficulty"`
	Messages          []string          `json:"Messages"`
	LevelStats        entity.LevelStats `json:"LevelStats"`
	TotalStats        entity.RunResult  `json:"TotalStats"`

	// Подмена типов на обертки
	Player       *JsonPlayerWrapper `json:"Player"`
	CurrentLevel *JsonLevelWrapper  `json:"CurrentLevel"`
}

func pointToString(p entity.Point) string {
	return fmt.Sprintf("%d,%d", p.X, p.Y)
}

func stringToPoint(s string) entity.Point {
	var p entity.Point
	fmt.Sscanf(s, "%d,%d", &p.X, &p.Y)
	return p
}
