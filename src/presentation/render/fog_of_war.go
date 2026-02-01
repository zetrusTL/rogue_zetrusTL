package render

import (
    "example.com/rogue/domain/entity"
)

type FogOfWar struct {
    explored [][]bool
    visible  [][]bool
    radius   int
}

func NewFogOfWar(width, height, radius int) *FogOfWar {
    explored := make([][]bool, height)
    visible := make([][]bool, height)
    
    for y := range explored {
        explored[y] = make([]bool, width)
        visible[y] = make([]bool, width)
    }
    
    return &FogOfWar{
        explored: explored,
        visible:  visible,
        radius:   radius,
    }
}

// обновляет видимые клетки относительно позиции игрока
func (f *FogOfWar) Update(playerPos entity.Point, level *entity.Level) {
    // сброс видимых клеток
    for y := range f.visible {
        for x := range f.visible[y] {
            f.visible[y][x] = false
        }
    }
    
    // Используем  алгоритм лучей в 8 направлениях
    
    if playerPos.Y >= 0 && playerPos.Y < len(f.visible) && 
       playerPos.X >= 0 && playerPos.X < len(f.visible[0]) {
        f.visible[playerPos.Y][playerPos.X] = true
        f.explored[playerPos.Y][playerPos.X] = true
    }
    
    for dy := -f.radius; dy <= f.radius; dy++ {
        for dx := -f.radius; dx <= f.radius; dx++ {
            if dx == 0 && dy == 0 {
                continue
            }
            
            x := playerPos.X + dx
            y := playerPos.Y + dy
            
            if x < 0 || x >= len(f.visible[0]) || y < 0 || y >= len(f.visible) {
                continue
            }
            
            distance := dx*dx + dy*dy
            if distance > f.radius*f.radius {
                continue
            }
            
            if f.hasLineOfSight(playerPos.X, playerPos.Y, x, y, level) {
                f.visible[y][x] = true
                f.explored[y][x] = true
            }
        }
    }
}

//  проверяет есть ли прямая видимость между двумя точками
func (f *FogOfWar) hasLineOfSight(x0, y0, x1, y1 int, level *entity.Level) bool {
    dx := abs(x1 - x0)
    dy := abs(y1 - y0)
    sx := 1
    if x0 > x1 {
        sx = -1
    }
    sy := 1
    if y0 > y1 {
        sy = -1
    }
    err := dx - dy
    
    x, y := x0, y0
    
    for {
        if level.Tiles[y][x] == entity.TileWall {
            return false
        }
        
        if x == x1 && y == y1 {
            break
        }
        
        e2 := 2 * err
        if e2 > -dy {
            err -= dy
            x += sx
        }
        if e2 < dx {
            err += dx
            y += sy
        }
    }
    
    return true
}

// для проверки видна ли клетка
func (f *FogOfWar) IsVisible(x, y int) bool {
    if y < 0 || y >= len(f.visible) || x < 0 || x >= len(f.visible[0]) {
        return false
    }
    return f.visible[y][x]
}

// для проверки исследована ли клетка
func (f *FogOfWar) IsExplored(x, y int) bool {
    if y < 0 || y >= len(f.explored) || x < 0 || x >= len(f.explored[0]) {
        return false
    }
    return f.explored[y][x]
}