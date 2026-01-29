package domain

import (
	"fmt"
	"math/rand/v2"
	"time"
)

const (
	TreasureType ItemType = iota
	FoodType
	ElixirType
	ScrollType
	WeaponType
	NoType
)

const (
	HealthFeature int = iota
	DexterityFeature
	StrengthFeature
	NoFeature
)

type WeaponKind int

const (
	NoWeapon WeaponKind = iota
	Knife
	Sword
	Bow
	Phaser
	MachineGun
	Bomb
	NumWeapon
)

type Logo struct {
	Logo      rune
	ColorPair int16
}

type Item interface {
	GetItemType() ItemType
	GetX() int
	GetY() int
	GetLogo() Logo
	GetValue() int
	GetHealth() HealthType
	GetStrength() StrengthType
	GetFeature() FeatureType
	GetDuration() DurationType
	SetPosition(p Position)
	GetPosition() Position
	UseItem(game *GameSession)
	FromMap(data map[string]interface{}) error
}

type BaseItem struct {
	itype    ItemType
	position Position
	logo     Logo
}

func (w *BaseItem) FromMap(data map[string]interface{}) error {
	w.position = Position{}
	w.position.FromMap(data["position"].(map[string]interface{}))
	w.itype = ItemType(data["itype"].(float64))
	return nil
}

func (i *BaseItem) GetItemType() ItemType {
	return i.itype
}

func (i *BaseItem) GetX() int {
	return i.position.x
}

func (i *BaseItem) GetY() int {
	return i.position.y
}

func (i *BaseItem) GetLogo() Logo {
	return i.logo
}

func (i *BaseItem) GetValue() int {
	return 0
}

func (i *BaseItem) GetHealth() HealthType {
	return 0
}

func (i *BaseItem) GetStrength() StrengthType {
	return 0
}

func (i *BaseItem) GetFeature() FeatureType {
	return FeatureType(NoFeature)
}

func (i *BaseItem) GetDuration() DurationType {
	return 0
}

func (i *BaseItem) SetPosition(p Position) {
	i.position = p
}

func (i *BaseItem) GetPosition() Position {
	return i.position
}

func ItemToMap(i Item) map[string]interface{} {
	return map[string]interface{}{
		"type":     "Item",
		"itype":    i.GetItemType(),
		"position": i.GetPosition().ToMap(),
		"health":   i.GetHealth(),
		"strength": i.GetStrength(),
		"feature":  i.GetFeature(),
		"duration": i.GetDuration(),
		"value":    i.GetValue(),
	}
}

func ItemFromMap(data map[string]interface{}) (Item, error) {
	itemType, ok := data["itype"].(float64)
	if !ok {
		return nil, fmt.Errorf("item type not found")
	}

	var item Item
	switch ItemType(itemType) {
	case TreasureType:
		item = &Treasure{}
	case WeaponType:
		item = &Weapon{}
	case FoodType:
		item = &Food{}
	case ElixirType:
		item = &Elixir{}
	case ScrollType:
		item = &Scroll{}
	default:
		item = &BaseItem{}
	}

	if err := item.FromMap(data); err != nil {
		return nil, err
	}

	return item, nil
}

type Treasure struct {
	BaseItem
	value int
}

func (w *Treasure) FromMap(data map[string]interface{}) error {
	w.BaseItem.FromMap(data)
	w.value = int(data["value"].(float64))
	return nil
}

func (t *Treasure) GetValue() int {
	return t.value
}

type Food struct {
	BaseItem
	health HealthType
}

func (w *Food) FromMap(data map[string]interface{}) error {
	w.BaseItem.FromMap(data)
	w.health = HealthType(data["health"].(float64))
	return nil
}

func (f *Food) GetHealth() HealthType {
	return f.health
}

type Weapon struct {
	BaseItem
	strength StrengthType
}

func (w *Weapon) FromMap(data map[string]interface{}) error {
	w.BaseItem.FromMap(data)
	w.strength = StrengthType(data["strength"].(float64))
	return nil
}

func (w *Weapon) GetStrength() StrengthType {
	return w.strength
}

type Elixir struct {
	BaseItem
	feature  FeatureType
	value    int
	duration DurationType
}

