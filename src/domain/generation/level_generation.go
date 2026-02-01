package generation

import (
	"example.com/rogue/domain"
	"example.com/rogue/domain/entity"
)

type LevelGenerator struct {
	repo       domain.GameRepository
	level      *entity.Level
	startRoom  *entity.Room
	difficulty int
}

func NewLevelGenerator(repo domain.GameRepository) *LevelGenerator {
	return &LevelGenerator{repo: repo}
}

func (lg *LevelGenerator) GenerateLevel(numLevel int, balanceDifficulty int) *entity.Level {
	lg.difficulty = numLevel + balanceDifficulty
	lg.level, lg.startRoom = lg.generateMap(numLevel)
	lg.generateEntities()

	return lg.level
}
