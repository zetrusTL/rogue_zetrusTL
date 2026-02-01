package implementation

import (
	"math/rand"
	"strconv"

	"example.com/rogue/domain/entity"
)

func (g *GameService) UpdateWorld() *entity.Level {
	g.updateBuffs()

	g.writeStepsTakenStat()

	g.updateMonsters()

	if g.GetPlayer().Health <= 0 {
		g.endGame(false)
	}

	return g.state.CurrentLevel
}

func (g *GameService) updateBuffs() {
	player := g.state.Player

	activeBuffs := make([]*entity.Buff, 0)
	for _, buff := range player.ActiveBuffs {
		if buff.TurnsLeft <= 0 {
			g.removeBuffEffect(player, buff)
			g.logMessage("Effect of " + buff.Stat.String() + " expired.")
		} else {
			activeBuffs = append(activeBuffs, buff)
			buff.TurnsLeft--
		}
	}
	player.ActiveBuffs = activeBuffs
}

func (g *GameService) removeBuffEffect(p *entity.Player, buff *entity.Buff) {
	switch buff.Stat {
	case entity.StatStrength:
		p.Strength -= buff.Increase
	case entity.StatDexterity:
		p.Dexterity -= buff.Increase
	case entity.StatMaxHealth:
		p.MaxHealth -= buff.Increase
		p.Health -= buff.Increase

		if p.Health <= 0 {
			p.Health = 1
			g.logMessage("The magic fades, leaving you barely alive...")
		}
	}
}

func (g *GameService) updateMonsters() {
	// Копируем список, чтобы безопасно менять карту позиций
	monsters := make([]*entity.Monster, 0, len(g.state.CurrentLevel.Monsters))
	for _, m := range g.state.CurrentLevel.Monsters {
		monsters = append(monsters, m)
	}

	for _, m := range monsters {
		// Ogre Logic: Пропуск хода (отдых после атаки)
		if m.Type == entity.Ogre && m.IsResting {
			m.IsResting = false
			g.logMessage("The Ogre catches its breath.")
			continue
		}

		g.processMonsterTurn(m)
	}
}

func (g *GameService) processMonsterTurn(monster *entity.Monster) {
	// Очки действия (Action Points)
	// Ogre двигается 2 клетки за ход (или Движение + Атака)
	actions := 1
	if monster.Type == entity.Ogre {
		actions = 2
	}

	for i := 0; i < actions; i++ {
		if _, ok := g.state.CurrentLevel.Monsters[monster.Pos]; !ok {
			break
		}

		// Если Огр атаковал на первом действии, он устает и не делает второе действие
		if monster.Type == entity.Ogre && monster.IsResting {
			break
		}

		dist := distance(monster.Pos, g.GetPlayer().Pos)

		if monster.Type == entity.Mimic && dist <= float64(monster.Hostility) {
			monster.IsDisguised = false
		}

		if dist < 1.5 {
			if monster.Type == entity.Ghost {
				monster.IsVisible = true
			}
			g.monsterAttack(monster)
			continue
		}

		if monster.Type == entity.Mimic && monster.IsDisguised {
			//мимик замаскировался
			break
		}

		moved := g.monsterMove(dist, monster)
		if !moved {
			// Если не смог походить (зажат), теряет действие
			break
		}

		if monster.Type == entity.Ghost && dist >= 1.5 {
			if rand.Intn(100) < 50 {
				monster.IsVisible = !monster.IsVisible
			}
		}
	}
}

func (g *GameService) monsterAttack(monster *entity.Monster) {
	player := g.state.Player

	hitChance := calculateHitChance(monster.CharacterStats, player.CharacterStats)

	roll := rand.Intn(100)
	// Ogre: Rest after attack
	if monster.Type == entity.Ogre {
		monster.IsResting = true
	}
	if roll >= hitChance {
		g.logMessage("The " + monster.Type.String() + " misses you.")
		return
	}

	damage := monster.Strength

	g.writeTotalHitsStat(damage)

	player.Health -= damage
	g.logMessage("The " + monster.Type.String() + " hits you for " + strconv.Itoa(damage) + " damage!")

	// Vampire: Drains Max HP
	if monster.Type == entity.Vampire {
		player.MaxHealth -= entity.VampireMaxHealthDrained
		if player.MaxHealth < 1 {
			player.MaxHealth = 1
		}
		if player.Health > player.MaxHealth {
			player.Health = player.MaxHealth
		}
		g.logMessage("You feel weaker (Max HP drained)!")
	}

	// Snake-Mage: Sleep Chance
	if monster.Type == entity.SnakeMage {
		if rand.Float64() < entity.SnakeAttackSleepChance {
			player.SleepTurns += 1
			g.logMessage("You fall asleep from the venom!")
		}
	}
}

func (g *GameService) monsterMove(dist float64, monster *entity.Monster) bool {
	var nextPos entity.Point
	var hasMove bool

	// Особая логика движения Призрака (Ghost) - Телепортация
	if monster.Type == entity.Ghost {
		nextPos, hasMove = g.moveGhost(dist, monster)
	} else {
		nextPos, hasMove = g.moveMonster(dist, monster)
	}

	// Выполнение движения
	if hasMove {
		delete(g.state.CurrentLevel.Monsters, monster.Pos)

		monster.Pos = nextPos

		g.state.CurrentLevel.Monsters[monster.Pos] = monster
		return true
	}

	return false
}

