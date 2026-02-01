package implementation

import (
	"testing"

	"example.com/rogue/domain/entity"
)

// --- MOCKS ---

type MockRepo struct {
	SavedState *entity.GameState
	HasSave    bool
}

func (m *MockRepo) SaveSession(state *entity.GameState) error {
	m.SavedState = state
	m.HasSave = true
	return nil
}

func (m *MockRepo) LoadSession() (*entity.GameState, error) {
	if m.HasSave {
		return m.SavedState, nil
	}
	return nil, nil
}

func (m *MockRepo) HasSavedSession() bool {
	return m.HasSave
}

func (m *MockRepo) DeleteSession() error {
	m.HasSave = false
	m.SavedState = nil
	return nil
}

func (m *MockRepo) SaveRunResult(result entity.RunResult) error {
	return nil
}

func (m *MockRepo) GetLeaderboard() ([]entity.RunResult, error) {
	return nil, nil
}

func (m *MockRepo) CreatePlayer() *entity.Player {
	return &entity.Player{
		CharacterStats: entity.CharacterStats{
			Health:    100,
			MaxHealth: 100,
			Strength:  10,
			Dexterity: 10,
		},
		Backpack: entity.NewBackpack(),
	}
}

func (m *MockRepo) CreateMonsterWithBaseCharacteristic(t entity.MonsterType) *entity.Monster {
	return &entity.Monster{
		CharacterStats: entity.CharacterStats{
			Health:    20,
			MaxHealth: 20,
			Strength:  5,
			Dexterity: 5,
		},
		Type:      t,
		Hostility: 5,
	}
}

func (m *MockRepo) CreateRandomItemWithBaseCharacteristic(t entity.ItemType) entity.Item {
	base := entity.BaseItem{Name: "Test Item"}
	switch t {
	case entity.ItemFood:
		return &entity.Food{BaseItem: base, HealAmount: 10}
	case entity.ItemWeapon:
		return &entity.Weapon{BaseItem: base, StrengthBonus: 5}
	case entity.ItemKey:
		return &entity.Key{BaseItem: base, Color: entity.RedKey}
	default:
		return &entity.Treasure{BaseItem: base, Value: 10}
	}
}

// --- SETUP HELPERS ---

func setupService() *GameService {
	repo := &MockRepo{}
	service := NewGameService(repo)
	// Инициализируем новую игру
	service.init("TestPlayer", false)
	return service
}

// Создает мини-уровень для тестов движения
func createTestLevel(playerPos entity.Point) *entity.Level {
	lvl := entity.NewLevel(1)
	// Заполняем стенами
	for y := 0; y < entity.MapHeight; y++ {
		for x := 0; x < entity.MapWidth; x++ {
			lvl.SetTile(x, y, entity.TileWall)
		}
	}
	// Делаем площадку 3x3 пола вокруг игрока
	for y := playerPos.Y - 1; y <= playerPos.Y+1; y++ {
		for x := playerPos.X - 1; x <= playerPos.X+1; x++ {
			lvl.SetTile(x, y, entity.TileFloor)
		}
	}
	lvl.StartPos = playerPos
	return lvl
}

// --- TESTS ---

func TestStartNewGame(t *testing.T) {
	service := setupService()

	if service.GetPlayer() == nil {
		t.Error("Player should be initialized")
	}
	if service.GetCurrentLevel() == nil {
		t.Error("Level should be generated")
	}
	if service.GetStatus() != entity.StatusRunning {
		t.Errorf("Expected status Running, got %v", service.GetStatus())
	}
}

