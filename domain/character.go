package domain

import (
	"fmt"
	"math/rand/v2"
	"rogue/utils"
	"strconv"
	"time"
)

type Backpack struct {
	items    map[ItemType][]Item
	treasure ValueType
}

func (b *Backpack) ToMap() map[string]interface{} {
	items := map[string]interface{}{}
	for itype, arr := range b.items {
		itype_str := fmt.Sprintf("%d", itype)
		items[itype_str] = utils.ArrayMap(arr, func(x Item) map[string]interface{} { return ItemToMap(x) })
	}
	return map[string]interface{}{
		"type":     "Backpack",
		"items":    items,
		"treasure": b.treasure,
	}
}

func (b *Backpack) FromMap(data map[string]interface{}) error {
	b.treasure = ValueType(data["treasure"].(float64))

	b.items = make(map[ItemType][]Item)

	arr, _ := data["items"].(map[string]interface{})
	for itype_str, v := range arr {
		items_int, _ := v.([]interface{})
		items := utils.ArrayMap(items_int,
			func(x interface{}) Item {
				m, _ := x.(map[string]interface{})
				i, _ := ItemFromMap(m)
				return i
			})
		it, err := strconv.Atoi(itype_str)
		if err != nil {
			panic(err)
		}
		b.items[ItemType(it)] = items
	}
	return nil
}

func MakeBackpack() Backpack {
	return Backpack{
		items:    make(map[ItemType][]Item),
		treasure: 0,
	}
}

type BuffType struct {
	value  int
	effect int
	start  int
}

func (b BuffType) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":   "BuffType",
		"value":  b.value,
		"effect": b.effect,
		"start":  b.start,
	}
}

func (b *BuffType) FromMap(data map[string]interface{}) error {
	b.value = int(data["value"].(float64))
	b.effect = int(data["effect"].(float64))
	b.start = int(data["start"].(float64))
	return nil
}

type BuffsType struct {
	health    []BuffType
	dexterity []BuffType
	strength  []BuffType
}

func (b BuffsType) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":      "BuffsType",
		"dexterity": utils.ArrayMap(b.dexterity, func(x BuffType) map[string]interface{} { return x.ToMap() }),
		"health":    utils.ArrayMap(b.health, func(x BuffType) map[string]interface{} { return x.ToMap() }),
		"strength":  utils.ArrayMap(b.strength, func(x BuffType) map[string]interface{} { return x.ToMap() }),
	}
}

func (b *BuffsType) FromMap(data map[string]interface{}) error {
	var arr []interface{}

	arr, _ = data["dexterity"].([]interface{})
	b.dexterity = utils.ArrayMap(arr,
		func(x interface{}) BuffType {
			m, _ := x.(map[string]interface{})
			var d = BuffType{}
			d.FromMap(m)
			return d
		})

	arr, _ = data["health"].([]interface{})
	b.health = utils.ArrayMap(arr,
		func(x interface{}) BuffType {
			m, _ := x.(map[string]interface{})
			var d = BuffType{}
			d.FromMap(m)
			return d
		})

	arr, _ = data["strength"].([]interface{})
	b.strength = utils.ArrayMap(arr,
		func(x interface{}) BuffType {
			m, _ := x.(map[string]interface{})
			var d = BuffType{}
			d.FromMap(m)
			return d
		})

	return nil
}

type Player struct {
	elixir_buffs BuffsType
	maximum      HealthType
	dexterity    DexterityType
	strength     StrengthType
	health       HealthType
	weapon       Item
	position     Position
	backpack     Backpack
	put_to_sleep bool
}

func (Player Player) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"type":         "Player",
		"elixir_buffs": Player.elixir_buffs.ToMap(),
		"maximum":      Player.maximum,
		"dexterity":    Player.dexterity,
		"strength":     Player.strength,
		"health":       Player.health,
		"position":     Player.position.ToMap(),
		"backpack":     Player.backpack.ToMap(),
		"put_to_sleep": Player.put_to_sleep,
	}
	if Player.weapon != nil {
		data["weapon"] = ItemToMap(Player.weapon)
	} else {
		data["weapon"] = nil
	}
	return data
}

func (Player *Player) FromMap(data map[string]interface{}) error {
	Player.elixir_buffs = BuffsType{}
	Player.elixir_buffs.FromMap(data["elixir_buffs"].(map[string]interface{}))
	Player.maximum = HealthType(data["maximum"].(float64))
	Player.dexterity = DexterityType(data["dexterity"].(float64))
	Player.strength = StrengthType(data["strength"].(float64))
	Player.health = HealthType(data["health"].(float64))
	Player.position = Position{}
	Player.position.FromMap(data["position"].(map[string]interface{}))
	Player.backpack = Backpack{}
	Player.backpack.FromMap(data["backpack"].(map[string]interface{}))
	Player.put_to_sleep = data["put_to_sleep"].(bool)

	Player.weapon, _ = ItemFromMap(data["weapon"].(map[string]interface{}))

	return nil
}

func (Player *Player) AddBuff(feature FeatureType, value, effect int) {
	if feature == FeatureType(HealthFeature) && len(Player.elixir_buffs.health) <= 9 {
		Player.elixir_buffs.health = append(Player.elixir_buffs.health, BuffType{
			value:  value,
			effect: effect,
			start:  int(time.Now().Unix()),
		})
		Player.maximum += HealthType(value)
	}
	if feature == FeatureType(DexterityFeature) && len(Player.elixir_buffs.dexterity) <= 9 {
		Player.elixir_buffs.dexterity = append(Player.elixir_buffs.dexterity, BuffType{
			value:  value,
			effect: effect,
			start:  int(time.Now().Unix()),
		})
		Player.dexterity += DexterityType(value)
	}
	if feature == FeatureType(StrengthFeature) && len(Player.elixir_buffs.strength) <= 9 {
		Player.elixir_buffs.strength = append(Player.elixir_buffs.strength, BuffType{
			value:  value,
			effect: effect,
			start:  int(time.Now().Unix()),
		})
		Player.strength += StrengthType(value)
	}
}

