package json_adapters

import (
	"encoding/json"

	"example.com/rogue/domain/entity"
)

type JsonLevelWrapper struct {
	*entity.Level
}

type LevelJsonData struct {
	LevelNum  int
	MapWidth  int
	MapHeight int
	Tiles     [][]entity.TileType
	Rooms     []RoomJsonData
	Corridors []*entity.Corridor
	StartPos  entity.Point
	ExitPos   entity.Point

	// Поля, которые требуют конвертации
	Monsters map[string]*entity.Monster   `json:"Monsters"`
	Items    map[string][]jsonItemWrapper `json:"Items"`
}

type RoomJsonData struct {
	Object  entity.Object
	GridPos entity.Point

	Doors [4]*entity.Door

	Connections [4]*entity.Point

	IsLocked bool
}

func (w *JsonLevelWrapper) MarshalJSON() ([]byte, error) {
	if w.Level == nil {
		return []byte("null"), nil
	}

	temp := LevelJsonData{
		LevelNum:  w.Level.LevelNum,
		MapWidth:  entity.MapWidth,
		MapHeight: entity.MapHeight,
		Tiles:     w.Level.Tiles,
		Rooms:     convertRoomsToJson(w.Level.Rooms),
		Corridors: w.Level.Corridors,
		StartPos:  w.Level.StartPos,
		ExitPos:   w.Level.ExitPos,
		Monsters:  make(map[string]*entity.Monster),
		Items:     make(map[string][]jsonItemWrapper),
	}

	for p, m := range w.Level.Monsters {
		temp.Monsters[pointToString(p)] = m
	}
	for p, items := range w.Level.Items {
		var wrappers []jsonItemWrapper
		for _, item := range items {
			wrappers = append(wrappers, newItemWrapper(item))
		}
		temp.Items[pointToString(p)] = wrappers
	}

	return json.Marshal(temp)
}

func convertRoomsToJson(rooms []*entity.Room) []RoomJsonData {
	result := make([]RoomJsonData, 0, len(rooms))
	for _, room := range rooms {
		var connections [4]*entity.Point
		roomJson := RoomJsonData{
			Object:   room.Object,
			GridPos:  room.GridPos,
			Doors:    room.Doors,
			IsLocked: room.IsLocked,
		}

		for i, connectedRoom := range room.Connections {
			if connectedRoom != nil {
				connections[i] = &connectedRoom.GridPos
			}
		}
		roomJson.Connections = connections

		result = append(result, roomJson)
	}

	return result
}

func (w *JsonLevelWrapper) UnmarshalJSON(data []byte) error {
	var temp LevelJsonData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if w.Level == nil {
		w.Level = &entity.Level{}
	}

	w.Level.LevelNum = temp.LevelNum
	w.Level.Width = temp.MapWidth
	w.Level.Height = temp.MapHeight
	w.Level.Tiles = temp.Tiles
	w.Level.Rooms = convertJsonToRooms(temp.Rooms)
	w.Level.Corridors = temp.Corridors
	w.Level.StartPos = temp.StartPos
	w.Level.ExitPos = temp.ExitPos

	w.Level.Monsters = make(map[entity.Point]*entity.Monster)
	w.Level.Items = make(map[entity.Point][]entity.Item)
	w.Level.Doors = make(map[entity.Point]*entity.Door)

	for k, m := range temp.Monsters {
		w.Level.Monsters[stringToPoint(k)] = m
	}

	for _, room := range temp.Rooms {
		for _, door := range room.Doors {
			if door != nil {
				w.Level.Doors[door.Pos] = door
			}
		}
	}

	for k, wrappers := range temp.Items {
		p := stringToPoint(k)
		var items []entity.Item
		for _, wrap := range wrappers {
			item, err := wrap.unwrap()
			if err != nil {
				return err
			}
			items = append(items, item)
		}
		w.Level.Items[p] = items
	}

	return nil
}

func convertJsonToRooms(jsonRooms []RoomJsonData) []*entity.Room {
	var mapGridPosRoom = make(map[entity.Point]*entity.Room)

	result := make([]*entity.Room, 0, len(jsonRooms))
	for _, jsonRoom := range jsonRooms {
		room := &entity.Room{
			Object:   jsonRoom.Object,
			GridPos:  jsonRoom.GridPos,
			IsLocked: jsonRoom.IsLocked,
			Doors:    jsonRoom.Doors,
		}

		mapGridPosRoom[room.GridPos] = room

		result = append(result, room)
	}

	for i, room := range result {
		for j, connection := range jsonRooms[i].Connections {
			if connection == nil {
				continue
			}
			room.Connections[j] = mapGridPosRoom[*connection]
		}

	}

	return result
}