func TestMovePlayer_Movement(t *testing.T) {
	service := setupService()
	startPos := entity.Point{X: 10, Y: 10}

	// Подменяем уровень на контролируемый
	service.state.CurrentLevel = createTestLevel(startPos)
	service.GetPlayer().Pos = startPos

	// 1. Движение вправо (на пол)
	service.MovePlayer(entity.Right)
	expected := entity.Point{X: 11, Y: 10}
	if service.GetPlayer().Pos != expected {
		t.Errorf("Player failed to move right. Got %v, expected %v", service.GetPlayer().Pos, expected)
	}

	// 2. Попытка движения в стену (вправо от 11,10 - это 12,10, там стена по логике createTestLevel, т.к. пол только 3x3)
	// createTestLevel(10,10) -> x: 9,10,11 are floors. 12 is wall.
	// Сейчас мы на 11,10. Следующий шаг вправо -> 12,10 (Wall).
	service.MovePlayer(entity.Right)
	if service.GetPlayer().Pos != expected { // Должен остаться на месте
		t.Errorf("Player moved into wall! Pos: %v", service.GetPlayer().Pos)
	}

	// Проверка сообщения
	msgs := service.GetMessages()
	found := false
	for _, m := range msgs {
		if m == "You tried to move wall" {
			found = true
		}
	}
	if !found {
		t.Error("Expected 'You tried to move wall' message")
	}
}

func TestPlayerAttack(t *testing.T) {
	service := setupService()
	playerPos := entity.Point{X: 10, Y: 10}
	monsterPos := entity.Point{X: 11, Y: 10}

	service.state.CurrentLevel = createTestLevel(playerPos)
	service.GetPlayer().Pos = playerPos

	// Создаем монстра (Зомби)
	monster := service.repo.CreateMonsterWithBaseCharacteristic(entity.Zombie)
	monster.Pos = monsterPos
	// Гарантируем попадание для теста (высокая ловкость игрока)
	service.GetPlayer().Dexterity = 100
	monster.Dexterity = 0
	monster.Health = 10 // Мало ХП, чтобы убить быстро

	service.state.CurrentLevel.Monsters[monsterPos] = monster

	// Атакуем (двигаемся на монстра)
	service.MovePlayer(entity.Right)

	// Игрок не должен сдвинуться
	if service.GetPlayer().Pos == monsterPos {
		t.Error("Player should not move into monster tile during attack")
	}

	// Проверяем урон
	if monster.Health >= 10 {
		t.Error("Monster should take damage")
	}
}

func TestVampireFirstHitMiss(t *testing.T) {
	service := setupService()
	playerPos := entity.Point{X: 10, Y: 10}
	monsterPos := entity.Point{X: 11, Y: 10}

	service.state.CurrentLevel = createTestLevel(playerPos)
	service.GetPlayer().Pos = playerPos
	service.GetPlayer().Dexterity = 100 // 100% шанс попадания обычно

	vampire := service.repo.CreateMonsterWithBaseCharacteristic(entity.Vampire)
	vampire.Pos = monsterPos
	vampire.WasHit = false
	service.state.CurrentLevel.Monsters[monsterPos] = vampire

	// Первый удар
	service.MovePlayer(entity.Right)

	// Должен быть промах
	if vampire.Health < vampire.MaxHealth {
		t.Error("Vampire should dodge the first hit")
	}
	if !vampire.WasHit {
		t.Error("Vampire WasHit flag should be set to true")
	}

	// Второй удар
	service.MovePlayer(entity.Right)
	if vampire.Health >= vampire.MaxHealth {
		t.Error("Vampire should take damage on second hit")
	}
}

func TestPickupAndUseFood(t *testing.T) {
	service := setupService()
	service.GetPlayer().Health = 50
	service.GetPlayer().MaxHealth = 100

	// Создаем еду
	food := &entity.Food{
		BaseItem:   entity.BaseItem{Name: "Apple"},
		HealAmount: 20,
	}

	// Кладем в инвентарь (через логику подбора или напрямую)
	// Тестируем UseItem, так что положим напрямую
	service.state.Player.Backpack.Items[entity.ItemFood] = append(service.state.Player.Backpack.Items[entity.ItemFood], food)

	// Используем (индекс 1, т.к. пользовательский ввод 1-based)
	err := service.UseItem(entity.ItemFood, 1)
	if err != nil {
		t.Fatalf("UseItem returned error: %v", err)
	}

	if service.GetPlayer().Health != 70 {
		t.Errorf("Expected HP 70, got %d", service.GetPlayer().Health)
	}

	if len(service.GetPlayer().Backpack.Items[entity.ItemFood]) != 0 {
		t.Error("Food should be consumed (removed from backpack)")
	}
}

