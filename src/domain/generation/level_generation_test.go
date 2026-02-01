package generation

import (
	"testing"

	"example.com/rogue/domain/entity"
)

// --- MOCK REPOSITORY ---

type MockRepository struct{}

func (m *MockRepository) SaveSession(state *entity.GameState) error   { return nil }
func (m *MockRepository) LoadSession() (*entity.GameState, error)     { return nil, nil }
func (m *MockRepository) HasSavedSession() bool                       { return false }
func (m *MockRepository) DeleteSession() error                        { return nil }
func (m *MockRepository) SaveRunResult(result entity.RunResult) error { return nil }
func (m *MockRepository) GetLeaderboard() ([]entity.RunResult, error) { return nil, nil }
func (m *MockRepository) CreatePlayer() *entity.Player                { return &entity.Player{} }

func (m *MockRepository) CreateMonsterWithBaseCharacteristic(t entity.MonsterType) *entity.Monster {
	// Возвращаем монстра с базовыми статами для тестов
	return &entity.Monster{
		CharacterStats: entity.CharacterStats{
			Health:    10,
			MaxHealth: 10,
			Strength:  2,
			Dexterity: 2,
			Object:    entity.Object{Size: entity.Size{Width: 1, Height: 1}},
		},
		Type:      t,
		Hostility: 5,
	}
}

func (m *MockRepository) CreateRandomItemWithBaseCharacteristic(t entity.ItemType) entity.Item {
	// Возвращаем базовый предмет
	base := entity.BaseItem{Name: "Test Item", Pos: entity.Point{}}
	switch t {
	case entity.ItemFood:
		return &entity.Food{BaseItem: base, HealAmount: 5}
	case entity.ItemWeapon:
		return &entity.Weapon{BaseItem: base, StrengthBonus: 2}
	default:
		return &entity.Treasure{BaseItem: base, Value: 10}
	}
}

// --- TESTS ---

func TestGenerateLevel_Structure(t *testing.T) {
	repo := &MockRepository{}
	gen := NewLevelGenerator(repo)

	level := gen.GenerateLevel(1, 0)

	if level == nil {
		t.Fatal("GenerateLevel returned nil")
	}

	// 1. Проверяем наличие комнат
	if len(level.Rooms) != 9 {
		t.Errorf("Expected 9 rooms, got %d", len(level.Rooms))
	}

	// 2. Проверяем старт и выход
	if gen.startRoom == nil {
		t.Fatal("StartRoom is nil")
	}

	// Проверяем, что выход не в стартовой комнате (в большинстве случаев, рандом может совпасть только при баге логики placeStartAndExit)
	// В коде placeStartAndExit есть цикл for endRoom == startRoom, так что они должны быть разные.
	if level.StartPos == level.ExitPos {
		t.Error("StartPos and ExitPos should not be the same")
	}

	// 3. Проверяем тайл выхода
	tileAtExit := level.GetTile(level.ExitPos.X, level.ExitPos.Y)
	if tileAtExit != entity.TileDownStairs {
		t.Errorf("Expected TileDownStairs at ExitPos %v, got %v", level.ExitPos, tileAtExit)
	}

	// 4. Проверяем, что стартовая позиция внутри стартовой комнаты
	sr := gen.startRoom
	if level.StartPos.X < sr.Pos.X || level.StartPos.X >= sr.Pos.X+sr.Size.Width ||
		level.StartPos.Y < sr.Pos.Y || level.StartPos.Y >= sr.Pos.Y+sr.Size.Height {
		t.Errorf("StartPos %v is outside StartRoom bounds %v (Size %v)", level.StartPos, sr.Pos, sr.Size)
	}
}

func TestStartRoomSafe(t *testing.T) {
	repo := &MockRepository{}
	gen := NewLevelGenerator(repo)

	// Генерируем несколько раз, чтобы исключить случайное совпадение
	for i := 0; i < 5; i++ {
		level := gen.GenerateLevel(1, 0)
		startRoom := gen.startRoom

		// Проверяем всех монстров на уровне
		for pos, _ := range level.Monsters {
			// Проверяем, попадает ли монстр в прямоугольник стартовой комнаты
			if pos.X >= startRoom.Pos.X && pos.X < startRoom.Pos.X+startRoom.Size.Width &&
				pos.Y >= startRoom.Pos.Y && pos.Y < startRoom.Pos.Y+startRoom.Size.Height {
				t.Errorf("Found monster at %v inside StartRoom!", pos)
			}
		}
	}
}

func TestLevelConnectivity(t *testing.T) {
	// Проверка того, что от старта можно дойти до выхода (геометрически)
	repo := &MockRepository{}
	gen := NewLevelGenerator(repo)
	level := gen.GenerateLevel(1, 0)

	start := level.StartPos
	end := level.ExitPos

	// BFS
	queue := []entity.Point{start}
	visited := make(map[entity.Point]bool)
	visited[start] = true
	found := false

	// Смещения для соседей (верх, низ, лево, право)
	dirs := []entity.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			found = true
			break
		}

		for _, d := range dirs {
			next := entity.Point{X: curr.X + d.X, Y: curr.Y + d.Y}

			// Проверка границ карты
			// В коде генератора есть entity.MapWidth, entity.MapHeight
			if next.X < 0 || next.Y < 0 || next.X >= entity.MapWidth || next.Y >= entity.MapHeight {
				continue
			}

			if visited[next] {
				continue
			}

			tile := level.GetTile(next.X, next.Y)
			// Считаем проходимыми: Пол, Коридор, Дверь (закрытая или открытая), Лестница
			isPassable := tile == entity.TileFloor ||
				tile == entity.TileCorridor ||
				tile == entity.TileDoorLocked ||
				tile == entity.TileDoorOpened ||
				tile == entity.TileDownStairs

			if isPassable {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}

	if !found {
		t.Error("Exit is not reachable from Start position (connectivity broken)")
	}
}

func TestDifficultyScaling(t *testing.T) {
	repo := &MockRepository{}
	gen := NewLevelGenerator(repo)

	// Уровень 1
	lvl1 := gen.GenerateLevel(1, 0)
	var totalHP1, count1 int
	for _, m := range lvl1.Monsters {
		totalHP1 += m.MaxHealth
		count1++
	}
	avgHP1 := float64(totalHP1)
	if count1 > 0 {
		avgHP1 /= float64(count1)
	}

	// Уровень 10 (с балансом 0)
	lvl10 := gen.GenerateLevel(10, 0)
	var totalHP10, count10 int
	for _, m := range lvl10.Monsters {
		totalHP10 += m.MaxHealth
		count10++
	}
	avgHP10 := float64(totalHP10)
	if count10 > 0 {
		avgHP10 /= float64(count10)
	}

	// Проверки
	// 1. Здоровье должно расти
	if avgHP10 <= avgHP1 {
		t.Errorf("Expected Level 10 monsters to have more HP than Level 1. Lvl1 Avg: %.2f, Lvl10 Avg: %.2f", avgHP1, avgHP10)
	}

	// 2. Количество монстров должно расти (проверяем raw count, хотя рандом может влиять, разница в 9 уровней должна быть заметна)
	if count10 < count1 {
		t.Logf("Warning: Level 10 has fewer or equal monsters (%d) than Level 1 (%d). Could be RNG, but expected increase.", count10, count1)
	} else {
		t.Logf("Scaling Check Passed: Level 1 monsters: %d (AvgHP: %.2f), Level 10 monsters: %d (AvgHP: %.2f)", count1, avgHP1, count10, avgHP10)
	}
}
