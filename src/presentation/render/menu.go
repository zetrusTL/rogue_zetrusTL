package render

import (
	"github.com/gdamore/tcell/v2"

)

func (r *Render) drawMenu() {
	if !r.showMenu {
		return
	}
	
	// полупрозрачный фон
	for y := 3; y < r.height-6; y++ {
		for x := 0; x < r.width; x++ {
			r.screen.SetContent(x, y, ' ', nil,
				tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)))
		}
	}
	
	width, height := 30, len(r.menuItems)+6
	x, y := (r.width-width)/2, (r.height-height)/2
	
	r.drawBox(x, y, width, height, "МЕНЮ", tcell.Color(14), true)
	
	r.drawTextCentered(y+1, "Игра приостановлена", tcell.ColorWhite)
	
	// пунткы
	for i, item := range r.menuItems {
		textY := y + 3 + i
		if i == r.menuIndex {
			r.drawText(x+3, textY, ">", tcell.ColorYellow)
		}
		color := tcell.ColorWhite
		r.drawText(x+5, textY, item, color)
	}
	
	help := "[↑↓] Выбор   [Enter] Выбрать   [ESC] Отмена"
	r.drawTextCentered(y+height-2, help, tcell.ColorDarkGray)
}

func (r *Render) handleMenuInput(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		r.toggleMenu()
	case tcell.KeyUp:
		r.menuIndex = (r.menuIndex - 1 + len(r.menuItems)) % len(r.menuItems)
	case tcell.KeyDown:
		r.menuIndex = (r.menuIndex + 1) % len(r.menuItems)
	case tcell.KeyEnter:
		r.executeMenuAction()
	}
}

func (r *Render) executeMenuAction() {
	// Расширенное меню, если есть gameRoot
	if r.gameRoot != nil {
		switch r.menuIndex {
		case 0: 
			r.toggleMenu()

		case 1: 
			session, err := r.gameRoot.StartNewGame("Player")
			if err != nil {
				return
			}
			r.game = session
			r.showInventory = false
			r.showItemSelect = false
			r.showStats = false
			r.toggleMenu()

		case 2: 
			if !r.gameRoot.HasSavedGame() {
				return
			}
			session, err := r.gameRoot.LoadLastGame()
			if err != nil {
				return
			}
			r.game = session
			r.showInventory = false
			r.showItemSelect = false
			r.showStats = false
			r.toggleMenu()

		case 3: 
			r.toggleMenu()
			r.toggleInventory()

		case 4: 
			r.toggleMenu()
			r.showStats = true

		case 5: 

			r.toggleMenu()

		case 6: 
			r.Stop()
		}
		return
	}

	// Обычное меню (без gameRoot)
	switch r.menuIndex {
	case 0: 
		r.toggleMenu()
	case 1: 
		r.toggleMenu()
		r.toggleInventory()
	case 2: 
		r.toggleMenu()
		r.showStats = true
	case 3: 
		r.toggleMenu()
	case 4: 
		r.Stop()
	}
}

