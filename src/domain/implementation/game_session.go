package implementation

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"example.com/rogue/domain"
	"example.com/rogue/domain/entity"
	"example.com/rogue/domain/generation"
)

type GameService struct {
	repo domain.GameRepository

	levelGenerator *generation.LevelGenerator
	state          *entity.GameState
}

func NewGameService(repo domain.GameRepository) *GameService {
	return &GameService{repo: repo, levelGenerator: generation.NewLevelGenerator(repo)}
}

func (g *GameService) GetInventoryItems(itemType entity.ItemType) []entity.Item {
	return g.state.Player.Backpack.Items[itemType]
}

func (g *GameService) UseItem(itemType entity.ItemType, index int) error {
	index--
	if len(g.state.Player.Backpack.Items[itemType]) <= index || index < 0 {
		return errors.New("invalid index")
	}

	g.clearMessages()

	switch itemType {
	case entity.ItemFood:
		food := g.state.Player.Backpack.Items[itemType][index].(*entity.Food)
		restored := g.state.Player.Health
		g.state.Player.Health = min(g.state.Player.Health+food.HealAmount, g.state.Player.MaxHealth)
		restored = g.state.Player.Health - restored
		g.logMessage(fmt.Sprintf("You ate '%s'", food.GetName()))
		g.writeAmountConsumedFoodStat()
	case entity.ItemElixir:
		elixir := g.state.Player.Backpack.Items[itemType][index].(*entity.Elixir)
		g.state.Player.ActiveBuffs = append(g.state.Player.ActiveBuffs, &entity.Buff{Stat: elixir.Stat, TurnsLeft: elixir.Duration, Increase: elixir.Increase})
		stat := increaseStat(g.state.Player, elixir.Stat, elixir.Increase)

		g.logMessage(fmt.Sprintf("You used '%s'", elixir.GetName()), fmt.Sprintf("Increased %s on %d points for %d turns", stat, elixir.Increase, elixir.Duration))
		g.writeNumberElixirDrunkStat()
	case entity.ItemScroll:
		scroll := g.state.Player.Backpack.Items[itemType][index].(*entity.Scroll)
		stat := increaseStat(g.state.Player, scroll.Stat, scroll.Increase)

		g.logMessage(fmt.Sprintf("You used '%s'", scroll.GetName()), fmt.Sprintf("Increased %s on %d points permanently", stat, scroll.Increase))
		g.writeNumberScrollsReadStat()
	default:
		return errors.New("invalid item type")
	}

	g.state.Player.Backpack.Items[itemType] = append(g.state.Player.Backpack.Items[itemType][:index], g.state.Player.Backpack.Items[itemType][index+1:]...)

	g.UpdateWorld()

	return nil
}

func increaseStat(player *entity.Player, stat entity.StatType, amount int) string {
	switch stat {
	case entity.StatMaxHealth:
		player.Health += amount
		player.MaxHealth += amount
		return "Max health"
	case entity.StatDexterity:
		player.Dexterity += amount
		return "Dexterity"
	case entity.StatStrength:
		player.Strength += amount
		return "Strength"
	default:
		panic("invalid elixir stat")
	}
}

func (g *GameService) DropItem(itemType entity.ItemType, index int) error {
	index--
	if len(g.state.Player.Backpack.Items[itemType]) <= index || index < 0 {
		return errors.New("cannot drop item")
	}

	g.clearMessages()

	err := g.dropItem(itemType, index)

	if err != nil {
		return err
	}

	g.UpdateWorld()

	return nil
}

func (g *GameService) GetPlayer() *entity.Player {
	return g.state.Player
}

func (g *GameService) GetCurrentLevel() *entity.Level {
	return g.state.CurrentLevel
}

func (g *GameService) GetMessages() []string {
	return g.state.Messages
}

func (g *GameService) GetStatus() entity.GameStatus {
	return g.state.Status
}

func (g *GameService) WaitTurn() {
	g.UpdateWorld()
}

func (g *GameService) EquipWeapon(index int) {
	g.clearMessages()

	var indexToRemove = -1
	index--
	if g.state.Player.Weapon != nil {
		for i, weapon := range g.state.Player.Backpack.Items[entity.ItemWeapon] {
			if weapon == g.state.Player.Weapon {
				indexToRemove = i
				break
			}
		}
	}

	if index == -1 {
		err := g.dropItem(entity.ItemWeapon, indexToRemove)
		if err == nil {
			g.UpdateWorld()
		}

		return
	}

	if indexToRemove == index {
		g.logMessage("You already holding this weapon")
		return
	}

	if indexToRemove != -1 {
		err := g.DropItem(entity.ItemWeapon, indexToRemove)
		if err != nil {
			return
		}
		if index > indexToRemove {
			index--
		}
	}

	g.state.Player.Weapon = g.state.Player.Backpack.Items[entity.ItemWeapon][index].(*entity.Weapon)
	g.logMessage(fmt.Sprintf("You equiped %s", g.state.Player.Weapon.GetName()))
	g.UpdateWorld()
}

