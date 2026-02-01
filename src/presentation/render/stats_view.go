package render

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
)

type runStat struct {
	Name               string
	Gold               int
	DeepestLevel       int
	EnemiesDefeated    int
	Steps              int
	Victory            bool
	AmountConsumedFood int
	NumberElixirDrunk  int
	NumberScrollsRead  int
	TotalHits          int
}

func (r *Render) drawStats() {
	// фон
	for y := 3; y < r.height-6; y++ {
		for x := 0; x < r.width; x++ {
			r.screen.SetContent(x, y, ' ', nil,
				tcell.StyleDefault.Background(tcell.NewHexColor(0x222222)))
		}
	}
	
	width, height := 100, 20
	x, y := (r.width-width)/2, (r.height-height)/2
	if width > r.width {
		width = r.width - 2
		x = 1
	}
	if height > r.height-6 {
		height = r.height - 8
		y = 4
	}

	r.drawBox(x, y, width, height, "СТАТИСТИКА ПРОХОЖДЕНИЙ", tcell.Color(14), true)

	var rows []runStat

	if r.game != nil && !r.statsFromMain {
		s := r.game.GetStats()
		rows = append(rows, runStat{
			Name:               s.PlayerName,
			Gold:               s.GoldCollected,
			DeepestLevel:       s.DeepestLevel,
			EnemiesDefeated:    s.EnemiesDefeated,
			Steps:              s.StepsTaken,
			Victory:            s.IsVictory,
			AmountConsumedFood: s.AmountConsumedFood,
			NumberElixirDrunk:  s.NumberElixirDrunk,
			NumberScrollsRead:  s.NumberScrollsRead,
			TotalHits:          s.TotalHits,
		})
	}

	if r.gameRoot != nil {
		runs := r.gameRoot.GetLeaderboard(20)
		for _, rr := range runs {
			rows = append(rows, runStat{
				Name:               rr.PlayerName,
				Gold:               rr.GoldCollected,
				DeepestLevel:       rr.DeepestLevel,
				EnemiesDefeated:    rr.EnemiesDefeated,
				Steps:              rr.StepsTaken,
				Victory:            rr.IsVictory,
				AmountConsumedFood: rr.AmountConsumedFood,
				NumberElixirDrunk:  rr.NumberElixirDrunk,
				NumberScrollsRead:  rr.NumberScrollsRead,
				TotalHits:          rr.TotalHits,
			})
		}
	}


	sort.Slice(rows, func(i, j int) bool { return rows[i].Gold > rows[j].Gold })

	headY := y + 2
	r.drawText(x+2, headY, "Игрок", tcell.ColorSilver)
	r.drawText(x+22, headY, "Золото", tcell.ColorGold)
	r.drawText(x+32, headY, "Ур.", tcell.ColorYellow)
	r.drawText(x+37, headY, "Убийств", tcell.ColorRed)
	r.drawText(x+46, headY, "Еда", tcell.ColorGreen)
	r.drawText(x+51, headY, "Эликс", tcell.ColorBlue)
	r.drawText(x+58, headY, "Свитки", tcell.ColorAqua)
	r.drawText(x+66, headY, "Удары", tcell.ColorSilver)
	r.drawText(x+74, headY, "Шаги", tcell.ColorSilver)

	startY := headY + 2
	maxRows := height - 6
	if len(rows) < maxRows {
		maxRows = len(rows)
	}

	if len(rows) == 0 {
		r.drawTextCentered(y+height/2, "Пока нет записей", tcell.ColorDarkGray)
	} else {
		for i := 0; i < maxRows; i++ {
			rowY := startY + i
			rs := rows[i]

			name := rs.Name
			if name == "" {
				name = "Player"
			}
			if rs.Victory {
				name = "★ " + name
			}

			r.drawText(x+2, rowY, trimTo(name, 18), tcell.ColorWhite)
			r.drawText(x+22, rowY, fmt.Sprintf("%d", rs.Gold), tcell.ColorGold)
			r.drawText(x+32, rowY, fmt.Sprintf("%d", rs.DeepestLevel), tcell.ColorYellow)
			r.drawText(x+37, rowY, fmt.Sprintf("%d", rs.EnemiesDefeated), tcell.ColorRed)
			r.drawText(x+46, rowY, fmt.Sprintf("%d", rs.AmountConsumedFood), tcell.ColorGreen)
			r.drawText(x+51, rowY, fmt.Sprintf("%d", rs.NumberElixirDrunk), tcell.ColorBlue)
			r.drawText(x+58, rowY, fmt.Sprintf("%d", rs.NumberScrollsRead), tcell.ColorAqua)
			r.drawText(x+66, rowY, fmt.Sprintf("%d", rs.TotalHits), tcell.ColorSilver)
			r.drawText(x+74, rowY, fmt.Sprintf("%d", rs.Steps), tcell.ColorSilver)
		}
	}

	r.drawTextCentered(y+height-2, "[ESC] Назад", tcell.ColorDarkGray)
}

func trimTo(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
