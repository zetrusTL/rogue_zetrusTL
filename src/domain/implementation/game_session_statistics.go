package implementation

import "example.com/rogue/domain/entity"

func (g *GameService) writeDefeatedEnemyStat() {
	g.state.TotalStats.EnemiesDefeated++
	g.state.LevelStats.EnemiesDefeated++
}

func (g *GameService) writeStepsTakenStat() {
	g.state.TotalStats.StepsTaken++
}

func (g *GameService) writeAmountConsumedFoodStat() {
	g.state.TotalStats.AmountConsumedFood++
	g.state.LevelStats.AmountConsumedFood++
}

func (g *GameService) writeNumberElixirDrunkStat() {
	g.state.TotalStats.NumberElixirDrunk++
	g.state.LevelStats.NumberElixirDrunk++
}

func (g *GameService) writeNumberScrollsReadStat() {
	g.state.TotalStats.NumberScrollsRead++
}

func (g *GameService) writeTotalHitsStat(hits int) {
	g.state.TotalStats.TotalHits++
}

func (g *GameService) GetStats() entity.RunResult {
	g.refreshStats()
	if g.state.Status == entity.StatusVictory {
		g.state.TotalStats.IsVictory = true
	} else {
		g.state.TotalStats.IsVictory = false
	}
	return g.state.TotalStats
}

func (g *GameService) refreshStats() {
	g.state.TotalStats.DeepestLevel = g.GetCurrentLevel().LevelNum
	g.state.TotalStats.GoldCollected = g.GetPlayer().Backpack.Gold
}
