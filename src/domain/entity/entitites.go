package entity

type Point struct {
	X, Y int
}

type Size struct {
	Width, Height int
}

type Object struct {
	Pos  Point
	Size Size
}

type CharacterStats struct {
	Object
	Health    int
	MaxHealth int
	Dexterity int
	Strength  int
}

// Item - общий интерфейс для всего, что можно положить в рюкзак
type Item interface {
	GetType() ItemType
	GetName() string
	GetPosition() Point
	SetPosition(p Point)
}

type BaseItem struct {
	Name string
	Pos  Point
}

func (b *BaseItem) GetName() string     { return b.Name }
func (b *BaseItem) GetPosition() Point  { return b.Pos }
func (b *BaseItem) SetPosition(p Point) { b.Pos = p }

type Treasure struct {
	BaseItem
	Value int
}

func (t *Treasure) GetType() ItemType { return ItemTreasure }

type Food struct {
	BaseItem
	HealAmount int
}

func (f *Food) GetType() ItemType { return ItemFood }

type Elixir struct {
	BaseItem
	Stat     StatType
	Increase int
	Duration int // В ходах (turns)
}

func (e *Elixir) GetType() ItemType { return ItemElixir }

type Scroll struct {
	BaseItem
	Stat     StatType
	Increase int // Permanent increase
}

func (s *Scroll) GetType() ItemType { return ItemScroll }

type Weapon struct {
	BaseItem
	StrengthBonus int
}

func (w *Weapon) GetType() ItemType { return ItemWeapon }

type Key struct {
	BaseItem
	Color KeyColor
}

func (k *Key) GetType() ItemType { return ItemKey }

type Backpack struct {
	Items map[ItemType][]Item
	Gold  int
}

func NewBackpack() Backpack {
	return Backpack{
		Items: map[ItemType][]Item{
			ItemFood:   make([]Item, 0),
			ItemElixir: make([]Item, 0),
			ItemScroll: make([]Item, 0),
			ItemWeapon: make([]Item, 0),
			ItemKey:    make([]Item, 0),
		},
		Gold: 0,
	}
}

type Buff struct {
	Stat      StatType
	Increase  int
	TurnsLeft int // Время действия в ходах (через сколько закончится)
}

type Player struct {
	CharacterStats
	Backpack    Backpack
	Weapon      *Weapon
	ActiveBuffs []*Buff

	// Для 3D режима
	Facing Direction

	SleepTurns int
}

type Monster struct {
	CharacterStats
	Type      MonsterType
	Hostility int

	MoveDir Direction

	// Vampire
	WasHit bool // Был ли уже атакован (для уворота от первого удара)

	// Ghost
	IsVisible bool // false = невидимость

	// Ogre
	IsResting bool // Пропуск хода после атаки

	// Mimic
	IsDisguised  bool     // Замаскирован ли
	DisguiseItem ItemType // Предмет, под который маскируется
}
