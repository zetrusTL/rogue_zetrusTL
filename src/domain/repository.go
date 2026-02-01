package domain

import "example.com/rogue/domain/entity"

type GameRepository interface {
	// Работа с текущей сессией (Save/Load/Resume)
	SaveSession(state *entity.GameState) error
	LoadSession() (*entity.GameState, error)
	HasSavedSession() bool
	DeleteSession() error

	SaveRunResult(result entity.RunResult) error
	GetLeaderboard() ([]entity.RunResult, error)

	CreateMonsterWithBaseCharacteristic(monsterType entity.MonsterType) *entity.Monster
	CreateRandomItemWithBaseCharacteristic(itemType entity.ItemType) entity.Item
	CreatePlayer() *entity.Player
}
