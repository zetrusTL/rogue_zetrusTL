package render

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"time"
	"unicode/utf8"
	"example.com/rogue/presentation"
	"example.com/rogue/domain/entity"
)
type Render struct {
	game          presentation.GameSession
	gameRoot presentation.Game
	screen        tcell.Screen
	width         int
	height        int
	mode2D        bool
	running       bool
	fogOfWar      *FogOfWar //туман
	//меню
	showMenu      bool
	menuItems     []string
	menuIndex     int
	//инвентарь
	showInventory bool
	inventoryTab  int
	showItemSelect bool
	selectItemType entity.ItemType

	messages []string
	lastPrintedMsgCount int
	showStats bool
	
	lastMoveAt time.Time
	moveCooldown time.Duration
	lastLevel *entity.Level

	showMainMenu  bool
	mainMenuItems []string
	mainMenuIndex int
	statsFromMain bool
}

func NewRender(game presentation.GameSession, gameRoot presentation.Game, width, height int) (*Render, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать экран: %v", err)
	}
	r := &Render{
		game:      game,
		gameRoot: gameRoot,
		screen:    screen,
		width:     width,
		height:    height,
		mode2D:    true,
		running:   true,
		showMenu:  false,
		menuIndex: 0,
		showMainMenu:  true,
		mainMenuIndex: 0,
		menuItems: func() []string {
			if gameRoot != nil {
				return []string{"Продолжить", "Новая игра", "Загрузить", "Инвентарь", "Статистика", "Сохранить", "Выйти"}
			}
			return []string{"Продолжить", "Инвентарь", "Статистика", "Сохранить", "Выйти"}
	}(),
	}
	r.moveCooldown = 90 * time.Millisecond

	r.refreshMainMenuItems()

	return r, nil
}


func (r *Render) Run() error {
	os.Setenv("LANG", "en_US.UTF-8")

	if err := r.screen.Init(); err != nil {
		return fmt.Errorf("Не удалось инициализировать экран: %v", err)
	}
	defer r.screen.Fini()

	r.screen.SetStyle(
		tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite),
	)
	mainLoop:

	
	for r.running {
		
		if r.showMainMenu {
		r.drawMainMenu()
		r.screen.Show()
		ev := r.screen.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			r.handleMainMenuInput(key)
		}
		time.Sleep(33 * time.Millisecond)
		continue
		}

		r.screen.Clear()
		r.draw()
		r.screen.Show()
		r.handleInput()

		if r.game.GetStatus() != entity.StatusRunning {
		r.drawGameOver()

			for {
				ev := r.screen.PollEvent()
				if key, ok := ev.(*tcell.EventKey); ok {

					if key.Key() == tcell.KeyEscape {
						r.running = false
						return nil
					}

					if key.Rune() == 'r' || key.Rune() == 'R' {
						session, err := r.gameRoot.StartNewGame("Player")
						if err == nil {
							r.game = session
							r.fogOfWar = nil
							r.lastLevel = nil 
							r.messages = nil  
							continue mainLoop 
						}
					}
				}
			}
		}
		time.Sleep(66 * time.Millisecond)
	}

	return nil
}

func (r *Render) Stop() {
	r.running = false
}

func (r *Render) draw() {
	if r.showStats {
		r.drawStats()
		r.drawUI()
		return
	}

	if r.showInventory {
		r.drawInventory()
	} else {
		if r.mode2D {
			r.draw2D()
		} else {
			r.draw3D()
		}
	}
	r.drawUI()
	if !r.showInventory && !r.showItemSelect {  //рисуем только если не открытьт выбор предмета
		r.drawMenu()
	}
	r.drawItemSelect()
}

func (r *Render) draw2D() {
	level := r.game.GetCurrentLevel()
	if r.lastLevel != level {
		r.fogOfWar = NewFogOfWar(level.Width, level.Height, 6)
		r.lastLevel = level
	}

	gameHeight := r.height - 9
	r.clearArea(0, 3, r.width, gameHeight)

	r.drawGameWorld()
	r.drawItems()
	r.drawMonsters()
	r.drawExit()
	// игрока рисуем поверх всего
	centerX, centerY := r.width/2, (r.height-9)/2+3
	r.drawPlayerAt(centerX, centerY)
}

