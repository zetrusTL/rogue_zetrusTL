package generation

import (
	"math/rand"

	"example.com/rogue/domain/entity"
)

type generator struct {
	level     *entity.Level
	grid      [entity.RoomsInRow][entity.RoomsInCol]*entity.Room
	startRoom *entity.Room
}

func (lg *LevelGenerator) generateMap(numLevel int) (*entity.Level, *entity.Room) {
	gen := &generator{
		level: entity.NewLevel(numLevel),
	}

	for y := 0; y < entity.MapHeight; y++ {
		for x := 0; x < entity.MapWidth; x++ {
			gen.level.SetTile(x, y, entity.TileWall)
		}
	}

	gen.generateSectors()
	gen.generateConnections()
	gen.generateRoomsGeometry()
	gen.generateCorridors()
	gen.renderToMap()
	gen.placeStartAndExit()
	gen.generateKeysAndLocks()

	return gen.level, gen.startRoom
}

func (g *generator) generateSectors() {
	for i := 0; i < entity.RoomsInRow; i++ {
		for j := 0; j < entity.RoomsInCol; j++ {
			room := &entity.Room{
				GridPos: entity.Point{X: j, Y: i},
			}
			g.grid[i][j] = room
			g.level.Rooms = append(g.level.Rooms, room)
		}
	}
}

