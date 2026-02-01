package entity

type KeyColor int

const (
	NoKey      KeyColor = -1
	RedKey     KeyColor = 0
	BlueKey    KeyColor = 1
	YellowKey  KeyColor = 2
	ColorCount int      = 3
)

type ItemType int

const (
	ItemTreasure ItemType = iota
	ItemFood
	ItemElixir
	ItemScroll
	ItemWeapon
	ItemKey
)

type Direction int

// переименовать top, right, bottom, left
const (
	Top Direction = iota
	Right
	Bottom
	Left
	TopRight
	BottomRight
	BottomLeft
	TopLeft
)

type StatType int

const (
	StatMaxHealth StatType = iota
	StatDexterity
	StatStrength
)

func (stat StatType) String() string {
	switch stat {
	case StatMaxHealth:
		return "Health"
	case StatDexterity:
		return "Dexterity"
	case StatStrength:
		return "Strength"
	default:
		return "Unknown"
	}
}

type MonsterType int

const (
	Zombie MonsterType = iota
	Vampire
	Ghost
	Ogre
	SnakeMage
	Mimic
)

func (monster MonsterType) String() string {
	switch monster {
	case Zombie:
		return "Zombie"
	case Vampire:
		return "Vampire"
	case Ghost:
		return "Ghost"
	case Ogre:
		return "Ogre"
	case SnakeMage:
		return "SnakeMage"
	case Mimic:
		return "Mimic"
	default:
		return "Unknown"
	}
}

type TileType int

const (
	TileWall TileType = iota
	TileFloor
	TileDoorLocked
	TileDoorOpened
	TileCorridor
	TileDownStairs // Выход
)

type GameStatus int

const (
	// StatusRunning игра идёт
	StatusRunning GameStatus = iota
	// StatusGameOver игра проиграна
	StatusGameOver
	// StatusVictory игра выйграна
	StatusVictory
)

const DoorChance float64 = 0.5
const KeyChance float64 = 0.4

const RoomsInRow int = 3
const RoomsInCol int = 3

const MapHeight = 30
const MapWidth = 90

const MaxRoomsNumber = RoomsInRow * RoomsInCol
const MaxCorridorsNumber = 12

const MaxItemsInBackpackOfOneType = 9

const SectorHeight = MapHeight / RoomsInRow
const SectorWidth = MapWidth / RoomsInCol

const MaxLevel = 21

const CornerVertRange = 8
const CornerHorRange = 15

const (
	MaxEnemiesPerRoom = 3
	MaxItemsPerRoom   = 3
)

const (
	MinBalanceDifficulty   = -4
	MaxBalanceDifficulty   = 4
	BalanceDifficultyDelta = 2
)

const VampireMaxHealthDrained = 2
const SnakeAttackSleepChance = 0.25
