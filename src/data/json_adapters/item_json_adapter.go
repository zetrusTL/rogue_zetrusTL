package json_adapters

import (
	"encoding/json"
	"fmt"

	"example.com/rogue/domain/entity"
)

type jsonItemWrapper struct {
	Type entity.ItemType `json:"type"`
	Data json.RawMessage `json:"data"`
}

func newItemWrapper(item entity.Item) jsonItemWrapper {
	data, _ := json.Marshal(item)
	return jsonItemWrapper{
		Type: item.GetType(),
		Data: data,
	}
}

func (w jsonItemWrapper) unwrap() (entity.Item, error) {
	return restoreItem(w.Type, w.Data)
}

func restoreItem(t entity.ItemType, data json.RawMessage) (entity.Item, error) {
	switch t {
	case entity.ItemFood:
		var i entity.Food
		err := json.Unmarshal(data, &i)
		return &i, err
	case entity.ItemWeapon:
		var i entity.Weapon
		err := json.Unmarshal(data, &i)
		return &i, err
	case entity.ItemElixir:
		var i entity.Elixir
		err := json.Unmarshal(data, &i)
		return &i, err
	case entity.ItemScroll:
		var i entity.Scroll
		err := json.Unmarshal(data, &i)
		return &i, err
	case entity.ItemKey:
		var i entity.Key
		err := json.Unmarshal(data, &i)
		return &i, err
	case entity.ItemTreasure:
		var i entity.Treasure
		err := json.Unmarshal(data, &i)
		return &i, err
	default:
		return nil, fmt.Errorf("unknown item type: %v", t)
	}
}
