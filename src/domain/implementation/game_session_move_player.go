package implementation

import (
	"math/rand"
	"strconv"

	"example.com/rogue/domain/entity"
)

func (g *GameService) MovePlayer(dir entity.Direction) {
	g.clearMessages()

	if g.GetPlayer().SleepTurns > 0 {
		g.GetPlayer().SleepTurns--
	}

	destination := calculateDestination(g.GetPlayer().Pos, dir)

	if g.GetCurrentLevel().Monsters[destination] != nil {
		g.playerAttack(g.GetCurrentLevel().Monsters[destination])
		g.UpdateWorld()
	} else if len(g.GetCurrentLevel().Items[destination]) != 0 {
		g.pickUpItems(destination)
		g.UpdateWorld()
	} else if g.GetCurrentLevel().GetTile(destination.X, destination.Y) == entity.TileWall {
		g.logMessage("You tried to move wall")
	} else if g.GetCurrentLevel().GetTile(destination.X, destination.Y) == entity.TileDownStairs {
		g.nextLevel(g.GetCurrentLevel().LevelNum, g.state.BalanceDifficulty)
	} else if g.GetCurrentLevel().GetTile(destination.X, destination.Y) == entity.TileDoorLocked {
		g.tryGoThroughDoor(destination)
	} else {
		g.GetPlayer().Pos = destination
		g.UpdateWorld()
	}
}

func (g *GameService) playerAttack(monster *entity.Monster) {
	hit := g.hitCheck(monster)

	if !hit {
		return
	}

	g.damageMonster(monster)

	if monster.Health <= 0 {
		g.defeatMonster(monster)
	}
}

func (g *GameService) hitCheck(monster *entity.Monster) bool {
	if monster.Type == entity.Vampire && !monster.WasHit {
		g.logMessage("You miss the Vampire! (He is too fast for the first strike)")
		monster.WasHit = true
		return false
	}

	hitChance := calculateHitChance(g.GetPlayer().CharacterStats, monster.CharacterStats)

	roll := rand.Intn(100)
	if roll >= hitChance {
		g.logMessage("You miss the " + monster.Type.String() + ".")
		return false
	}

	return true
}

func (g *GameService) damageMonster(monster *entity.Monster) {
	player := g.state.Player
	damage := player.Strength
	if player.Weapon != nil {
		damage += player.Weapon.StrengthBonus
	}

	g.writeTotalHitsStat(damage)

	monster.Health -= damage
	g.logMessage("You hit the " + monster.Type.String() + " for " + strconv.Itoa(damage) + " damage.")
}

func (g *GameService) defeatMonster(monster *entity.Monster) {
	g.logMessage("The " + monster.Type.String() + " dies!")

	treasureValue := int(float64(monster.MaxHealth+monster.Strength+monster.Dexterity+int(monster.Hostility)*5) / 2 * g.treasureMultiplier())

	gold := &entity.Treasure{
		BaseItem: entity.BaseItem{
			Name: "Gold Coin",
			Pos:  monster.Pos,
		},
		Value: treasureValue,
	}

	g.state.CurrentLevel.Items[monster.Pos] = append(g.state.CurrentLevel.Items[monster.Pos], gold)
	g.logMessage("It dropped " + strconv.Itoa(treasureValue) + " gold.")
	g.writeDefeatedEnemyStat()

	delete(g.state.CurrentLevel.Monsters, monster.Pos)
}

func (g *GameService) treasureMultiplier() float64 {
	k := 1.0 + 0.1*float64(g.state.CurrentLevel.LevelNum+g.state.BalanceDifficulty-1)
	if k > 4 {
		k = 4
	}
	return k
}

func (g *GameService) pickUpItems(destination entity.Point) {
	newItems := make([]entity.Item, 0, len(g.GetCurrentLevel().Items[destination]))

	for _, item := range g.GetCurrentLevel().Items[destination] {
		picked := tryPickupItem(g.GetPlayer(), item)
		if picked {
			g.logMessage("You picked up " + item.GetName())
		} else {
			newItems = append(newItems, item)
			g.logMessage("You have not enough space for " + item.GetName())
		}
	}
	g.state.CurrentLevel.Items[destination] = newItems
	g.GetPlayer().Pos = destination
}

func tryPickupItem(player *entity.Player, item entity.Item) bool {
	if item.GetType() == entity.ItemTreasure {
		player.Backpack.Gold += item.(*entity.Treasure).Value
		return true
	}

	if len(player.Backpack.Items[item.GetType()]) == entity.MaxItemsInBackpackOfOneType {
		return false
	}

	player.Backpack.Items[item.GetType()] = append(player.Backpack.Items[item.GetType()], item)

	return true
}

func (g *GameService) tryGoThroughDoor(destination entity.Point) {
	door := g.GetCurrentLevel().Doors[destination]
	if hasKey(g.GetPlayer(), door) {
		g.logMessage("You opened the door")
		door.Opened = true
		g.state.CurrentLevel.SetTile(door.Pos.X, door.Pos.Y, entity.TileDoorOpened)
		g.GetPlayer().Pos = destination
		g.UpdateWorld()
	} else {
		g.logMessage("You need a key to open the door")
	}
}

func hasKey(player *entity.Player, door *entity.Door) bool {
	for _, key := range player.Backpack.Items[entity.ItemKey] {
		if key.(*entity.Key).Color == door.Color {
			return true
		}
	}

	return false
}
