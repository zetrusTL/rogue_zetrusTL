package entity

type GameState struct {
	Player       *Player
	CurrentLevel *Level
	TotalStats   RunResult
	LevelStats   LevelStats
	//используется для создания баланса сложности(складывается с номером уровня и получается сложность).
	BalanceDifficulty int
	Status            GameStatus
	Messages          []string
}
