package render

import (
	"github.com/gdamore/tcell/v2"
	"example.com/rogue/domain/entity"
)

func (r *Render) drawMonsters() {
	level := r.game.GetCurrentLevel()
	player := r.game.GetPlayer()
	if level == nil || player == nil {
		return
	}

	for _, monster := range level.Monsters {
		dx := monster.Pos.X - player.Pos.X
		dy := monster.Pos.Y - player.Pos.Y

		screenX := r.width/2 + dx
		screenY := (r.height-9)/2 + 3 + dy

		if screenX < 0 || screenX >= r.width || screenY < 3 || screenY >= r.height-6 {
			continue
		}

		// Fog of war
		if r.fogOfWar != nil && !r.fogOfWar.IsVisible(int(monster.Pos.X), int(monster.Pos.Y)) {
			continue
		}

		// не рисуем монстра на игроке
		if monster.Pos.X == player.Pos.X && monster.Pos.Y == player.Pos.Y {
			continue
		}
		if monster.Type == entity.Ghost && !monster.IsVisible {
			continue
		}

		monsterInfo := GetMonsterInfo(monster.Type)

		if monster.Type == entity.Mimic && monster.IsDisguised {
			monsterInfo = GetItemInfoByType(monster.DisguiseItem)
		}

		style := tcell.StyleDefault.Foreground(monsterInfo.Color)
		if monsterInfo.Bold {
			style = style.Bold(true)
		}

		r.screen.SetContent(screenX, screenY, monsterInfo.Symbol, nil, style)
	}
}

func (r *Render) drawItems() {
	level := r.game.GetCurrentLevel()
	player := r.game.GetPlayer()
	if level == nil || player == nil {
		return
	}

	for _, items := range level.Items {
		if len(items) == 0 {
			continue
		}

		it := items[len(items)-1]
		pos := it.GetPosition()

		dx := pos.X - player.Pos.X
		dy := pos.Y - player.Pos.Y

		screenX := r.width/2 + dx
		screenY := (r.height-9)/2 + 3 + dy

		if screenX < 0 || screenX >= r.width || screenY < 3 || screenY >= r.height-6 {
			continue
		}

		// Fog of war
		if r.fogOfWar != nil && !r.fogOfWar.IsVisible(int(pos.X), int(pos.Y)) {
			// предметы в тумане не рисуем
			continue
		}

		tileInfo := GetItemInfo(it)
		style := tcell.StyleDefault.Foreground(tileInfo.Color)
		if tileInfo.Bold {
			style = style.Bold(true)
		}
		r.screen.SetContent(screenX, screenY, tileInfo.Symbol, nil, style)
	}
}