func TestDoorsAndKeys(t *testing.T) {
	service := setupService()
	pPos := entity.Point{X: 10, Y: 10}
	doorPos := entity.Point{X: 11, Y: 10}

	service.state.CurrentLevel = createTestLevel(pPos)
	service.GetPlayer().Pos = pPos

	// Ставим дверь
	service.state.CurrentLevel.SetTile(doorPos.X, doorPos.Y, entity.TileDoorLocked)
	door := &entity.Door{Object: entity.Object{Pos: doorPos}, Color: entity.RedKey, Opened: false}
	service.state.CurrentLevel.Doors[doorPos] = door

	// 1. Попытка пройти без ключа
	service.MovePlayer(entity.Right)
	if service.GetPlayer().Pos == doorPos {
		t.Error("Should not pass locked door without key")
	}

	// 2. Даем ключ
	key := &entity.Key{BaseItem: entity.BaseItem{Name: "Red Key"}, Color: entity.RedKey}
	service.state.Player.Backpack.Items[entity.ItemKey] = append(service.state.Player.Backpack.Items[entity.ItemKey], key)

	// 3. Попытка пройти с ключом
	service.MovePlayer(entity.Right)
	if service.GetPlayer().Pos != doorPos {
		t.Error("Should open door and move into it with key")
	}
	if !door.Opened {
		t.Error("Door should be marked as Opened")
	}
}

func TestOgreMechanics(t *testing.T) {
	service := setupService()
	// Игрок рядом с огром
	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}
	ogrePos := entity.Point{X: 11, Y: 10}

	service.state.CurrentLevel = createTestLevel(service.GetPlayer().Pos)

	ogre := service.repo.CreateMonsterWithBaseCharacteristic(entity.Ogre)
	ogre.Pos = ogrePos
	ogre.Dexterity = 100 // Чтобы попадал по игроку
	service.state.CurrentLevel.Monsters[ogrePos] = ogre

	// Ход игрока (WaitTurn вызывает UpdateWorld, где ходят монстры)
	// Огр атакует
	initialHP := service.GetPlayer().Health
	service.WaitTurn()

	if service.GetPlayer().Health >= initialHP {
		t.Error("Ogre should have hit player")
	}

	// Теперь Огр должен отдыхать (IsResting = true)
	if !ogre.IsResting {
		t.Error("Ogre should be resting after attack")
	}

	// Следующий ход - Огр пропускает
	hpAfterFirstHit := service.GetPlayer().Health
	service.WaitTurn()

	if service.GetPlayer().Health < hpAfterFirstHit {
		t.Error("Ogre should skip turn (resting), no damage expected")
	}
	if ogre.IsResting {
		t.Error("Ogre should wake up after resting turn")
	}
}

func TestStatsTracking(t *testing.T) {
	service := setupService()
	service.state.TotalStats.StepsTaken = 0

	service.state.CurrentLevel = createTestLevel(entity.Point{X: 10, Y: 10})
	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}

	// Делаем шаг
	service.MovePlayer(entity.Right) // 10,10 -> 11,10

	if service.state.TotalStats.StepsTaken != 1 {
		t.Errorf("Expected 1 step taken, got %d", service.state.TotalStats.StepsTaken)
	}
}

func TestLevelTransition(t *testing.T) {
	service := setupService()
	// Устанавливаем текущий уровень 1
	service.state.CurrentLevel.LevelNum = 1

	// Создаем лестницу
	stairsPos := entity.Point{X: 11, Y: 10}
	service.state.CurrentLevel.SetTile(stairsPos.X, stairsPos.Y, entity.TileDownStairs)

	// Телепортируем игрока к лестнице и заставляем наступить на нее
	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}
	service.MovePlayer(entity.Right) // Идем на 11,10

	// Проверяем:
	// 1. Уровень должен повыситься
	if service.GetCurrentLevel().LevelNum != 2 {
		t.Errorf("Expected Level 2, got %d", service.GetCurrentLevel().LevelNum)
	}
	// 2. Игрок должен быть в новой стартовой позиции (обычно это генерируется, но проверим, что он не на лестнице старого уровня)
	if service.GetPlayer().Pos == stairsPos {
		t.Error("Player should be moved to new start position, but stayed at stairs coordinates")
	}
}

