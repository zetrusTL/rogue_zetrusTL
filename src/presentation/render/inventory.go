package render

import (
	"fmt"
	"sort"

	"example.com/rogue/domain/entity"
	"github.com/gdamore/tcell/v2"
)

const (
	InvTabAll int = iota
	InvTabWeapons
	InvTabFood
	InvTabElixirs
	InvTabScrolls
	InvTabCount
)
type sortedItem struct {
	item     entity.Item
	itemType entity.ItemType
	index    int // индекс ВНУТРИ слайса items данного типа
	priority int
}
// сортировка предметов
func sortItems(items []entity.Item, itemType entity.ItemType) []sortedItem {
	out := make([]sortedItem, 0, len(items))

	for i, it := range items {
		priority := 0

		switch it.GetType() {
		case entity.ItemWeapon:
			priority = 1000
			if w, ok := it.(*entity.Weapon); ok {
				priority += w.StrengthBonus * 10
			}
		case entity.ItemFood:
			priority = 900
			if f, ok := it.(*entity.Food); ok {
				priority += f.HealAmount
			}
		case entity.ItemElixir:
			priority = 800
		case entity.ItemScroll:
			priority = 700
		case entity.ItemTreasure:
			priority = 600
			if t, ok := it.(*entity.Treasure); ok {
				priority += t.Value / 10
			}
		default:
			priority = 500
		}

		// небольшой вес по имени
		priority += len(it.GetName())

		out = append(out, sortedItem{
			item:     it,
			itemType: itemType,
			index:    i,
			priority: priority,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].priority != out[j].priority {
			return out[i].priority > out[j].priority
		}
		return out[i].item.GetName() < out[j].item.GetName()
	})

	return out
}

// отрисовка инвентаря
func (r *Render) drawInventory() {
	// заливаем фон
	for y := 3; y < r.height-6; y++ {
		for x := 0; x < r.width; x++ {
			r.screen.SetContent(x, y, ' ', nil,
				tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)))
		}
	}

	player := r.game.GetPlayer()
	if player == nil {
		r.drawTextCentered(r.height/2, "Нет игрока", tcell.ColorRed)
		return
	}
	// рисуем окно инвентаря
	width, height := 74, 22
	x, y := (r.width-width)/2, (r.height-height)/2
	r.drawBox(x, y, width, height, "ИНВЕНТАРЬ", tcell.Color(12), true)

	totalItems := r.getTotalInventoryCount()
	capText := fmt.Sprintf("Предметы: %d/%d", totalItems, 36) 
	goldText := fmt.Sprintf("Золото: %d", player.Backpack.Gold)

	r.drawText(x+2, y+1, capText, tcell.ColorSilver)
	r.drawText(x+width-runeLen(goldText)-2, y+1, goldText, tcell.ColorGold)

	tabs := []string{"Все", "Оружие", "Еда", "Эликсиры", "Свитки"}
	tabNames := []string{"Все предметы", "Оружие", "Еда", "Эликсиры", "Свитки"}
	tabX := x + 2

	for i := 0; i < InvTabCount; i++ {
		color := tcell.ColorWhite
		if i == r.inventoryTab {
			color = tcell.ColorYellow
			r.drawText(tabX-1, y+2, ">", color)
		}
		r.drawText(tabX, y+2, tabs[i], color)
		tabX += runeLen(tabs[i]) + 4
	}

	items := r.getFilteredAndSortedItems()

	startY := y + 4
	maxItems := height - 8

	if len(items) == 0 {
		emptyText := "Инвентарь пуст"
		if r.inventoryTab != InvTabAll {
			emptyText = fmt.Sprintf("Нет предметов: %s", tabNames[r.inventoryTab])
		}
		r.drawTextCentered(y+height/2, emptyText, tcell.ColorDarkGray)
	} else {
		r.drawText(x+2, startY-1, "№", tcell.ColorDarkGray)
		r.drawText(x+6, startY-1, "Предмет", tcell.ColorDarkGray)
		r.drawText(x+40, startY-1, "Характеристики", tcell.ColorDarkGray)
		r.drawText(x+62, startY-1, "Действия", tcell.ColorDarkGray)

		limit := maxItems
		if len(items) < limit {
			limit = len(items)
		}

		for i := 0; i < limit; i++ {
			si := items[i]
			it := si.item
			itemY := startY + i

			// номера 1-9
			if i < 9 {
				r.drawText(x+2, itemY, fmt.Sprintf("%d.", i+1), tcell.ColorYellow)
			} else {
				r.drawText(x+2, itemY, "·", tcell.ColorDarkGray)
			}

			info := GetItemInfo(it)
			r.drawText(x+4, itemY, string(info.Symbol), info.Color)

			name := it.GetName()
			if w, ok := it.(*entity.Weapon); ok && player.Weapon == w {
				name = "[" + name + "]"
			}
			r.drawText(x+6, itemY, name, info.Color)

			statsX := x + 40
			switch v := it.(type) {
			case *entity.Food:
				r.drawText(statsX, itemY, fmt.Sprintf("Восстановление: +%d HP", v.HealAmount), tcell.ColorGreen)
			case *entity.Weapon:
				r.drawText(statsX, itemY, fmt.Sprintf("Сила: +%d", v.StrengthBonus), tcell.ColorRed)
			case *entity.Treasure:
				r.drawText(statsX, itemY, fmt.Sprintf("Ценность: %d золота", v.Value), tcell.ColorGold)
			case *entity.Elixir:
				r.drawText(statsX, itemY, "Эликсир", tcell.ColorBlue)
			case *entity.Scroll:
				r.drawText(statsX, itemY, "Свиток", tcell.ColorAqua)
			}

			if i < 9 {
				// для оружия: Use -> Equip
				if si.itemType == entity.ItemWeapon {
					r.drawText(x+62, itemY, fmt.Sprintf("[%d]Экип", i+1), tcell.ColorDarkGray)
				} else {
					r.drawText(x+62, itemY, fmt.Sprintf("[%d]Исп", i+1), tcell.ColorDarkGray)
				}
			}
		}

		if len(items) > maxItems {
			moreText := fmt.Sprintf("... и еще %d предметов", len(items)-maxItems)
			r.drawTextCentered(y+height-3, moreText, tcell.ColorDarkGray)
		}
	}

	hints := []string{
		"[1-9] Использовать/экипировать",
		"[TAB/←→] Вкладки",
		"[ESC/I] Закрыть",
	}
	hintX := x + 2
	for _, htxt := range hints {
		r.drawText(hintX, y+height-2, htxt, tcell.ColorDarkGray)
		hintX += runeLen(htxt) + 3
	}

	tabInfo := fmt.Sprintf("Вкладка: %s (%d предметов)", tabNames[r.inventoryTab], len(items))
	r.drawTextCentered(y+height-1, tabInfo, tcell.ColorSilver)
}

