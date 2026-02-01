package domain

import "example.com/rogue/domain/entity"

type GameSession interface {
	GetPlayer() *entity.Player
	GetCurrentLevel() *entity.Level
	GetMessages() []string
	GetStatus() entity.GameStatus

	//движение игрока
	MovePlayer(dir entity.Direction)
	//пропуск кода игроком
	WaitTurn()

	// Инвентарь
	//Получить предметы конкретного типа
	GetInventoryItems(itemType entity.ItemType) []entity.Item
	//1-9 выбор предмета
	UseItem(itemType entity.ItemType, index int) error
	DropItem(itemType entity.ItemType, index int) error

	//0 - убирает оружие из руки, 1-9 выбор оружия
	EquipWeapon(index int)

	// 3D Actions
	RotatePlayer(dir entity.Direction)

	GetStats() entity.RunResult
}

type Game interface {
	StartNewGame(playerName string) (GameSession, error)
	LoadLastGame() (GameSession, error)
	HasSavedGame() bool

	GetLeaderboard(limit int) []entity.RunResult
}
