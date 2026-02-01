package implementation

import (
	"math"

	"example.com/rogue/domain/entity"
)

func calculateDestination(pos entity.Point, dir entity.Direction) entity.Point {
	adjacentTiles := []entity.Point{
		{X: pos.X, Y: pos.Y - 1},
		{X: pos.X + 1, Y: pos.Y},
		{X: pos.X, Y: pos.Y + 1},
		{X: pos.X - 1, Y: pos.Y},
		{X: pos.X + 1, Y: pos.Y - 1},
		{X: pos.X + 1, Y: pos.Y + 1},
		{X: pos.X - 1, Y: pos.Y + 1},
		{X: pos.X - 1, Y: pos.Y - 1},
	}
	destinations := map[entity.Direction]entity.Point{
		entity.Top:         adjacentTiles[0],
		entity.Right:       adjacentTiles[1],
		entity.Bottom:      adjacentTiles[2],
		entity.Left:        adjacentTiles[3],
		entity.TopRight:    adjacentTiles[4],
		entity.BottomRight: adjacentTiles[5],
		entity.BottomLeft:  adjacentTiles[6],
		entity.TopLeft:     adjacentTiles[7],
	}

	return destinations[dir]
}

func distance(p1, p2 entity.Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(float64(dx*dx + dy*dy))
}

func calculateHitChance(attacking entity.CharacterStats, defeating entity.CharacterStats) int {
	hitChance := 50 + (attacking.Dexterity-defeating.Dexterity)*2

	if hitChance < 10 {
		hitChance = 10
	}
	if hitChance > 100 {
		hitChance = 100
	}

	return hitChance
}

func (g *GameService) getValidNeighborsForMonster(p entity.Point, isChasing bool) []entity.Point {
	res := make([]entity.Point, 0, 8)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			np := entity.Point{X: p.X + dx, Y: p.Y + dy}
			if g.monsterCanMoveTo(np, isChasing) {
				res = append(res, np)
			}
		}
	}
	return res
}

func (g *GameService) monsterCanMoveTo(p entity.Point, isChasing bool) bool {
	// Границы карты
	if p.X < 0 || p.Y < 0 || p.X >= entity.MapWidth || p.Y >= entity.MapHeight {
		return false
	}

	// Стены и закрытые двери
	tile := g.state.CurrentLevel.GetTile(p.X, p.Y)
	if tile == entity.TileWall || tile == entity.TileDoorLocked || tile == entity.TileDownStairs {
		return false
	}

	//Если монстр не преследует игрока, тогда двигаемся в рамках коридора или комнаты
	if !isChasing && tile == entity.TileDoorOpened {
		return false
	}

	// Другие монстры
	if _, occupied := g.state.CurrentLevel.Monsters[p]; occupied {
		return false
	}

	return true
}

func (g *GameService) getValidNeighbors(p entity.Point) []entity.Point {
	res := make([]entity.Point, 0, 8)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			np := entity.Point{X: p.X + dx, Y: p.Y + dy}
			if g.canMoveTo(np) {
				res = append(res, np)
			}
		}
	}
	return res
}

func (g *GameService) canMoveTo(p entity.Point) bool {
	// Границы карты
	if p.X < 0 || p.Y < 0 || p.X >= entity.MapWidth || p.Y >= entity.MapHeight {
		return false
	}

	// Стены и закрытые двери
	tile := g.state.CurrentLevel.GetTile(p.X, p.Y)
	if tile == entity.TileWall || tile == entity.TileDoorLocked {
		return false
	}

	// Игрок
	if p == g.GetPlayer().Pos {
		return false
	}

	return true
}