func (g *generator) generateConnections() {
	type neighbor struct {
		room *entity.Room
		dir  int
	}

	startRoom := g.grid[rand.Intn(3)][rand.Intn(3)]
	visited := make(map[*entity.Room]bool)
	stack := []*entity.Room{startRoom}
	visited[startRoom] = true

	for len(stack) > 0 {
		current := stack[len(stack)-1]

		var unvisitedNeighbors []neighbor

		if current.GridPos.Y > 0 {
			n := g.grid[current.GridPos.Y-1][current.GridPos.X]
			if !visited[n] {
				unvisitedNeighbors = append(unvisitedNeighbors, neighbor{n, 0})
			}
		}
		if current.GridPos.X < entity.RoomsInCol-1 {
			n := g.grid[current.GridPos.Y][current.GridPos.X+1]
			if !visited[n] {
				unvisitedNeighbors = append(unvisitedNeighbors, neighbor{n, 1})
			}
		}
		if current.GridPos.Y < entity.RoomsInRow-1 {
			n := g.grid[current.GridPos.Y+1][current.GridPos.X]
			if !visited[n] {
				unvisitedNeighbors = append(unvisitedNeighbors, neighbor{n, 2})
			}
		}
		if current.GridPos.X > 0 {
			n := g.grid[current.GridPos.Y][current.GridPos.X-1]
			if !visited[n] {
				unvisitedNeighbors = append(unvisitedNeighbors, neighbor{n, 3})
			}
		}

		if len(unvisitedNeighbors) > 0 {
			next := unvisitedNeighbors[rand.Intn(len(unvisitedNeighbors))]

			current.Connections[next.dir] = next.room
			next.room.Connections[opposite(next.dir)] = current

			visited[next.room] = true
			stack = append(stack, next.room)
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	for i := 0; i < entity.RoomsInRow; i++ {
		for j := 0; j < entity.RoomsInCol; j++ {
			room := g.grid[i][j]

			// Пробуем соединить с правым соседом
			if j < entity.RoomsInCol-1 {
				right := g.grid[i][j+1]
				if room.Connections[1] == nil && rand.Float64() < 0.3 { // 30% шанс добавить лишнюю связь
					room.Connections[1] = right
					right.Connections[3] = room
				}
			}
			// Пробуем соединить с нижним соседом
			if i < entity.RoomsInRow-1 {
				bottom := g.grid[i+1][j]
				if room.Connections[2] == nil && rand.Float64() < 0.3 {
					room.Connections[2] = bottom
					bottom.Connections[0] = room
				}
			}
		}
	}
}

func (g *generator) generateRoomsGeometry() {
	for _, room := range g.level.Rooms {
		offsetY := room.GridPos.Y * entity.SectorHeight
		offsetX := room.GridPos.X * entity.SectorWidth

		maxY := (room.GridPos.Y + 1) * entity.SectorHeight
		maxX := (room.GridPos.X + 1) * entity.SectorWidth

		room.Pos.X = rand.Intn(entity.CornerHorRange)/3 + offsetX + 2  // TopLeft X
		room.Pos.Y = rand.Intn(entity.CornerVertRange)/3 + offsetY + 2 // TopLeft Y

		room.Size.Width = min(rand.Intn(entity.CornerHorRange-2)+4, maxX-room.Pos.X-2)
		room.Size.Height = min(rand.Intn(entity.CornerVertRange-2)+4, maxY-room.Pos.Y-2)

		g.generateDoors(room)
	}
}

func (g *generator) generateDoors(room *entity.Room) {
	if room.Connections[0] != nil {
		d := &entity.Door{}
		d.Pos.Y = room.Pos.Y
		d.Pos.X = randRange(room.Pos.X+1, room.Pos.X+room.Size.Width-1)
		d.Color = entity.KeyColor(randRange(0, entity.ColorCount))
		if shouldPlaceDoor() {
			room.Doors[entity.Top] = d
		}
	}
	if room.Connections[1] != nil {
		d := &entity.Door{}
		d.Pos.X = room.Pos.X + room.Size.Width
		d.Pos.Y = randRange(room.Pos.Y+1, room.Pos.Y+room.Size.Height-1)
		d.Color = entity.KeyColor(randRange(0, entity.ColorCount))
		if shouldPlaceDoor() {
			room.Doors[entity.Right] = d
		}
	}
	if room.Connections[2] != nil {
		d := &entity.Door{}
		d.Pos.Y = room.Pos.Y + room.Size.Height
		d.Pos.X = randRange(room.Pos.X+1, room.Pos.X+room.Size.Width-1)
		d.Color = entity.KeyColor(randRange(0, entity.ColorCount))
		if shouldPlaceDoor() {
			room.Doors[entity.Bottom] = d
		}
	}
	if room.Connections[3] != nil {
		d := &entity.Door{}
		d.Pos.X = room.Pos.X
		d.Pos.Y = randRange(room.Pos.Y+1, room.Pos.Y+room.Size.Height-1)
		d.Color = entity.KeyColor(randRange(0, entity.ColorCount))
		if shouldPlaceDoor() {
			room.Doors[entity.Left] = d
		}
	}
}

func (g *generator) generateCorridors() {
	for _, room := range g.level.Rooms {
		if neighbor := room.Connections[entity.Right]; neighbor != nil {
			start := getConnectionPoint(room, entity.Right)
			end := getConnectionPoint(neighbor, entity.Left)
			g.createCorridor(start, end, true)
		}
		if neighbor := room.Connections[entity.Bottom]; neighbor != nil {
			start := getConnectionPoint(room, entity.Bottom)
			end := getConnectionPoint(neighbor, entity.Top)
			g.createCorridor(start, end, false)
		}
	}
}

func getConnectionPoint(r *entity.Room, dir entity.Direction) entity.Point {
	if r.Doors[dir] != nil {
		return r.Doors[dir].Pos
	}
	switch dir {
	case entity.Top:
		return entity.Point{X: r.Pos.X + r.Size.Width/2, Y: r.Pos.Y}
	case entity.Right:
		return entity.Point{X: r.Pos.X + r.Size.Width, Y: r.Pos.Y + r.Size.Height/2}
	case entity.Bottom:
		return entity.Point{X: r.Pos.X + r.Size.Width/2, Y: r.Pos.Y + r.Size.Height}
	case entity.Left:
		return entity.Point{X: r.Pos.X, Y: r.Pos.Y + r.Size.Height/2}
	default:
		panic("invalid dir")
	}
	return entity.Point{}
}

func (g *generator) createCorridor(start, end entity.Point, horizontal bool) {
	corridor := &entity.Corridor{}
	corridor.Points = append(corridor.Points, start)

	if horizontal {
		midX := start.X + (end.X-start.X)/2
		p1 := entity.Point{X: midX, Y: start.Y}
		p2 := entity.Point{X: midX, Y: end.Y}
		corridor.Points = append(corridor.Points, p1, p2)
	} else {
		midY := start.Y + (end.Y-start.Y)/2
		p1 := entity.Point{X: start.X, Y: midY}
		p2 := entity.Point{X: end.X, Y: midY}
		corridor.Points = append(corridor.Points, p1, p2)
	}
	corridor.Points = append(corridor.Points, end)
	g.level.Corridors = append(g.level.Corridors, corridor)
}

func (g *generator) renderToMap() {
	for _, r := range g.level.Rooms {
		for y := r.Pos.Y + 1; y < r.Pos.Y+r.Size.Height; y++ {
			for x := r.Pos.X + 1; x < r.Pos.X+r.Size.Width; x++ {
				g.level.SetTile(x, y, entity.TileFloor)
			}
		}
		for _, d := range r.Doors {
			if d != nil {
				g.level.SetTile(d.Pos.X, d.Pos.Y, entity.TileDoorLocked)
				g.level.Doors[d.Pos] = d
			}
		}
	}
	for _, c := range g.level.Corridors {
		for i := 0; i < len(c.Points)-1; i++ {
			drawLineCorridor(g.level, c.Points[i], c.Points[i+1])
		}
	}
}

func (g *generator) placeStartAndExit() {
	startRoom := g.level.Rooms[rand.Intn(len(g.level.Rooms))]
	endRoom := g.level.Rooms[len(g.level.Rooms)-1]
	for endRoom == startRoom {
		endRoom = g.level.Rooms[rand.Intn(len(g.level.Rooms))]
	}

	g.level.StartPos = center(startRoom)
	g.level.ExitPos = center(endRoom)

	g.startRoom = startRoom

	g.level.SetTile(g.level.ExitPos.X, g.level.ExitPos.Y, entity.TileDownStairs)
}

func (g *generator) generateKeysAndLocks() {
	accessibleRooms := make(map[*entity.Room]bool)
	queue := []*entity.Room{g.startRoom}

	accessibleRooms[g.startRoom] = true

	hasKey := make(map[entity.KeyColor]bool)
	hasKey[entity.NoKey] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for i, neighbor := range curr.Connections {
			if neighbor == nil || accessibleRooms[neighbor] {
				continue
			}

			g.fixLock(hasKey, curr.Doors[i], accessibleRooms)
			g.fixLock(hasKey, neighbor.Doors[opposite(i)], accessibleRooms)

			accessibleRooms[neighbor] = true
			queue = append(queue, neighbor)
		}
	}
}

func (g *generator) fixLock(hasKey map[entity.KeyColor]bool, door *entity.Door, accessibleRooms map[*entity.Room]bool) {
	doorColor := entity.NoKey
	if door != nil {
		doorColor = door.Color
	}

	if !hasKey[doorColor] {
		g.placeKeyInAccessibleArea(doorColor, accessibleRooms)

		hasKey[doorColor] = true
	}
}

func (g *generator) placeKeyInAccessibleArea(color entity.KeyColor, accessibleRooms map[*entity.Room]bool) {
	candidates := make([]*entity.Room, 0, len(accessibleRooms))
	for r := range accessibleRooms {
		candidates = append(candidates, r)
	}

	if len(candidates) == 0 {
		panic("can not place key in accessible rooms")
	}

	targetRoom := candidates[rand.Intn(len(candidates))]

	keyName := "Unknown Key"
	switch color {
	case entity.RedKey:
		keyName = "Red Key"
	case entity.BlueKey:
		keyName = "Blue Key"
	case entity.YellowKey:
		keyName = "Yellow Key"
	}

	key := &entity.Key{
		BaseItem: entity.BaseItem{
			Name: keyName,
			Pos:  getRandomPlaceInRoom(g.level, targetRoom),
		},
		Color: color,
	}

	g.level.Items[key.Pos] = append(g.level.Items[key.Pos], key)
}

func getRandomPlaceInRoom(level *entity.Level, room *entity.Room) (point entity.Point) {
	x := rand.Intn(room.Size.Width-2) + room.Pos.X + 1
	y := rand.Intn(room.Size.Height-2) + room.Pos.Y + 1

	for level.GetTile(x, y) == entity.TileDownStairs || (x == level.StartPos.X && y == level.StartPos.Y) {
		x = rand.Intn(room.Size.Width-2) + room.Pos.X + 1
		y = rand.Intn(room.Size.Height-2) + room.Pos.Y + 1
	}

	point.X = x
	point.Y = y

	return point
}

func opposite(dir int) int {
	return (dir + 2) % 4
}

func randRange(min, max int) int {
	if min >= max {
		return min
	}
	return rand.Intn(max-min) + min
}

func shouldPlaceDoor() bool {
	return rand.Float64() < entity.DoorChance
}

func center(r *entity.Room) entity.Point {
	return entity.Point{
		X: r.Pos.X + r.Size.Width/2,
		Y: r.Pos.Y + r.Size.Height/2,
	}
}

func drawLineCorridor(lvl *entity.Level, p1, p2 entity.Point) {
	if p1.X == p2.X {
		minY, maxY := min(p1.Y, p2.Y), max(p1.Y, p2.Y)
		for y := minY; y <= maxY; y++ {
			if lvl.GetTile(p1.X, y) == entity.TileWall {
				lvl.SetTile(p1.X, y, entity.TileCorridor)
			}
		}
	} else {
		minX, maxX := min(p1.X, p2.X), max(p1.X, p2.X)
		for x := minX; x <= maxX; x++ {
			if lvl.GetTile(x, p1.Y) == entity.TileWall {
				lvl.SetTile(x, p1.Y, entity.TileCorridor)
			}
		}
	}
}