func (r *Render) handleInventoryInput(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		r.toggleInventory()
	case tcell.KeyTab:
		r.inventoryTab = (r.inventoryTab + 1) % InvTabCount
	case tcell.KeyLeft:
		r.inventoryTab = (r.inventoryTab - 1 + InvTabCount) % InvTabCount
	case tcell.KeyRight:
		r.inventoryTab = (r.inventoryTab + 1) % InvTabCount
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'i', 'I', 'е', 'Е':
			r.toggleInventory()
		case 'h', 'H':
			r.inventoryTab = int(entity.ItemWeapon)
			return

		case 'j', 'J':
			r.inventoryTab = int(entity.ItemFood)
			return

		case 'k', 'K':
			r.inventoryTab = int(entity.ItemElixir)
			return

		case 'e', 'E':
			r.inventoryTab = int(entity.ItemScroll)
			return
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			idx := int(ev.Rune() - '1') 
			items := r.getFilteredAndSortedItems()
			if idx < 0 || idx >= len(items) {
				return
			}

			si := items[idx]

			realIndex := 0
			sameType := r.game.GetInventoryItems(si.itemType)
			for i, it := range sameType {
				if it == si.item {
					realIndex = i
					break
				}
			}

			if si.itemType == entity.ItemWeapon {
				r.game.EquipWeapon(realIndex + 1)
				r.toggleInventory()
				return
			}

			if err := r.game.UseItem(si.itemType, realIndex+1); err != nil { 
				r.messages = append(r.messages, "Не удалось использовать: "+err.Error())
			} else {
				r.messages = append(r.messages, "Предмет использован")
			}


			r.toggleInventory()

		}
	}
}

func (r *Render) getTotalInventoryCount() int {
	total := 0
	for _, t := range []entity.ItemType{entity.ItemWeapon, entity.ItemFood, entity.ItemElixir, entity.ItemScroll} {
		total += len(r.game.GetInventoryItems(t))
	}
	return total
}

// получить отфильтрованные и отсортированные предметы
func (r *Render) getFilteredAndSortedItems() []sortedItem {
	switch r.inventoryTab {
	case InvTabWeapons:
		return sortItems(r.game.GetInventoryItems(entity.ItemWeapon), entity.ItemWeapon)
	case InvTabFood:
		return sortItems(r.game.GetInventoryItems(entity.ItemFood), entity.ItemFood)
	case InvTabElixirs:
		return sortItems(r.game.GetInventoryItems(entity.ItemElixir), entity.ItemElixir)
	case InvTabScrolls:
		return sortItems(r.game.GetInventoryItems(entity.ItemScroll), entity.ItemScroll)
	case InvTabAll:
		// собираем все типы + сортируем единым списком
		all := make([]sortedItem, 0, 32)

		for _, t := range []entity.ItemType{entity.ItemWeapon, entity.ItemFood, entity.ItemElixir, entity.ItemScroll} {
			part := sortItems(r.game.GetInventoryItems(t), t)
			all = append(all, part...)
		}

		sort.Slice(all, func(i, j int) bool {
			if all[i].priority != all[j].priority {
				return all[i].priority > all[j].priority
			}
			return all[i].item.GetName() < all[j].item.GetName()
		})

		return all
	default:
		return []sortedItem{}
	}
}

