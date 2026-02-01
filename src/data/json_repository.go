package data

import (
	"encoding/json"
	"math/rand"
	"os"
	"sort"

	"example.com/rogue/data/config"
	"example.com/rogue/data/json_adapters"
	"example.com/rogue/domain/entity"
)

const (
	gameStateFilename      = "game_state.json"
	statisticFilename      = "statistic.json"
	monstersConfigFilename = "monsters_config.json"
	itemsConfigFilename    = "items_config.json"
	playerConfigFilename   = "player_config.json"
)

type JsonGameRepository struct {
	pathGameState string
	pathStats     string

	monstersRegistry *config.MonstersRegistry
	itemsRegistry    *config.ItemsRegistry
	playerRegistry   *config.PlayerConfig
}

func (repo *JsonGameRepository) CreatePlayer() *entity.Player {
	return &entity.Player{
		CharacterStats: entity.CharacterStats{
			Object:    entity.Object{Size: entity.Size{Width: 1, Height: 1}},
			MaxHealth: repo.playerRegistry.MaxHealth,
			Health:    repo.playerRegistry.MaxHealth,
			Dexterity: repo.playerRegistry.Dexterity,
			Strength:  repo.playerRegistry.Strength,
		},
		Backpack:    entity.NewBackpack(),
		Weapon:      nil,
		ActiveBuffs: make([]*entity.Buff, 0),
		Facing:      entity.Top,
	}
}

func (repo *JsonGameRepository) CreateRandomItemWithBaseCharacteristic(itemType entity.ItemType) entity.Item {
	switch itemType {
	case entity.ItemTreasure:
		itemConfig := repo.itemsRegistry.Treasures[rand.Intn(len(repo.itemsRegistry.Treasures))]
		return &entity.Treasure{Value: itemConfig.Amount, BaseItem: entity.BaseItem{Name: itemConfig.Name}}
	case entity.ItemFood:
		itemConfig := repo.itemsRegistry.Food[rand.Intn(len(repo.itemsRegistry.Food))]

		return &entity.Food{HealAmount: itemConfig.Amount, BaseItem: entity.BaseItem{Name: itemConfig.Name}}
	case entity.ItemElixir:
		itemConfig := repo.itemsRegistry.Potions[rand.Intn(len(repo.itemsRegistry.Potions))]

		return &entity.Elixir{BaseItem: entity.BaseItem{Name: itemConfig.Name},
			Stat:     getStatType(itemConfig.StatAffected),
			Increase: itemConfig.Amount,
			Duration: itemConfig.Duration}
	case entity.ItemScroll:
		itemConfig := repo.itemsRegistry.Scrolls[rand.Intn(len(repo.itemsRegistry.Scrolls))]

		return &entity.Scroll{BaseItem: entity.BaseItem{Name: itemConfig.Name},
			Stat:     getStatType(itemConfig.StatAffected),
			Increase: itemConfig.Amount}
	case entity.ItemWeapon:
		itemConfig := repo.itemsRegistry.Weapons[rand.Intn(len(repo.itemsRegistry.Weapons))]

		return &entity.Weapon{BaseItem: entity.BaseItem{Name: itemConfig.Name},
			StrengthBonus: itemConfig.Amount}
	default:
		panic("Invalid item type")
	}
}

func getStatType(statName string) entity.StatType {
	switch statName {
	case "max_health":
		return entity.StatMaxHealth
	case "strength":
		return entity.StatStrength
	case "dexterity":
		return entity.StatDexterity
	default:
		panic("Invalid stat name")
	}
}