func TestGameOver(t *testing.T) {
	service := setupService()
	service.GetPlayer().Health = 5

	// Создаем сильного монстра рядом
	monster := service.repo.CreateMonsterWithBaseCharacteristic(entity.Ogre)
	monster.Strength = 100 // Убивает с удара
	monster.Dexterity = 100
	monster.Pos = entity.Point{X: 11, Y: 10}
	service.state.CurrentLevel.Monsters[monster.Pos] = monster

	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}

	// Пропускаем ход, чтобы монстр ударил
	service.WaitTurn()

	if service.GetPlayer().Health > 0 {
		t.Fatal("Player should be dead")
	}

	if service.GetStatus() != entity.StatusGameOver {
		t.Errorf("Expected StatusGameOver, got %v", service.GetStatus())
	}
}

func TestScrollsAndElixirs(t *testing.T) {
	service := setupService()
	player := service.GetPlayer()
	initialStr := player.Strength

	// --- ТЕСТ СВИТКА (Scroll) ---
	scroll := &entity.Scroll{
		BaseItem: entity.BaseItem{Name: "Scroll of Power"},
		Stat:     entity.StatStrength,
		Increase: 5,
	}
	player.Backpack.Items[entity.ItemScroll] = append(player.Backpack.Items[entity.ItemScroll], scroll)

	// Используем свиток
	service.UseItem(entity.ItemScroll, 1)

	if player.Strength != initialStr+5 {
		t.Errorf("Scroll should permanently increase Strength. Got %d, expected %d", player.Strength, initialStr+5)
	}

	// --- ТЕСТ ЭЛИКСИРА (Elixir) ---
	elixir := &entity.Elixir{
		BaseItem: entity.BaseItem{Name: "Strength Potion"},
		Stat:     entity.StatStrength,
		Increase: 10,
		Duration: 2, // Действует 2 хода
	}
	player.Backpack.Items[entity.ItemElixir] = append(player.Backpack.Items[entity.ItemElixir], elixir)

	// Используем эликсир
	service.UseItem(entity.ItemElixir, 1)

	// Сила должна вырасти: База(была+5) + Эликсир(10)
	expectedStrWithBuff := initialStr + 5 + 10
	if player.Strength != expectedStrWithBuff {
		t.Errorf("Elixir should temporarily increase Strength. Got %d, expected %d", player.Strength, expectedStrWithBuff)
	}

	// Проверяем истечение времени баффа
	// Ход 1
	service.WaitTurn()
	if player.Strength != expectedStrWithBuff {
		t.Error("Buff should still be active")
	}
	// Ход 2 (здесь бафф должен истечь)
	service.WaitTurn()

	// Проверяем, что сила вернулась к значению после свитка
	if player.Strength != initialStr+5 {
		t.Errorf("Buff should expire. Got %d, expected %d", player.Strength, initialStr+5)
	}
}

