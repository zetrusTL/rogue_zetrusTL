package entity

type Door struct {
	Object
	Color  KeyColor
	Opened bool
}

type Room struct {
	Object
	GridPos Point

	// Массив указателей на двери (Top, Right, Bottom, Left).
	// nil если двери нет.
	Doors [4]*Door

	// Связи с соседними комнатами (для графа).
	// Indices: 0:Top, 1:Right, 2:Bottom, 3:Left
	Connections [4]*Room

	IsLocked bool
}

type Corridor struct {
	// Точки изгиба для отрисовки
	Points []Point
}

type Level struct {
	Width, Height int
	LevelNum      int
	Rooms         []*Room
	Corridors     []*Corridor
	ExitPos       Point
	StartPos      Point

	Monsters map[Point]*Monster
	Items    map[Point][]Item
	Doors    map[Point]*Door

	// Карта занятости для быстрого поиска (Collision Map)
	Tiles [][]TileType
}

func (level *Level) GetTile(x, y int) TileType {
	return level.Tiles[y][x]
}

func (level *Level) SetTile(x, y int, tile TileType) {
	level.Tiles[y][x] = tile
}

func NewLevel(numLevel int) *Level {
	tiles := make([][]TileType, MapHeight)

	for y := 0; y < MapHeight; y++ {
		tiles[y] = make([]TileType, MapWidth)
	}

	return &Level{
		LevelNum:  numLevel,
		Width:     MapWidth,
		Height:    MapHeight,
		Tiles:     tiles,
		Rooms:     make([]*Room, 0),
		Corridors: make([]*Corridor, 0),
		Monsters:  make(map[Point]*Monster),
		Items:     make(map[Point][]Item),
		Doors:     make(map[Point]*Door),
	}
}
