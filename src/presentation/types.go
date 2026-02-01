package presentation

import "github.com/gdamore/tcell/v2"

type TileInfo struct {
    Symbol rune
    Color  tcell.Color
    Bold   bool
}

func GetTileInfo(tile TileType) TileInfo {
    switch tile {
    case TileWall:
        return TileInfo{Symbol: '█', Color: tcell.NewHexColor(0x555555), Bold: false}
    case TileFloor:
        return TileInfo{Symbol: '.', Color: tcell.NewHexColor(0x666666), Bold: false}
    case TileDoor:
        return TileInfo{Symbol: '+', Color: tcell.ColorYellow, Bold: true}
    case TileCorridor:
        return TileInfo{Symbol: '#', Color: tcell.NewHexColor(0x444444), Bold: false}
    default:
        return TileInfo{Symbol: '?', Color: tcell.ColorRed, Bold: false}
    }
}

type TileType int

const (
    TileWall TileType = iota
    TileFloor
    TileDoor
    TileCorridor
)