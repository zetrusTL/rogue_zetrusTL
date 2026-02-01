package entity

type RunResult struct {
	PlayerName         string
	DeepestLevel       int
	GoldCollected      int
	EnemiesDefeated    int
	StepsTaken         int
	IsVictory          bool
	AmountConsumedFood int
	NumberElixirDrunk  int
	NumberScrollsRead  int
	TotalHits          int
}

type LevelStats struct {
	AmountConsumedFood int
	NumberElixirDrunk  int
	EnemiesDefeated    int
	MaxHealthBegin     int
}
