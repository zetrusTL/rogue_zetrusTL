package data

import (
	"os"
	"path/filepath"
	"testing"

	"example.com/rogue/domain/entity"
)

// Helper to create a repo with temp files
func setupTestRepo(t *testing.T) (*JsonGameRepository, string) {
	tmpDir, err := os.MkdirTemp("", "rogue_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	repo := &JsonGameRepository{
		pathGameState: filepath.Join(tmpDir, "game_state.json"),
		pathStats:     filepath.Join(tmpDir, "statistic.json"),
		// Registry поля оставляем nil, так как тестируем только сохранение/загрузку
	}

	return repo, tmpDir
}

func cleanupTestRepo(t *testing.T, dir string) {
	os.RemoveAll(dir)
}

func TestSaveAndLoadSession_Basic(t *testing.T) {
	repo, dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	// 1. Создаем простое состояние
	state := &entity.GameState{
		Status:            entity.StatusRunning,
		BalanceDifficulty: 5,
		Player: &entity.Player{
			CharacterStats: entity.CharacterStats{
				Health:    100,
				MaxHealth: 100,
			},
			Backpack: entity.NewBackpack(),
		},
		CurrentLevel: &entity.Level{
			LevelNum: 1,
		},
	}

	// 2. Сохраняем
	err := repo.SaveSession(state)
	if err != nil {
		t.Fatalf("SaveSession failed: %v", err)
	}

	// 3. Проверяем, что файл создан
	if !repo.HasSavedSession() {
		t.Error("HasSavedSession returned false after saving")
	}

	// 4. Загружаем
	loadedState, err := repo.LoadSession()
	if err != nil {
		t.Fatalf("LoadSession failed: %v", err)
	}

	// 5. Сравниваем основные поля
	if loadedState.BalanceDifficulty != state.BalanceDifficulty {
		t.Errorf("Expected difficulty %d, got %d", state.BalanceDifficulty, loadedState.BalanceDifficulty)
	}
	if loadedState.Player.Health != state.Player.Health {
		t.Errorf("Expected player health %d, got %d", state.Player.Health, loadedState.Player.Health)
	}
}

func TestDeleteSession(t *testing.T) {
	repo, dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	// Сохраняем фиктивный файл
	state := &entity.GameState{Status: entity.StatusRunning}
	_ = repo.SaveSession(state)

	if !repo.HasSavedSession() {
		t.Fatal("Setup failed: session not saved")
	}

	// Удаляем
	err := repo.DeleteSession()
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	// Проверяем
	if repo.HasSavedSession() {
		t.Error("Session file should be deleted")
	}
}

func TestLeaderboard(t *testing.T) {
	repo, dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	// 1. Проверяем пустой лидерборд
	lb, err := repo.GetLeaderboard()
	if err != nil {
		t.Fatalf("GetLeaderboard on empty failed: %v", err)
	}
	if len(lb) != 0 {
		t.Error("Expected empty leaderboard")
	}

	// 2. Добавляем результаты (в разнобой по золоту)
	results := []entity.RunResult{
		{PlayerName: "Player1", GoldCollected: 100, DeepestLevel: 2},
		{PlayerName: "Player2", GoldCollected: 500, DeepestLevel: 5}, // Лидер
		{PlayerName: "Player3", GoldCollected: 50, DeepestLevel: 1},
	}

	for _, r := range results {
		err := repo.SaveRunResult(r)
		if err != nil {
			t.Fatalf("SaveRunResult failed: %v", err)
		}
	}

	// 3. Загружаем и проверяем сортировку
	lb, err = repo.GetLeaderboard()
	if err != nil {
		t.Fatalf("GetLeaderboard failed: %v", err)
	}

	if len(lb) != 3 {
		t.Errorf("Expected 3 results, got %d", len(lb))
	}

	// Проверка порядка (должно быть: 500, 100, 50)
	if lb[0].GoldCollected != 500 || lb[1].GoldCollected != 100 || lb[2].GoldCollected != 50 {
		t.Error("Leaderboard is not sorted by gold descending")
	}
}

// Этот тест проверяет сохранение сложных структур (интерфейсов)
// Он может упасть, если JSON Unmarshal некорректно восстановит типы интерфейсов
func TestSaveAndLoadSession_ComplexData(t *testing.T) {
	repo, dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	player := &entity.Player{
		Backpack: entity.NewBackpack(),
	}

	// Добавляем конкретный предмет (Еду)
	apple := &entity.Food{
		BaseItem:   entity.BaseItem{Name: "Apple", Pos: entity.Point{X: 1, Y: 1}},
		HealAmount: 10,
	}
	player.Backpack.Items[entity.ItemFood] = append(player.Backpack.Items[entity.ItemFood], apple)

	state := &entity.GameState{
		Player: player,
	}

	// Сохраняем
	if err := repo.SaveSession(state); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Загружаем
	loaded, err := repo.LoadSession()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Проверяем тип восстановленного предмета
	items := loaded.Player.Backpack.Items[entity.ItemFood]
	if len(items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(items))
	}

	// ВАЖНО: Здесь будет проверка типа.
	// Если json.Unmarshal вернул map[string]interface{}, этот Type Assertion упадет или вернет !ok
	loadedApple, ok := items[0].(*entity.Food)
	if !ok {
		t.Errorf("Failed to cast item back to *entity.Food. Got type: %T", items[0])
		t.Log("HINT: You might need a custom UnmarshalJSON implementation for Backpack to handle interfaces.")
	} else {
		if loadedApple.HealAmount != 10 {
			t.Errorf("Data corruption: expected HealAmount 10, got %d", loadedApple.HealAmount)
		}
	}
}
