package generation

import (
	"math/rand"

	"example.com/rogue/domain/entity"
)

var (
	// Типы монстров для рандома
	enemyTypes = []entity.MonsterType{
		entity.Zombie, entity.Vampire, entity.Ghost,
		entity.Ogre, entity.SnakeMage, entity.Mimic,
	}

	// Типы предметов для рандома
	itemTypes = []entity.ItemType{
		entity.ItemFood,
		entity.ItemElixir,
		entity.ItemScroll,
	}

	mimicItemTypes = []entity.ItemType{
		entity.ItemFood,
		entity.ItemElixir,
		entity.ItemScroll,
		entity.ItemWeapon,
	}
)

func (lg *LevelGenerator) generateEntities() {
	lg.generateEnemies()
	lg.generateItems()
}

func (lg *LevelGenerator) generateEnemies() {
	maxCountEntities := int(entity.MaxEnemiesPerRoom * lg.enemyCountMultiplier())
	enemyStatsMultiplier := lg.enemyCountMultiplier()
	for _, room := range lg.level.Rooms {
		if room == lg.startRoom {
			continue
		}

		enemiesCnt := rand.Intn(maxCountEntities)
		for i := 0; i < enemiesCnt; i++ {
			enemy := lg.createRandomEnemy()
			enemy.Object.Pos = getRandomPlaceInRoom(lg.level, room)
			enemy.MaxHealth = int(float64(enemy.MaxHealth) * enemyStatsMultiplier)
			enemy.Health = enemy.MaxHealth
			enemy.Strength = int(float64(enemy.Strength) * enemyStatsMultiplier)
			enemy.Dexterity = int(float64(enemy.Hostility) * enemyStatsMultiplier)

			lg.level.Monsters[enemy.Pos] = enemy
		}
	}
}

func (lg *LevelGenerator) generateItems() {
	usefulChance := lg.itemUsefulChance()

	for _, room := range lg.level.Rooms {
		for i := 0; i < entity.MaxItemsPerRoom; i++ {
			var item entity.Item
			if rand.Float64() <= 0.125 {
				item = lg.repo.CreateRandomItemWithBaseCharacteristic(entity.ItemWeapon)
			} else if rand.Float64() > usefulChance {
				continue
			} else {
				item = lg.createRandomUsefulItem()
			}

			pos := getRandomPlaceInRoom(lg.level, room)
			item.SetPosition(pos)
			lg.level.Items[pos] = append(lg.level.Items[pos], item)
		}
	}
}

func (lg *LevelGenerator) createRandomEnemy() *entity.Monster {
	mType := enemyTypes[rand.Intn(len(enemyTypes))]

	enemy := lg.repo.CreateMonsterWithBaseCharacteristic(mType)
	if enemy.Type == entity.Mimic {
		enemy.IsDisguised = true
		enemy.DisguiseItem = mimicItemTypes[rand.Intn(len(mimicItemTypes))]
	}
	return enemy
}

func (lg *LevelGenerator) createRandomUsefulItem() entity.Item {
	iType := itemTypes[rand.Intn(len(itemTypes))]

	return lg.repo.CreateRandomItemWithBaseCharacteristic(iType)
}

func (lg *LevelGenerator) enemyCountMultiplier() float64 {
	if lg.difficulty <= 1 {
		return 1
	}

	k := 1.0 + 0.1*float64(lg.difficulty-1)
	if k > 2 {
		k = 2
	}
	return k
}

func (lg *LevelGenerator) enemyStatMultiplier() float64 {
	k := 1.0 + 0.05*float64(lg.difficulty-1)
	if k > 3 {
		k = 3
	}
	return k
}

func (lg *LevelGenerator) itemUsefulChance() float64 {
	base := 1 - 0.04*float64(lg.difficulty-1)
	if base < 0.25 {
		base = 0.25
	}
	return base
}