func TestEquipAndDropWeapon(t *testing.T) {
	service := setupService()
	player := service.GetPlayer()

	// Создаем меч
	sword := &entity.Weapon{
		BaseItem:      entity.BaseItem{Name: "Iron Sword"},
		StrengthBonus: 10,
	}
	player.Backpack.Items[entity.ItemWeapon] = append(player.Backpack.Items[entity.ItemWeapon], sword)

	// 1. Экипировка (EquipWeapon)
	// Индекс 1 (1-based)
	service.EquipWeapon(1)

	if player.Weapon != sword {
		t.Error("Weapon was not equipped")
	}

	// 2. Выбрасывание (DropItem)
	// Сначала снимем (EquipWeapon с индексом 0 снимает, но тут мы тестируем DropItem конкретно)
	// В задании сказано: "DropItem(itemType, index)"

	// Проверим, есть ли место куда бросить (вокруг игрока должны быть свободные клетки)
	service.state.CurrentLevel = createTestLevel(player.Pos)

	// Бросаем меч (он все еще в рюкзаке под индексом 1, даже если надет)
	err := service.DropItem(entity.ItemWeapon, 1)
	if err != nil {
		t.Fatalf("Failed to drop item: %v", err)
	}

	if len(player.Backpack.Items[entity.ItemWeapon]) != 0 {
		t.Error("Weapon should be removed from backpack")
	}

	// Проверяем, что меч лежит на полу рядом
	found := false
	for _, items := range service.GetCurrentLevel().Items {
		for _, item := range items {
			if item.GetName() == "Iron Sword" {
				found = true
			}
		}
	}
	if !found {
		t.Error("Dropped item not found on the floor")
	}
}

func TestSnakeMageSleep(t *testing.T) {
	// Сложный тест, так как там random. Но мы можем проверить саму логику наложения эффекта,
	// если "повезет" или если мы зациклим попытки.

	service := setupService()
	player := service.GetPlayer()

	snake := service.repo.CreateMonsterWithBaseCharacteristic(entity.SnakeMage)
	snake.Pos = entity.Point{X: 11, Y: 10}
	snake.Dexterity = 100 // Всегда попадает

	service.state.CurrentLevel.Monsters[snake.Pos] = snake
	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}

	// Атакуем много раз (пропускаем ходы), пока не уснем
	asleep := false
	for i := 0; i < 50; i++ {
		service.WaitTurn() // Монстр атакует
		if player.SleepTurns > 0 {
			asleep = true
			break
		}
		// Лечимся, чтобы не умереть в тесте
		player.Health = player.MaxHealth
	}

	if !asleep {
		t.Log("Warning: Snake Mage did not put player to sleep in 50 hits. Possible RNG issue or low chance.")
	} else {
		// Проверяем, что во сне нельзя ходить
		// MovePlayer содержит проверку: if SleepTurns > 0 { SleepTurns--; return }

		initialPos := player.Pos
		turns := player.SleepTurns

		service.MovePlayer(entity.Right)

		if player.Pos != initialPos {
			t.Error("Player moved while sleeping!")
		}
		if player.SleepTurns != turns-1 {
			t.Error("Sleep turns should decrement after trying to move")
		}
	}
}

func TestInventoryLimit(t *testing.T) {
	service := setupService()
	player := service.GetPlayer()

	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}
	service.state.CurrentLevel = createTestLevel(service.GetPlayer().Pos)
	service.state.CurrentLevel.SetTile(11, 10, entity.TileFloor)

	// Заполняем инвентарь едой (9 штук)
	for i := 0; i < entity.MaxItemsInBackpackOfOneType; i++ {
		player.Backpack.Items[entity.ItemFood] = append(player.Backpack.Items[entity.ItemFood], &entity.Food{BaseItem: entity.BaseItem{Name: "Apple"}})
	}

	// Создаем 10-е яблоко на полу
	itemPos := entity.Point{X: 11, Y: 10}
	service.state.CurrentLevel.Items[itemPos] = []entity.Item{
		&entity.Food{BaseItem: entity.BaseItem{Name: "Extra Apple"}, HealAmount: 10},
	}

	// Пытаемся наступить (подобрать)
	service.MovePlayer(entity.Right)

	// Проверки
	if len(player.Backpack.Items[entity.ItemFood]) != 9 {
		t.Errorf("Inventory should not exceed limit. Got %d items", len(player.Backpack.Items[entity.ItemFood]))
	}

	// Предмет должен остаться на уровне
	found := false
	if len(service.GetCurrentLevel().Items[itemPos]) > 0 {
		found = true
	}

	if !found {
		t.Error("Item should remain on floor if inventory is full")
	}

	// Проверяем сообщение
	msgs := service.GetMessages()
	msgFound := false
	for _, m := range msgs {
		if m == "You have not enough space for Extra Apple" {
			msgFound = true
		}
	}
	if !msgFound {
		t.Error("Expected warning message about full inventory")
	}
}

