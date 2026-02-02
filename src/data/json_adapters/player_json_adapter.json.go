package json_adapters

import (
	"encoding/json"

	"example.com/rogue/domain/entity"
)

type JsonPlayerWrapper struct {
	*entity.Player
}

type PlayerJsonData struct {
	CharacterStats entity.CharacterStats
	Backpack       backpackJsonData
	Weapon         *entity.Weapon
	ActiveBuffs    []entity.Buff
	Facing         entity.Direction
	SleepTurns     int
}

type backpackJsonData struct {
	Items map[entity.ItemType][]jsonItemWrapper
	Gold  int
}

func (w *JsonPlayerWrapper) MarshalJSON() ([]byte, error) {
	if w.Player == nil {
		return []byte("null"), nil
	}
	p := w.Player

	activeBuffs := make([]entity.Buff, 0, len(w.ActiveBuffs))
	for _, buff := range w.ActiveBuffs {
		activeBuffs = append(activeBuffs, *buff)
	}

	temp := PlayerJsonData{
		CharacterStats: p.CharacterStats,
		Weapon:         p.Weapon,
		ActiveBuffs:    activeBuffs,
		Facing:         p.Facing,
		SleepTurns:     p.SleepTurns,
		Backpack: backpackJsonData{
			Gold:  p.Backpack.Gold,
			Items: make(map[entity.ItemType][]jsonItemWrapper),
		},
	}

	for t, items := range p.Backpack.Items {
		wrappers := make([]jsonItemWrapper, 0)
		for _, item := range items {
			wrappers = append(wrappers, newItemWrapper(item))
		}
		temp.Backpack.Items[t] = wrappers
	}

	return json.Marshal(temp)
}

func (w *JsonPlayerWrapper) UnmarshalJSON(data []byte) error {
	var temp PlayerJsonData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if w.Player == nil {
		w.Player = &entity.Player{Backpack: entity.NewBackpack()}
	}
	p := w.Player

	activeBuffs := make([]*entity.Buff, 0, len(temp.ActiveBuffs))
	for i := range temp.ActiveBuffs {
		activeBuffs = append(activeBuffs, &temp.ActiveBuffs[i])
	}

	p.CharacterStats = temp.CharacterStats
	p.Weapon = temp.Weapon
	p.ActiveBuffs = activeBuffs
	p.Facing = temp.Facing
	p.SleepTurns = temp.SleepTurns
	p.Backpack.Gold = temp.Backpack.Gold

	p.Backpack.Items = make(map[entity.ItemType][]entity.Item)

	for t, wrappers := range temp.Backpack.Items {
		items := make([]entity.Item, 0)
		for _, wrap := range wrappers {
			item, err := wrap.unwrap()
			if err != nil {
				return err
			}
			items = append(items, item)
		}
		p.Backpack.Items[t] = items
	}

	return nil
}
