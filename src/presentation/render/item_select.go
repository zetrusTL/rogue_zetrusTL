package render

import (
	"fmt"

	"example.com/rogue/domain/entity"
	"github.com/gdamore/tcell/v2"
)

func (r *Render) openItemSelect(itemType entity.ItemType) {
	// если открыт инвентарь/меню — закрываем, чтобы не конфликтовать
	r.showMenu = false
	r.showInventory = false

	r.showItemSelect = true
	r.selectItemType = itemType
}

func (r *Render) closeItemSelect() {
	r.showItemSelect = false
}

func (r *Render) drawItemSelect() {
	if !r.showItemSelect {
		return
	}

	// фон (как в меню/инвентаре)
	for y := 3; y < r.height-6; y++ {
		for x := 0; x < r.width; x++ {
			r.screen.SetContent(x, y, ' ', nil,
				tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)))
		}
	}

	items := r.game.GetInventoryItems(r.selectItemType)

	title := r.itemTypeTitle(r.selectItemType)
	width, height := 56, 12
	x, y := (r.width-width)/2, (r.height-height)/2

	r.drawBox(x, y, width, height, title, tcell.Color(14), true)

	if len(items) == 0 {
		r.drawTextCentered(y+height/2, "Нет предметов этого типа", tcell.ColorDarkGray)
		r.drawTextCentered(y+height-2, "[ESC] Закрыть", tcell.ColorDarkGray)
		return
	}

	// подсказка для оружия
	if r.selectItemType == entity.ItemWeapon {
		r.drawText(x+2, y+1, "0 — убрать оружие", tcell.ColorDarkGray)
	}

	maxShow := 9
	if len(items) < maxShow {
		maxShow = len(items)
	}

	for i := 0; i < maxShow; i++ {
		it := items[i]
		lineY := y + 3 + i

		info := GetItemInfo(it)
		r.drawText(x+2, lineY, fmt.Sprintf("%d.", i+1), tcell.ColorYellow)
		r.drawText(x+6, lineY, string(info.Symbol), info.Color)
		r.drawText(x+8, lineY, it.GetName(), info.Color)
	}

	r.drawTextCentered(y+height-2, "[1-9] Выбрать   [ESC] Отмена", tcell.ColorDarkGray)
}

func (r *Render) itemTypeTitle(t entity.ItemType) string {
	switch t {
	case entity.ItemWeapon:
		return "ОРУЖИЕ"
	case entity.ItemFood:
		return "ЕДА"
	case entity.ItemElixir:
		return "ЭЛИКСИРЫ"
	case entity.ItemScroll:
		return "СВИТКИ"
	default:
		return "ПРЕДМЕТЫ"
	}
}

func (r *Render) handleItemSelectInput(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		r.closeItemSelect()
		return

	case tcell.KeyRune:
		switch ev.Rune() {
		case '0':
			if r.selectItemType == entity.ItemWeapon {
				r.game.EquipWeapon(0)
				r.closeItemSelect()
			}
			return

		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			idx := int(ev.Rune() - '1') // 0..8
			items := r.game.GetInventoryItems(r.selectItemType)
			if idx < 0 || idx >= len(items) {
				return
			}

			// оружие 1..9
			if r.selectItemType == entity.ItemWeapon {
				r.game.EquipWeapon(idx + 1)
				r.closeItemSelect()
				return
			}

			_ = r.game.UseItem(r.selectItemType, idx)
			r.closeItemSelect()
			return
		}
	}
}