func (w *Elixir) FromMap(data map[string]interface{}) error {
	w.BaseItem.FromMap(data)
	w.feature = FeatureType(data["feature"].(float64))
	w.value = int(data["value"].(float64))
	w.duration = DurationType(data["duration"].(float64))
	return nil
}

func (e *Elixir) GetFeature() FeatureType {
	return e.feature
}

func (e *Elixir) GetValue() int {
	return e.value
}

func (e *Elixir) GetDuration() DurationType {
	return e.duration
}

type Scroll struct {
	BaseItem
	feature FeatureType
	value   int
}

func (w *Scroll) FromMap(data map[string]interface{}) error {
	w.BaseItem.FromMap(data)
	w.feature = FeatureType(data["feature"].(float64))
	w.value = int(data["value"].(float64))
	return nil
}

func (s *Scroll) GetFeature() FeatureType {
	return s.feature
}

func (s *Scroll) GetValue() int {
	return s.value
}

func (f *Food) UseItem(game *GameSession) {
	game.Hero.SetHealth(f.GetHealth())
}

func (s *Scroll) UseItem(game *GameSession) {
	value := s.value
	switch s.feature {
	case FeatureType(HealthFeature):
		game.Hero.maximum += HealthType(value)
		game.Hero.health += HealthType(value)
	case FeatureType(DexterityFeature):
		game.Hero.dexterity += DexterityType(value)
	case FeatureType(StrengthFeature):
		game.Hero.strength += StrengthType(value)
	}
}

func (e *Elixir) UseItem(game *GameSession) {
	value := e.value
	switch e.feature {
	case FeatureType(HealthFeature):
		health_buff := BuffType{value: value, effect: int(e.duration), start: int(time.Now().Unix())}
		game.Hero.elixir_buffs.health = append(game.Hero.elixir_buffs.health, health_buff)
	case FeatureType(DexterityFeature):
		dexterity_buff := BuffType{value: value, effect: int(e.duration), start: int(time.Now().Unix())}
		game.Hero.elixir_buffs.dexterity = append(game.Hero.elixir_buffs.dexterity, dexterity_buff)
	case FeatureType(StrengthFeature):
		strength_buff := BuffType{value: value, effect: int(e.duration), start: int(time.Now().Unix())}
		game.Hero.elixir_buffs.strength = append(game.Hero.elixir_buffs.strength, strength_buff)
	}
}

func (w *Weapon) UseItem(game *GameSession) {
	prevWeapon := game.Hero.weapon
	game.Hero.weapon = w
	for i, item := range game.Hero.backpack.items[WeaponType] {
		if item.GetStrength() == w.GetStrength() {
			game.Hero.backpack.items[WeaponType] = append(game.Hero.backpack.items[WeaponType][:i], game.Hero.backpack.items[WeaponType][i+1:]...)
			break
		}
	}
	game.DropWeapon(prevWeapon)
}

func (b *BaseItem) UseItem(game *GameSession) {}

func MakeItem(t ItemType, f FeatureType) Item {
	var i Item
	// Default logo
	logo := Logo{Logo: '*', ColorPair: 5}
	switch t {
	case TreasureType:
		base := BaseItem{itype: t, logo: logo}
		i = &Treasure{
			BaseItem: base,
			value:    10 + rand.IntN(90), // 10-99
		}
	case FoodType:
		base := BaseItem{itype: t, logo: logo}
		i = &Food{
			BaseItem: base,
			health:   HealthType(5 + rand.IntN(15)), // 5-19
		}
	case WeaponType:
		base := BaseItem{itype: t, logo: logo}
		i = &Weapon{
			BaseItem: base,
			strength: StrengthType(5 + rand.IntN(10)), // 5-14
		}
	case ElixirType:
		base := BaseItem{itype: t, logo: logo}
		i = &Elixir{
			BaseItem: base,
			feature:  f,
			value:    10 + rand.IntN(20),
			duration: DurationType(10 + rand.IntN(20)),
		}
	case ScrollType:
		base := BaseItem{itype: t, logo: logo}
		i = &Scroll{
			BaseItem: base,
			feature:  f,
			value:    10 + rand.IntN(20),
		}
	}
	return i
}
