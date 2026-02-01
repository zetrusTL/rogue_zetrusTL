package render


import (
	"github.com/gdamore/tcell/v2"
	"example.com/rogue/domain/entity"
	"time"
)

func (r *Render) handleInput() {
	ev := r.screen.PollEvent()
	if ev == nil {
		return
	}
	if key, ok := ev.(*tcell.EventKey); ok {
		r.handleKeyEvent(key)
	}
}

func (r *Render) handleKeyEvent(ev *tcell.EventKey) {
	if r.showStats {
		if ev.Key() == tcell.KeyEscape {
			r.showStats = false
		}
		return
	}
	// если открыто меню, обрабатываем только ввод меню
	if r.showMenu {
		r.handleMenuInput(ev)
		return
	}
	
	if r.showItemSelect {
		r.handleItemSelectInput(ev)
		return
	}
	
	if r.showInventory {
		r.handleInventoryInput(ev)
		return
	}
	
	switch ev.Key() {
	case tcell.KeyEscape:
		r.toggleMenu()
	
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'v', 'V':
			r.SwitchMode()
		case 'w', 'W':
			if !r.canMoveNow() {
					return
				}
			r.game.MovePlayer(entity.Top)
		case 's', 'S':
			if !r.canMoveNow() {
				return
			}
			r.game.MovePlayer(entity.Bottom)
		case 'a', 'A':
			if !r.canMoveNow() {
				return
			}
			r.game.MovePlayer(entity.Left)
		case 'd', 'D':
			if !r.canMoveNow() {
				return
			}
			r.game.MovePlayer(entity.Right)
		case ' ': // пропуск хода
			r.game.WaitTurn()
		case 'i', 'I', 'е', 'Е':
			r.toggleInventory()
		case '0':
			// Снять оружие
			r.game.EquipWeapon(0)

		case 'h', 'H', 'р', 'Р':
			r.openItemSelect(entity.ItemWeapon)
		case 'j', 'J', 'о', 'О':
			r.openItemSelect(entity.ItemFood)
		case 'k', 'K', 'л', 'Л':
			r.openItemSelect(entity.ItemElixir)
		case 'e', 'E', 'у', 'У':
			r.openItemSelect(entity.ItemScroll)
		}
	}
}

func (r *Render) canMoveNow() bool {
    now := time.Now()
    if !r.lastMoveAt.IsZero() && now.Sub(r.lastMoveAt) < r.moveCooldown {
        return false
    }
    r.lastMoveAt = now
    return true
}