func (r *Render) draw3D() {
	gameHeight := r.height - 9
	r.clearArea(0, 3, r.width, gameHeight)

	startY := 5

	r.drawTextCentered(startY, "╔══════════════════════════════╗", tcell.ColorBlue)
	r.drawTextCentered(startY+1, "║       3D MODE ACTIVE         ║", tcell.ColorBlue)
	r.drawTextCentered(startY+2, "╚══════════════════════════════╝", tcell.ColorBlue)

	r.drawTextCentered(startY+4, "УПРАВЛЕНИЕ:", tcell.ColorYellow)
	r.drawTextCentered(startY+5, "WASD - Вперёд / Назад", tcell.ColorWhite)
	r.drawTextCentered(startY+6, "AD   - Поворот", tcell.ColorWhite)
	r.drawTextCentered(startY+7, "V    - Переключить на 2D режим", tcell.ColorWhite)
	r.drawTextCentered(startY+8, "ESC  - Выход из игры", tcell.ColorWhite)

	lineY := startY + 10
	if lineY < r.height-7 {
		r.drawHorizontalLine(lineY, tcell.ColorGray)
	}
}

func (r *Render) GetModeName() string {
	if r.mode2D {
		return "2D"
	}
	return "3D"
}

func (r *Render) SwitchMode() {
	r.mode2D = !r.mode2D
}

func (r *Render) toggleMenu() {
	r.showMenu = !r.showMenu
	r.menuIndex = 0
}

func (r *Render) toggleInventory() {
	r.showInventory = !r.showInventory
	if r.showInventory {
		r.inventoryTab = 0
	}
}

// Базовые методы для работы с экраном
func (r *Render) drawText(x, y int, text string, color tcell.Color) {
	style := tcell.StyleDefault.Foreground(color)
	i := 0
	for _, ch := range text {
		r.screen.SetContent(x+i, y, ch, nil, style)
		i++
	}
}

func (r *Render) drawTextBold(x, y int, text string, color tcell.Color) {
	style := tcell.StyleDefault.Foreground(color).Bold(true)
	i := 0
	for _, ch := range text {
		r.screen.SetContent(x+i, y, ch, nil, style)
		i++
	}
}

func (r *Render) drawTextCentered(y int, text string, color tcell.Color) {
	x := (r.width - runeLen(text)) / 2
	if x < 0 {
		x = 0
	}
	r.drawText(x, y, text, color)
}

func (r *Render) clearArea(x1, y1, width, height int) {
	for y := y1; y < y1+height && y < r.height; y++ {
		for x := x1; x < x1+width && x < r.width; x++ {
			r.screen.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
	}
}

func (r *Render) drawBox(x1, y1, width, height int, title string, color tcell.Color, bold bool) {
	r.clearArea(x1, y1, width, height)

	r.drawText(x1, y1, "┌", color)
	for x := x1 + 1; x < x1+width-1; x++ {
		r.drawText(x, y1, "─", color)
	}
	r.drawText(x1+width-1, y1, "┐", color)

	if title != "" {
		titleText := " " + title + " "
		tx := x1 + (width-runeLen(titleText))/2
		if tx < x1+1 {
			tx = x1 + 1
		}
		if bold {
			r.drawTextBold(tx, y1, titleText, color)
		} else {
			r.drawText(tx, y1, titleText, color)
		}
	}

	for y := y1 + 1; y < y1+height-1; y++ {
		r.drawText(x1, y, "│", color)
		r.drawText(x1+width-1, y, "│", color)
	}

	r.drawText(x1, y1+height-1, "└", color)
	for x := x1 + 1; x < x1+width-1; x++ {
		r.drawText(x, y1+height-1, "─", color)
	}
	r.drawText(x1+width-1, y1+height-1, "┘", color)
}

func (r *Render) drawHorizontalLine(y int, color tcell.Color) {
	for x := 0; x < r.width; x++ {
		r.drawText(x, y, "─", color)
	}
}

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}

func (r *Render) drawExit() {
	level := r.game.GetCurrentLevel()
	player := r.game.GetPlayer()
	if level == nil || player == nil {
		return
	}

	// координаты выхода в мире
	ex := level.ExitPos.X
	ey := level.ExitPos.Y

	// не рисуем на игроке
	if ex == player.Pos.X && ey == player.Pos.Y {
		return
	}

	//рисуем только если клетка видима
	if r.fogOfWar != nil && !r.fogOfWar.IsVisible(ex, ey) {
		return
	}

	// переводим в экранные координаты относительно игрока
	dx := ex - player.Pos.X
	dy := ey - player.Pos.Y

	screenX := r.width/2 + dx
	screenY := (r.height-9)/2 + 3 + dy

	if screenX < 0 || screenX >= r.width || screenY < 3 || screenY >= r.height-6 {
		return
	}

	info := GetExitInfo()
	style := tcell.StyleDefault.Foreground(info.Color)
	if info.Bold {
		style = style.Bold(true)
	}

	r.screen.SetContent(screenX, screenY, info.Symbol, nil, style)
}