func (Player *Player) CalcBuffs() {
	t := int(time.Now().Unix())
	for _, b := range Player.elixir_buffs.health {
		if b.start+b.effect < t {
			Player.maximum -= HealthType(b.value)
		}
	}
	if Player.maximum <= 0 {
		Player.health = 0
		Player.maximum = 100
	}
	for _, b := range Player.elixir_buffs.dexterity {
		if b.start+b.effect < t {
			Player.dexterity -= DexterityType(b.value)
		}
	}
	for _, b := range Player.elixir_buffs.strength {
		if b.start+b.effect < t {
			Player.strength -= StrengthType(b.value)
		}
	}
}

func (Player *Player) UpdateBuffs() {
	buffs := BuffsType{}
	before := Player.GetHealth()
	for _, b := range Player.elixir_buffs.health {
		if b.start+b.effect > int(time.Now().Unix()) {
			buffs.health = append(buffs.health, b)
		}
	}

	for _, b := range Player.elixir_buffs.dexterity {
		if b.start+b.effect > int(time.Now().Unix()) {
			buffs.dexterity = append(buffs.dexterity, b)
		}
	}
	for _, b := range Player.elixir_buffs.strength {
		if b.start+b.effect > int(time.Now().Unix()) {
			buffs.strength = append(buffs.strength, b)
		}
	}
	Player.elixir_buffs = buffs
	after := Player.GetHealth()
	if before != after {
		if after <= 0 {
			Player.health = 1
		}
		if Player.GetMaxHealth() <= 0 {
			Player.maximum = 1
		}
	}
	Player.health = min(Player.health, HealthType(Player.GetMaxHealth()))
}

func MakePlayer() Player {
	pl := Player{
		maximum:   150,
		dexterity: 100,
		strength:  150,
		health:    150,
		backpack:  MakeBackpack(),
	}
	pl.weapon = &Weapon{strength: 0}

	pl.elixir_buffs = BuffsType{
		health:    make([]BuffType, 0, 9),
		dexterity: make([]BuffType, 0, 9),
		strength:  make([]BuffType, 0, 9),
	}
	return pl
}

func (Player *Player) Hx() int              { return Player.position.x }
func (Player *Player) Hy() int              { return Player.position.y }
func (Player *Player) Treasures() int       { return int(Player.backpack.treasure) }
func (Player *Player) SetPosition(y, x int) { Player.position = Position{x, y} }

func (Player *Player) GetBagItems() map[ItemType][]Item {
	return Player.backpack.items
}

func (Player *Player) atack(game *GameSession, enemy *Enemy) {
	e := game.FightEnemy
	if Player.hit_check(enemy) {
		e.atack_counter++
		game.CurStats.AddHit()
		if e.enemy_type != VampireType || (e.enemy_type == VampireType && e.atack_counter > 1) {
			e.health -= HealthType(game.Hero.strength)
		}
	}
}

func (Player *Player) hit_check(enemy *Enemy) bool {
	return DexterityType(rand.Float64()*100) > enemy.dexterity
}

func (Player *Player) IsDead() bool {
	return Player.health <= 0
}

func (Player *Player) GetX() int {
	return Player.position.x
}

func (Player *Player) GetY() int {
	return Player.position.y
}

func (Player *Player) GetHealth() int {
	health := int(Player.health)
	for _, b := range Player.elixir_buffs.health {
		health += b.value
	}
	return health
}

func (Player *Player) SetHealth(health HealthType) {
	Player.health = min(health+Player.health, Player.maximum)

}

func (Player *Player) GetStrength() int {
	str := int(Player.strength)
	if Player.weapon != nil {
		str += int(Player.weapon.GetStrength())
	}
	for _, b := range Player.elixir_buffs.strength {
		str += b.value
	}
	return str
}

func (Player *Player) GetDexterity() int {
	dex := int(Player.dexterity)
	for _, b := range Player.elixir_buffs.dexterity {
		dex += b.value
	}
	return dex
}

func (Player *Player) GetMaxHealth() int {
	max := int(Player.maximum)
	for _, b := range Player.elixir_buffs.health {
		max += b.value
	}
	return max
}

func (Player *Player) GetGold() int {
	return int(Player.backpack.treasure)
}

func (player Player) GetItemByType(itemType ItemType) []Item {
	arr := player.backpack.items[itemType]
	return arr[:]
}

func (player *Player) GetHandWeapon() Item {
	return player.weapon
}

func (player *Player) EmptyHand() {
	handWeapon := player.GetHandWeapon()
	if handWeapon != nil && handWeapon.GetStrength() != 0 && len(player.backpack.items[WeaponType]) < 9 {
		player.backpack.items[WeaponType] = append(player.backpack.items[WeaponType], handWeapon)
		player.weapon = &Weapon{strength: 0}
	}
}

func (player *Player) GetBuffs() BuffsType {
	return player.elixir_buffs
}

func (buffs *BuffsType) GetHealthBuffs() []BuffType {
	return buffs.health
}

func (buffs *BuffsType) GetDexterityBuffs() []BuffType {
	return buffs.dexterity
}

func (buffs *BuffsType) GetStrengthBuffs() []BuffType {
	return buffs.strength
}

func (buff BuffType) GetValue() int {
	return buff.value
}

func (buff BuffType) GetStart() int {
	return buff.start
}

func (buff BuffType) GetEffect() int {
	return buff.effect
}
