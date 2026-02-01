package implementation

import (
	"example.com/rogue/domain"
	"example.com/rogue/domain/entity"
)

type GameImpl struct {
	repo domain.GameRepository
}

func NewGameImpl(repo domain.GameRepository) *GameImpl {
	return &GameImpl{repo: repo}
}

func (g GameImpl) StartNewGame(playerName string) (domain.GameSession, error) {
	gameSession := NewGameService(g.repo)
	err := gameSession.init(playerName, false)

	return gameSession, err
}

func (g GameImpl) LoadLastGame() (domain.GameSession, error) {
	gameSession := NewGameService(g.repo)
	err := gameSession.init("player", true)

	return gameSession, err
}

func (g GameImpl) HasSavedGame() bool {
	return g.repo.HasSavedSession()
}

func (g GameImpl) GetLeaderboard(limit int) []entity.RunResult {
	board, err := g.repo.GetLeaderboard()

	if err != nil {
		return make([]entity.RunResult, 0)
	}

	if len(board) >= limit {
		return board[:limit]
	}

	return board
}