func (g *GameService) RotatePlayer(dir entity.Direction) {
	g.state.Player.Facing = dir
}

// если true, то создается игра из сейва, false - новая игра
func (g *GameService) init(playerName string, saved bool) error {
	var err error
	if saved && g.repo.HasSavedSession() {
		g.state, err = g.repo.LoadSession()
		if err != nil {
			return err
		}

		return nil
	}

	g.state = &entity.GameState{}

	g.state.Player = g.repo.CreatePlayer()
	g.state.TotalStats = entity.RunResult{
		PlayerName: playerName,
	}
	g.nextLevel(0, 0)
	g.state.BalanceDifficulty = 0
	g.state.Status = entity.StatusRunning
	g.state.Messages = make([]string, 0)

	return nil
}

func (g *GameService) nextLevel(prevLevelNum int, balanceDifficulty int) {
	if prevLevelNum == entity.MaxLevel {
		g.endGame(true)
		return
	}

	g.GetPlayer().Backpack.Items[entity.ItemKey] = make([]entity.Item, 0)
	g.adjustBalanceAfterLevel()
	g.state.CurrentLevel = g.levelGenerator.GenerateLevel(prevLevelNum+1, balanceDifficulty)
	g.GetPlayer().Pos = g.state.CurrentLevel.StartPos

	err := g.repo.SaveSession(g.state)
	if err != nil {
		panic(err)
	}
}

func (g *GameService) adjustBalanceAfterLevel() {
	player := g.state.Player
	stats := g.state.LevelStats

	hpRatio := float64(player.Health) / float64(player.MaxHealth)

	easySignals := 0
	hardSignals := 0

	// --- EASY ---
	if hpRatio > 0.7 {
		easySignals++
	}
	if stats.EnemiesDefeated >= int(math.RoundToEven(0.7*float64(entity.MaxEnemiesPerRoom*9))) {
		easySignals++
	}
	if stats.AmountConsumedFood == 0 && stats.NumberElixirDrunk == 0 {
		easySignals++
	}

	// --- HARD ---
	if hpRatio < 0.3 {
		hardSignals++
	}
	if stats.AmountConsumedFood+stats.NumberElixirDrunk >= 3 {
		hardSignals++
	}
	if g.GetPlayer().MaxHealth < int(float64(stats.MaxHealthBegin)*0.3) {
		hardSignals++
	}

	if easySignals >= 2 && g.state.BalanceDifficulty < entity.MaxBalanceDifficulty {
		g.state.BalanceDifficulty += entity.BalanceDifficultyDelta
	} else if hardSignals >= 2 && g.state.BalanceDifficulty > entity.MinBalanceDifficulty {
		g.state.BalanceDifficulty += entity.BalanceDifficultyDelta
	}

	stats.EnemiesDefeated = 0
	stats.AmountConsumedFood = 0
	stats.NumberElixirDrunk = 0
	stats.MaxHealthBegin = g.GetPlayer().MaxHealth
}

func (g *GameService) clearMessages() {
	g.state.Messages = make([]string, 0)
}

func (g *GameService) logMessage(messages ...string) {
	g.state.Messages = append(g.state.Messages, messages...)
}

func (g *GameService) endGame(victory bool) {
	if victory {
		g.state.Status = entity.StatusVictory
	} else {
		g.state.Status = entity.StatusGameOver
	}

	g.state.TotalStats.IsVictory = victory

	err := g.repo.DeleteSession()

	if err != nil {
		panic(err)
	}

	g.refreshStats()
	err = g.repo.SaveRunResult(g.state.TotalStats)
	if err != nil {
		panic(err)
	}
}

func (g *GameService) dropItem(itemType entity.ItemType, index int) error {
	possibleCells := g.getValidNeighbors(g.GetPlayer().Pos)

	if len(possibleCells) == 0 {
		g.logMessage("There is no tile to drop item to")
		return errors.New("cannot drop item")
	}

	item := g.state.Player.Backpack.Items[itemType][index]

	item.SetPosition(possibleCells[rand.Intn(len(possibleCells))])
	g.state.CurrentLevel.Items[item.GetPosition()] = append(g.state.CurrentLevel.Items[item.GetPosition()], item)

	g.state.Player.Backpack.Items[itemType] = append(g.state.Player.Backpack.Items[itemType][:index], g.state.Player.Backpack.Items[itemType][index+1:]...)

	g.logMessage(fmt.Sprintf("You droped '%s'", item.GetName()))

	return nil
}
