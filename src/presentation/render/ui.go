package render

import (
	"fmt"
	"unicode/utf8"
	"example.com/rogue/domain/entity"
	"github.com/gdamore/tcell/v2"
)

func (r *Render) drawUI() {
    player := r.game.GetPlayer()
    level := r.game.GetCurrentLevel()

    r.drawBox(0, 0, r.width, 3, "ROGUE GAME", tcell.ColorWhite, true)

    hpPercent := 0.0
    if player.MaxHealth > 0 {
        hpPercent = float64(player.Health) / float64(player.MaxHealth)
    }

hpBar := generateHealthBar(float64(player.Health), float64(player.MaxHealth), 10)

    weaponInfo := ""
    if player.Weapon != nil {
        weaponInfo = fmt.Sprintf("%s", player.Weapon.GetName())
    } else {
        weaponInfo = "Без оружия"
    }
    
    rk, bk, yk := countKeys(player)
    keysText := fmt.Sprintf("Ключи: К:%s С:%s Ж:%s", onOff(rk), onOff(bk), onOff(yk))

    // собираем статус
    levelText := fmt.Sprintf("Ур.%d", level.LevelNum)
    hpText := fmt.Sprintf("HP: %s %d/%d", hpBar, player.Health, player.MaxHealth)
    statsText := fmt.Sprintf("С:%d Л:%d", player.Strength, player.Dexterity)
    goldText := fmt.Sprintf("$: %d", player.Backpack.Gold)
    weaponText := weaponInfo
    modeText := fmt.Sprintf("%s", r.GetModeName())
    
    parts := []string{levelText, hpText, statsText, goldText, weaponText, modeText}
    if keysText != "" {
        parts = append(parts, keysText)
    }
    // рассчитываем позиции
    x := 2
    for i, part := range parts {
        var color tcell.Color
        
        switch i {
        case 0: 
            color = tcell.ColorYellow
        case 1: 
            if hpPercent > 0.6 {
                color = tcell.ColorGreen
            } else if hpPercent > 0.3 {
                color = tcell.ColorYellow
            } else {
                color = tcell.ColorRed
            }
        case 2: 
            color = tcell.ColorAqua
        case 3: 
            color = tcell.ColorGold
        case 4: 
            color = tcell.ColorSilver
        case 5: 
            if r.mode2D {
                color = tcell.ColorGreen
            } else {
                color = tcell.ColorBlue
            }
        }
        
        r.drawText(x, 1, part, color)
        x += runeLen(part) + 3
        
        if i < len(parts)-1 {
            r.drawText(x-3, 1, " | ", tcell.ColorDarkGray)
        }
    }

    r.drawBox(0, r.height-6, r.width, 6, "СООБЩЕНИЯ", tcell.ColorYellow, false)

    messages := r.game.GetMessages()
    maxVisible := 4
    start := 0
    if len(messages) > maxVisible {
        start = len(messages) - maxVisible
    }

    for i := 0; i < maxVisible && start+i < len(messages); i++ {
        msg := messages[start+i]
        maxLen := r.width - 6
        if utf8.RuneCountInString(msg) > maxLen {
            msg = string([]rune(msg)[:maxLen])
        }
        
        color := getMessageColor(msg)
        r.drawText(3, r.height-5+i, msg, color)
    }
    
    var hotkeys string
    if r.showMenu {
        hotkeys = "[↑↓] Выбор   [Enter] Выбрать   [ESC] Отмена"
    } else if r.showItemSelect {
    hotkeys = "[1-9] Выбрать   [ESC] Отмена   (0 — убрать оружие)"
    } else if r.showInventory {
        hotkeys = "[1-9] Использовать   [TAB] Вкладки   [ESC/I] Закрыть"
    } else {
        hotkeys = "[ESC] Меню   [I] Инвентарь   [V] Режим   [WASD] Движение   [1-9] Предметы"
    }
    
    r.drawTextCentered(r.height-1, hotkeys, tcell.ColorDarkGray)
}

func onOff(v bool) string {
	if v {
		return "+"
	}
	return "-"
}

func countKeys(p *entity.Player) (red, blue, yellow bool) {
	if p == nil {
		return
	}
	keys := p.Backpack.Items[entity.ItemKey]
	for _, it := range keys {
		k, ok := it.(*entity.Key)
		if !ok {
			continue
		}
		switch k.Color {
		case entity.RedKey:
			red = true
		case entity.BlueKey:
			blue = true
		case entity.YellowKey:
			yellow = true
		}
	}
	return
}