func (g *GameService) moveMonster(dist float64, monster *entity.Monster) (nextPos entity.Point, hasMove bool) {
	isChasing := dist <= float64(monster.Hostility)

	// Попытка найти путь к игроку (Погоня)
	if isChasing {
		next := g.findPathBFS(monster.Pos, g.GetPlayer().Pos)
		if next != nil {
			nextPos = *next
			hasMove = true
		}
	}

	// Если не преследует или пути нет -> Блуждание
	if !hasMove {
		nextPos, hasMove = g.calculateNextPos(monster)
	}

	return nextPos, hasMove
}

func (g *GameService) calculateNextPos(monster *entity.Monster) (nextPos entity.Point, hasMove bool) {
	neighbors := g.getValidNeighborsForMonster(monster.Pos, false)

	// Snake-Mage Wandering: Только диагонали
	if monster.Type == entity.SnakeMage {
		diagNeighbors := make([]entity.Point, 0)
		for _, n := range neighbors {
			if n.X != monster.Pos.X && n.Y != monster.Pos.Y {
				diagNeighbors = append(diagNeighbors, n)
			}
		}
		neighbors = diagNeighbors
	}

	if len(neighbors) > 0 {
		nextPos = neighbors[rand.Intn(len(neighbors))]
		hasMove = true
	}

	return nextPos, hasMove
}

func (g *GameService) moveGhost(dist float64, monster *entity.Monster) (nextPos entity.Point, hasMove bool) {
	isChasing := dist <= float64(monster.Hostility)

	var candidates []entity.Point

	if isChasing {
		// При погоне пытается прыгнуть ближе к игроку
		// Собираем валидные точки в радиусе 3 от игрока
		pPos := g.GetPlayer().Pos
		for dy := -3; dy <= 3; dy++ {
			for dx := -3; dx <= 3; dx++ {
				p := entity.Point{X: pPos.X + dx, Y: pPos.Y + dy}
				if g.monsterCanMoveTo(p, true) && p != pPos {
					candidates = append(candidates, p)
				}
			}
		}
	} else {
		candidates = g.ghostRandomMove(candidates, monster)
	}

	if len(candidates) > 0 {
		target := candidates[rand.Intn(len(candidates))]
		nextPos = target
		hasMove = true
	}

	return nextPos, hasMove
}

func (g *GameService) ghostRandomMove(candidates []entity.Point, monster *entity.Monster) []entity.Point {
	// Определяем, в какой комнате монстр (если в комнате)
	var currentRoom *entity.Room
	for _, room := range g.GetCurrentLevel().Rooms {
		if isPointInRoom(monster.Pos, room) {
			currentRoom = room
			break
		}
	}

	if currentRoom != nil {
		candidates = g.ghostTeleport(currentRoom, candidates)
	} else {
		// Если в коридоре - прыжок на 1-2 клетки случайно
		for dy := -2; dy <= 2; dy++ {
			for dx := -2; dx <= 2; dx++ {
				p := entity.Point{X: monster.Pos.X + dx, Y: monster.Pos.Y + dy}
				if g.monsterCanMoveTo(p, false) {
					candidates = append(candidates, p)
				}
			}
		}
	}

	return candidates
}

func (g *GameService) ghostTeleport(currentRoom *entity.Room, candidates []entity.Point) []entity.Point {
	// Телепорт в любую точку этой комнаты
	// Генерируем несколько случайных попыток для эффективности
	for i := 0; i < 10; i++ {
		rx := currentRoom.Pos.X + rand.Intn(currentRoom.Size.Width)
		ry := currentRoom.Pos.Y + rand.Intn(currentRoom.Size.Height)
		p := entity.Point{X: rx, Y: ry}
		if g.monsterCanMoveTo(p, false) {
			candidates = append(candidates, p)
		}
	}

	return candidates
}

// --- Helpers ---

func (g *GameService) findPathBFS(start, target entity.Point) *entity.Point {
	queue := []entity.Point{start}
	cameFrom := make(map[entity.Point]entity.Point)
	visited := make(map[entity.Point]bool)
	visited[start] = true

	found := false

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == target {
			found = true
			break
		}

		points := g.getValidNeighborsForMonster(current, true)

		for _, next := range points {
			if !visited[next] {
				visited[next] = true
				cameFrom[next] = current
				queue = append(queue, next)
			}
		}
	}

	if !found {
		return nil
	}

	// Восстановление пути
	var path []entity.Point
	curr := target
	for curr != start {
		path = append([]entity.Point{curr}, path...)
		curr = cameFrom[curr]
	}
	return &path[0]
}

func isPointInRoom(p entity.Point, r *entity.Room) bool {
	return p.X >= r.Pos.X && p.X < r.Pos.X+r.Size.Width &&
		p.Y >= r.Pos.Y && p.Y < r.Pos.Y+r.Size.Height
}