func (r *Render) drawPlayerAt(x, y int) {
    style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
    r.screen.SetContent(x, y, '@', nil, style)
}

func (r *Render) drawGameOver() {
	r.screen.Clear()

	w, h := r.screen.Size()
	cx := w / 2

	title := "GAME OVER"
	titleColor := tcell.ColorRed
	if r.game.GetStatus() == entity.StatusVictory {
		title = "YOU WIN!"
		titleColor = tcell.ColorGreen
	}
	r.drawText(cx-len(title)/2, h/2-4, title, titleColor)


	stats := r.game.GetStats()

	lines := []string{
	fmt.Sprintf("Level reached: %d", stats.DeepestLevel),
	fmt.Sprintf("Gold collected: %d", stats.GoldCollected),
	}


	y := h/2 - 1
	for _, l := range lines {
		r.drawText(cx-len(l)/2, y, l, tcell.ColorWhite)
		y++
	}

	hint := "R - Restart   ESC - Exit"
	r.drawText(cx-len(hint)/2, y+2, hint, tcell.ColorDarkGray)

	r.screen.Show()
}

func (r *Render) refreshMainMenuItems() {
	items := []string{"Начать новую игру"}
	if r.gameRoot != nil && r.gameRoot.HasSavedGame() {
		items = append(items, "Загрузить сохранение")
	}
	items = append(items, "Статистика", "Выйти из игры")
	r.mainMenuItems = items
	if r.mainMenuIndex >= len(r.mainMenuItems) {
		r.mainMenuIndex = 0
	}
}

func (r *Render) drawMainMenu() {
	r.screen.Clear()

	r.refreshMainMenuItems()

	width := 46
	height := len(r.mainMenuItems) + 6
	x := (r.width - width) / 2
	y := (r.height - height) / 2

	r.drawBox(x, y, width, height, "ROGUE", tcell.Color(14), true)
	r.drawTextCentered(y+1, "Главное меню", tcell.ColorWhite)

	for i, item := range r.mainMenuItems {
		yy := y + 3 + i
		if i == r.mainMenuIndex {
			r.drawText(x+3, yy, ">", tcell.ColorYellow)
		}
		r.drawText(x+5, yy, item, tcell.ColorWhite)
	}

	help := "[↑↓] Выбор   [Enter] Выбрать   [ESC] Выход"
	r.drawTextCentered(y+height-2, help, tcell.ColorDarkGray)
}

func (r *Render) handleMainMenuInput(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		r.Stop()
	case tcell.KeyUp:
		r.mainMenuIndex = (r.mainMenuIndex - 1 + len(r.mainMenuItems)) % len(r.mainMenuItems)
	case tcell.KeyDown:
		r.mainMenuIndex = (r.mainMenuIndex + 1) % len(r.mainMenuItems)
	case tcell.KeyEnter:
		r.executeMainMenuAction()
	}
}

func (r *Render) executeMainMenuAction() {
	item := r.mainMenuItems[r.mainMenuIndex]

	switch item {
	case "Начать новую игру":
		if r.gameRoot == nil {
			return
		}
		session, err := r.gameRoot.StartNewGame("Player")
		if err != nil {
			return
		}
		r.game = session
		r.showMainMenu = false
		r.showStats = false
		r.statsFromMain = false
		r.showInventory = false
		r.showMenu = false
		r.fogOfWar = nil
		r.lastLevel = nil
		r.messages = nil

	case "Загрузить сохранение":
		if r.gameRoot == nil || !r.gameRoot.HasSavedGame() {
			return
		}
		session, err := r.gameRoot.LoadLastGame()
		if err != nil {
			return
		}
		r.game = session
		r.showMainMenu = false
		r.showStats = false
		r.statsFromMain = false
		r.showInventory = false
		r.showMenu = false
		r.fogOfWar = nil
		r.lastLevel = nil
		r.messages = nil

	case "Статистика":
		r.showStats = true
		r.statsFromMain = true

	case "Выйти из игры":
		r.Stop()
	}
}
