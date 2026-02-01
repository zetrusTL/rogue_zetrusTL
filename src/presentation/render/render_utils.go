package render

import (
	"example.com/rogue/domain/entity"
	"example.com/rogue/presentation"
	"github.com/gdamore/tcell/v2"
	"strings"
)

// возвращает символ и цвет для монстра
func GetMonsterInfo(monsterType entity.MonsterType) presentation.TileInfo {
	switch monsterType {
	case entity.Zombie:
		return presentation.TileInfo{Symbol: 'z', Color: tcell.ColorGreen, Bold: true}
	case entity.Vampire:
		return presentation.TileInfo{Symbol: 'v', Color: tcell.ColorRed, Bold: true}
	case entity.Ghost:
		return presentation.TileInfo{Symbol: 'g', Color: tcell.ColorWhite, Bold: true}
	case entity.Ogre:
		return presentation.TileInfo{Symbol: 'O', Color: tcell.ColorYellow, Bold: true}
	case entity.SnakeMage:
		return presentation.TileInfo{Symbol: 's', Color: tcell.ColorPurple, Bold: true}
	case entity.Mimic:
		return presentation.TileInfo{Symbol: 'm', Color: tcell.ColorWhite, Bold: true}
	default:
		return presentation.TileInfo{Symbol: 'M', Color: tcell.ColorRed, Bold: true}
	}
}

// возвращает символ и цвет для предмета
func GetItemInfo(item entity.Item) presentation.TileInfo {
	switch item.GetType() {

	case entity.ItemWeapon:
		return presentation.TileInfo{Symbol: '/', Color: tcell.ColorRed, Bold: true}

	case entity.ItemFood:
		return presentation.TileInfo{Symbol: '%', Color: tcell.ColorGreen, Bold: false}

	case entity.ItemElixir:
		return presentation.TileInfo{Symbol: '!', Color: tcell.ColorBlue, Bold: true}

	case entity.ItemScroll:
		return presentation.TileInfo{Symbol: '?', Color: tcell.ColorAqua, Bold: true}

	case entity.ItemTreasure:
		return presentation.TileInfo{Symbol: '$', Color: tcell.ColorGold, Bold: true}

	case entity.ItemKey:
		if k, ok := item.(*entity.Key); ok {
			switch k.Color {
			case entity.RedKey:
				return presentation.TileInfo{Symbol: 'k', Color: tcell.ColorRed, Bold: true}
			case entity.BlueKey:
				return presentation.TileInfo{Symbol: 'k', Color: tcell.ColorBlue, Bold: true}
			case entity.YellowKey:
				return presentation.TileInfo{Symbol: 'k', Color: tcell.ColorYellow, Bold: true}
			}
		}
		return presentation.TileInfo{Symbol: 'k', Color: tcell.ColorWhite, Bold: true}
	}
	return presentation.TileInfo{Symbol: '?', Color: tcell.ColorWhite, Bold: false}
}

func GetItemInfoByType(t entity.ItemType) presentation.TileInfo {
	switch t {
	case entity.ItemWeapon:
		return presentation.TileInfo{Symbol: '/', Color: tcell.ColorRed, Bold: true}
	case entity.ItemFood:
		return presentation.TileInfo{Symbol: '%', Color: tcell.ColorGreen, Bold: false}
	case entity.ItemElixir:
		return presentation.TileInfo{Symbol: '!', Color: tcell.ColorBlue, Bold: true}
	case entity.ItemScroll:
		return presentation.TileInfo{Symbol: '?', Color: tcell.ColorAqua, Bold: true}
	case entity.ItemTreasure:
		return presentation.TileInfo{Symbol: '$', Color: tcell.ColorGold, Bold: true}
	case entity.ItemKey:
		return presentation.TileInfo{Symbol: 'k', Color: tcell.ColorWhite, Bold: true}
	default:
		return presentation.TileInfo{Symbol: '?', Color: tcell.ColorWhite, Bold: false}
	}
}


func generateHealthBar(current, max float64, width int) string {
	if max <= 0 {
		return strings.Repeat("░", width)
	}

	percent := current / max
	filled := int(float64(width) * percent)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// определяет цвет сообщения по его содержанию
func getMessageColor(msg string) tcell.Color {
	switch {
	case strings.Contains(msg, "наносит") || strings.Contains(msg, "урон"):
		return tcell.ColorRed
	case strings.Contains(msg, "лечит") || strings.Contains(msg, "восстанавливает"):
		return tcell.ColorGreen
	case strings.Contains(msg, "нашел") || strings.Contains(msg, "золото"):
		return tcell.ColorGold
	case strings.Contains(msg, "уровень") || strings.Contains(msg, "спускается"):
		return tcell.ColorYellow
	case strings.Contains(msg, "умер") || strings.Contains(msg, "погиб"):
		return tcell.ColorDarkRed
	default:
		return tcell.ColorSilver
	}
}

// возвращает абсолютное значение для целых чисел
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func GetExitInfo() presentation.TileInfo {
	return presentation.TileInfo{Symbol: '>', Color: tcell.ColorGreen, Bold: true}
}

func GetTileInfoForEntityTile(tile entity.TileType) presentation.TileInfo {
	switch tile {
	case entity.TileDoorOpened:
		return presentation.TileInfo{Symbol: '+', Color: tcell.ColorYellow, Bold: true}

	case entity.TileDoorLocked:
		return presentation.TileInfo{Symbol: '+', Color: tcell.ColorGray, Bold: true}

	case entity.TileDownStairs:
		return presentation.TileInfo{Symbol: '>', Color: tcell.ColorWhite, Bold: true}

	default:
		return presentation.TileInfo{Symbol: '·', Color: tcell.ColorRed, Bold: false}
	}
}
