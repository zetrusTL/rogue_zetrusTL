package render

import (
    "github.com/gdamore/tcell/v2"
    "example.com/rogue/domain/entity"
)

func (r *Render) drawGameWorld() {
    level := r.game.GetCurrentLevel()
    player := r.game.GetPlayer()
    // инициализируем туман
    if r.fogOfWar == nil {
        r.fogOfWar = NewFogOfWar(level.Width, level.Height, 8)
    }
    
    r.fogOfWar.Update(player.Pos, level)
    
    viewWidth := r.width
    viewHeight := r.height - 9
    
    playerX, playerY := int(player.Pos.X), int(player.Pos.Y)
    
    startX := playerX - viewWidth/2
    startY := playerY - viewHeight/2
    
    for screenY := 3; screenY < r.height-6; screenY++ {
        worldY := startY + (screenY - 3)
        
        for screenX := 0; screenX < r.width; screenX++ {
            worldX := startX + screenX
            
            if worldX < 0 || worldX >= level.Width ||
                worldY < 0 || worldY >= level.Height {
                r.screen.SetContent(screenX, screenY, ' ', nil,
                    tcell.StyleDefault.Background(tcell.NewHexColor(0x111111)))
                continue
            }
            
            isVisible := r.fogOfWar.IsVisible(worldX, worldY)
            isExplored := r.fogOfWar.IsExplored(worldX, worldY)
            
            var symbol rune = ' '
            var color tcell.Color = tcell.ColorBlack
            var bgColor tcell.Color = tcell.ColorBlack
            var style = tcell.StyleDefault
            
            if isVisible {
                tileType := level.GetTile(worldX, worldY)
                tileInfo := GetTileInfoForEntityTile(tileType)
                doorColor := tcell.ColorDefault
                if tileType == entity.TileDoorLocked || tileType == entity.TileDoorOpened {
                    p := entity.Point{X: worldX, Y: worldY}
                    if door, ok := level.Doors[p]; ok {
                        switch int(door.Color) {
                        case 0:
                            doorColor = tcell.ColorRed
                        case 1:
                            doorColor = tcell.ColorBlue
                        case 2:
                            doorColor = tcell.ColorYellow
                        }
                    }
                }

                symbol = tileInfo.Symbol
                color = tcell.NewHexColor(0xdddddd)
                bgColor = tcell.NewHexColor(0x222222)
                if doorColor != tcell.ColorDefault {
                    color = doorColor
                }

                switch tileType {
                case entity.TileDoorOpened:
                    symbol = '+'
                    color = tcell.NewHexColor(0x444422)
                    bgColor = tcell.NewHexColor(0x333311)
                    case entity.TileDoorLocked:
                    symbol = '+'
                    if doorColor != tcell.ColorDefault {
                        color = doorColor
                    } else {
                        color = tcell.NewHexColor(0x666666)
                    }
                    bgColor = tcell.NewHexColor(0x1a1a1a)

                }

                style = tcell.StyleDefault.Foreground(color).Background(bgColor)
                if tileInfo.Bold {
                    style = style.Bold(true)
                }

            } else if isExplored {

                tileType := level.GetTile(worldX, worldY)
                tileInfo := GetTileInfoForEntityTile(tileType)

                symbol = tileInfo.Symbol
                color = tcell.NewHexColor(0x666666)
                bgColor = tcell.NewHexColor(0x111111)

                switch tileType {
                case entity.TileDoorOpened:
                    symbol = '+'
                    color = tcell.NewHexColor(0x444422)
                    bgColor = tcell.NewHexColor(0x222211)
                case entity.TileDoorLocked:
                    symbol = '+'
                    color = tcell.NewHexColor(0x555555)
                    bgColor = tcell.NewHexColor(0x1a1a1a)
                }

                style = tcell.StyleDefault.Foreground(color).Background(bgColor)
                if tileInfo.Bold {
                    style = style.Bold(true)
                }

            } else {
                bgColor = tcell.ColorBlack
                style = tcell.StyleDefault.Background(bgColor)
            }

            r.screen.SetContent(screenX, screenY, symbol, nil, style)
        }
    }
}