func TestGameVictory(t *testing.T) {
	service := setupService()
	// Эмулируем нахождение на 21 уровне (MaxLevel)
	// В service.go есть константа MaxLevel (обычно 21)
	service.state.CurrentLevel.LevelNum = entity.MaxLevel

	stairsPos := entity.Point{X: 11, Y: 10}
	service.state.CurrentLevel.SetTile(stairsPos.X, stairsPos.Y, entity.TileDownStairs)

	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}
	service.MovePlayer(entity.Right) // Идем на лестницу

	if service.GetStatus() != entity.StatusVictory {
		t.Errorf("Expected StatusVictory when finishing max level, got %v", service.GetStatus())
	}
}

func TestGhostTeleport(t *testing.T) {
	service := setupService()

	// 1. Создаем контролируемую среду (стены везде, кроме тестовой зоны)
	// Создаем пол 5x5 в области 10,10
	service.state.CurrentLevel = createTestLevel(entity.Point{X: 10, Y: 10}) // Это создаст пол 3x3
	// Расширим пол до 5x5 для комнаты (10,10 - 14,14)
	for y := 10; y < 15; y++ {
		for x := 10; x < 15; x++ {
			service.state.CurrentLevel.SetTile(x, y, entity.TileFloor)
		}
	}

	room := &entity.Room{
		Object: entity.Object{
			Pos:  entity.Point{X: 10, Y: 10},
			Size: entity.Size{Width: 5, Height: 5},
		},
	}
	service.state.CurrentLevel.Rooms = append(service.state.CurrentLevel.Rooms, room)

	// 3. Помещаем игрока и призрака на ВАЛИДНЫЕ клетки (пол)
	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}

	ghost := service.repo.CreateMonsterWithBaseCharacteristic(entity.Ghost)
	ghost.Pos = entity.Point{X: 14, Y: 14} // В том же углу комнаты

	ghost.Hostility = 0

	service.state.CurrentLevel.Monsters[ghost.Pos] = ghost

	// Ход
	service.WaitTurn()

	// Проверка
	if ghost.Pos == (entity.Point{X: 14, Y: 14}) {
		t.Log("Ghost stayed in place (Random chance or logic bug)")
	} else {
		// Проверяем, что он все еще в комнате (на полу)
		tile := service.GetCurrentLevel().GetTile(ghost.Pos.X, ghost.Pos.Y)
		if tile != entity.TileFloor {
			t.Errorf("Ghost teleported into a wall/void at %v", ghost.Pos)
		}
	}
}

func TestDynamicDifficulty(t *testing.T) {
	service := setupService()
	// Устанавливаем условия для ПОВЫШЕНИЯ сложности (HARD)
	// 1. HP полное
	service.GetPlayer().Health = service.GetPlayer().MaxHealth
	// 2. Убито много врагов
	service.state.LevelStats.EnemiesDefeated = 100
	// 3. Не пил эликсиров и не ел (уже 0 по дефолту)
	service.state.LevelStats.AmountConsumedFood = 0

	currentDiff := service.state.BalanceDifficulty

	// Триггерим переход уровня, где срабатывает adjustBalanceAfterLevel
	service.state.CurrentLevel.LevelNum = 1
	stairsPos := entity.Point{X: 11, Y: 10}
	service.state.CurrentLevel.SetTile(stairsPos.X, stairsPos.Y, entity.TileDownStairs)

	service.GetPlayer().Pos = entity.Point{X: 10, Y: 10}
	service.MovePlayer(entity.Right)

	if service.state.BalanceDifficulty <= currentDiff {
		t.Errorf("Difficulty should increase for strong player. Was %d, became %d", currentDiff, service.state.BalanceDifficulty)
	}
}

func TestPlayerRotation(t *testing.T) {
	service := setupService()
	service.GetPlayer().Facing = entity.Top

	service.RotatePlayer(entity.Right)

	if service.GetPlayer().Facing != entity.Right {
		t.Error("RotatePlayer failed to update facing direction")
	}
}