func (repo *JsonGameRepository) CreateMonsterWithBaseCharacteristic(monsterType entity.MonsterType) *entity.Monster {
	var monsterConfig config.MonsterConfig

	switch monsterType {
	case entity.Zombie:
		monsterConfig = repo.monstersRegistry.Zombie
	case entity.Vampire:
		monsterConfig = repo.monstersRegistry.Vampire
	case entity.Ghost:
		monsterConfig = repo.monstersRegistry.Ghost
	case entity.Ogre:
		monsterConfig = repo.monstersRegistry.Ogre
	case entity.SnakeMage:
		monsterConfig = repo.monstersRegistry.SnakeMage
	case entity.Mimic:
		monsterConfig = repo.monstersRegistry.Mimic
	default:
		panic("Invalid monster type")
	}

	return &entity.Monster{
		Hostility: monsterConfig.Hostility,
		Type:      monsterType,
		WasHit:    false,
		IsVisible: true,
		IsResting: false,
		CharacterStats: entity.CharacterStats{
			Health:    monsterConfig.MaxHealth,
			Strength:  monsterConfig.Strength,
			Dexterity: monsterConfig.Dexterity,
			MaxHealth: monsterConfig.MaxHealth,
			Object:    entity.Object{Size: entity.Size{Width: 1, Height: 1}}},
	}
}

func (repo *JsonGameRepository) SaveSession(state *entity.GameState) error {
	wrapper := json_adapters.GameStateWrapper{
		Status:            state.Status,
		BalanceDifficulty: state.BalanceDifficulty,
		Messages:          state.Messages,
		LevelStats:        state.LevelStats,
		TotalStats:        state.TotalStats,
		Player:            &json_adapters.JsonPlayerWrapper{Player: state.Player},
		CurrentLevel:      &json_adapters.JsonLevelWrapper{Level: state.CurrentLevel},
	}

	data, err := json.MarshalIndent(wrapper, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(repo.pathGameState, data, 0644)
}

func (repo *JsonGameRepository) LoadSession() (*entity.GameState, error) {
	fileData, err := os.ReadFile(repo.pathGameState)
	if err != nil {
		return nil, err
	}

	state := &entity.GameState{
		Player:       &entity.Player{Backpack: entity.NewBackpack()},
		CurrentLevel: &entity.Level{},
	}

	wrapper := json_adapters.GameStateWrapper{
		Player:       &json_adapters.JsonPlayerWrapper{Player: state.Player},
		CurrentLevel: &json_adapters.JsonLevelWrapper{Level: state.CurrentLevel},
	}

	err = json.Unmarshal(fileData, &wrapper)
	if err != nil {
		return nil, err
	}

	state.Status = wrapper.Status
	state.BalanceDifficulty = wrapper.BalanceDifficulty
	state.Messages = wrapper.Messages
	state.LevelStats = wrapper.LevelStats
	state.TotalStats = wrapper.TotalStats

	return state, nil
}

func (repo *JsonGameRepository) HasSavedSession() bool {
	if _, err := os.Stat(repo.pathGameState); err == nil {
		return true
	}

	return false
}

func (repo *JsonGameRepository) DeleteSession() error {
	return os.Remove(repo.pathGameState)
}

func (repo *JsonGameRepository) SaveRunResult(result entity.RunResult) error {
	leaderboard, err := repo.GetLeaderboard()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	leaderboard = append(leaderboard, result)

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].GoldCollected > leaderboard[j].GoldCollected
	})

	data, err := json.MarshalIndent(leaderboard, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(repo.pathStats, data, 0644)
}

func (repo *JsonGameRepository) GetLeaderboard() ([]entity.RunResult, error) {
	fileData, err := os.ReadFile(repo.pathStats)
	if err != nil {
		if os.IsNotExist(err) {
			return []entity.RunResult{}, nil
		}
		return nil, err
	}

	var leaderboard []entity.RunResult
	err = json.Unmarshal(fileData, &leaderboard)
	if err != nil {
		return nil, err
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].GoldCollected > leaderboard[j].GoldCollected
	})

	return leaderboard, nil
}

func NewJsonRepository(folder string) (*JsonGameRepository, error) {
	repo := &JsonGameRepository{
		pathGameState: folder + "/" + gameStateFilename,
		pathStats:     folder + "/" + statisticFilename}

	monstersRegistry, err := config.LoadMonsterConfigs(folder + "/" + monstersConfigFilename)
	itemsRegistry, err := config.LoadItemsConfigs(folder + "/" + itemsConfigFilename)
	playerRegistry, err := config.LoadPlayerConfigs(folder + "/" + playerConfigFilename)

	repo.monstersRegistry = monstersRegistry
	repo.itemsRegistry = itemsRegistry
	repo.playerRegistry = playerRegistry

	if err != nil {
		return nil, err
	}

	return repo, nil
}
