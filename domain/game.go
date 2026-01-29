package domain

import (
	"fmt"
	"math"
	"rogue/domain/serializer"
	"rogue/log"
	"rogue/utils"
	"time"
)

type ItemType int
type EnemyType int
type ItemSubtype int
type HealthType int
type DexterityType int
type StrengthType int
type ValueType int
type HostilityType int
type ExperienceType int
type HitsType int
type SizeType int
type ActionTime int
type ActionType string
type FeatureType int
type DurationType int
type DirectionType int

type Statistics struct {
	level      int
	SessFile   string
	treasures  ValueType
	enemies    int
	food       int
	elixirs    int
	scrolls    int
	hits       int
	enemy_hits int
	tiles      int
	state      int //GAMEOVER_S WIN_S OTHER
}

func MakeGameSession() GameSession {
	return GameSession{}
}

func (s Statistics) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":       "Statistics",
		"sess_file":  s.SessFile,
		"treasures":  s.treasures,
		"enemies":    s.enemies,
		"food":       s.food,
		"elixirs":    s.elixirs,
		"scrolls":    s.scrolls,
		"hits":       s.hits,
		"enemy_hits": s.enemy_hits,
		"tiles":      s.tiles,
		"state":      s.state,
		"level":      s.level,
	}
}

func (s *Statistics) FromMap(data map[string]interface{}) error {
	if v, ok := data["sess_file"].(string); ok {
		s.SessFile = v
	}
	if v, ok := data["level"].(float64); ok {
		s.level = int(v)
	}
	if v, ok := data["state"].(float64); ok {
		s.state = int(v)
	}
	if v, ok := data["treasures"].(float64); ok {
		s.treasures = ValueType(v)
	}
	if v, ok := data["enemies"].(float64); ok {
		s.enemies = int(v)
	}
	if v, ok := data["food"].(float64); ok {
		s.food = int(v)
	}
	if v, ok := data["elixirs"].(float64); ok {
		s.elixirs = int(v)
	}
	if v, ok := data["scrolls"].(float64); ok {
		s.scrolls = int(v)
	}
	if v, ok := data["hits"].(float64); ok {
		s.hits = int(v)
	}
	if v, ok := data["enemy_hits"].(float64); ok {
		s.enemy_hits = int(v)
	}
	if v, ok := data["tiles"].(float64); ok {
		s.tiles = int(v)
	}
	return nil
}

type GameSession struct {
	Stats      []Statistics
	Playground *Map
	Level      Level
	Hero       Player
	FightEnemy *Enemy
	Armistice  bool
	TimeCode   int64
	CurStats   Statistics
	Action     int
}

func (game *GameSession) HandleUserInput(action int) {
	game.Action = action
}

func (game *GameSession) LoadFromSession(index int) error {
	sessions := GetGameSessions()
	if index > len(sessions)-1 {
		log.DebugLog(fmt.Sprintf("сессий с индексом %d не существует", index))
		return fmt.Errorf("сессий с индексом %d не существует", index)
	}

	log.DebugLog("load game", utils.GetSessName(sessions[index].TimeCode))

	game.FillGame(&sessions[index])
	game.GenerateFog()
	return nil
}

func (game *GameSession) NewGame() {
	game.Playground = &Map{}
	game.Hero = MakePlayer()
	game.NextLevel(true)
	log.DebugLog("NewGame:", game.Level.value)
}

func (game *GameSession) GetSession() *GameSession {
	return game
}

func (game *GameSession) Store() {
	log.DebugLog("sess file", utils.GetSessName(game.TimeCode))
	serial := serializer.MakeSerializer()
	fileHandler := serializer.MakeFileHandler("sessions")
	if err := fileHandler.SaveObject(utils.GetSessName(game.TimeCode), game, serial); err != nil {
		log.DebugLog("Ошибка сохранения: ", err)
	}
}

func (game *GameSession) CalcStats() Statistics {
	return Statistics{
		SessFile:  utils.GetSessName(game.TimeCode),
		level:     game.Level.GetValue(),
		treasures: game.Hero.backpack.treasure,
		state:     game.CurStats.state,
	}
}

func (game *GameSession) GenerateFog() {
	if FOG == 1 {
		for y := range game.Level.Fog {
			for x := range game.Level.Fog[y] {
				game.Level.Fog[y][x] = 1
			}
		}
	}
}

func (game *GameSession) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":       "GameSession",
		"stats":      game.CurStats.ToMap(),
		"playground": game.Playground.ToMap(),
		"level":      game.Level.ToMap(),
		"hero":       game.Hero.ToMap(),
		"time_code":  game.TimeCode,
	}
}

func (game *GameSession) FromMap(data map[string]interface{}) error {
	// game.TimeCode = 0
	if v, ok := data["time_code"].(float64); ok {
		game.TimeCode = int64(v)
	}

	if playgroundData, ok := data["playground"].(map[string]interface{}); ok {
		game.Playground = &Map{}
		if err := game.Playground.FromMap(playgroundData); err != nil {
			return err
		}
	}
	if levelData, ok := data["level"].(map[string]interface{}); ok {
		game.Level = Level{}
		if err := game.Level.FromMap(levelData); err != nil {
			return err
		}
	}
	if heroData, ok := data["hero"].(map[string]interface{}); ok {
		game.Hero = Player{}
		if err := game.Hero.FromMap(heroData); err != nil {
			return err
		}
	}
	if statsData, ok := data["stats"].(map[string]interface{}); ok {
		game.CurStats = Statistics{}
		if err := game.CurStats.FromMap(statsData); err != nil {
			return err
		}
	}
	return nil
}

func (game *GameSession) FillGame(src *GameSession) {
	game.Playground = src.Playground
	game.Level = src.Level
	game.Hero = src.Hero
	game.Stats = src.Stats
	game.TimeCode = src.TimeCode
	game.CurStats = src.CurStats
}

func (game *GameSession) NextLevel(new_game bool) {
	if new_game {
		game.Level.value = 0
	}
	game.Level.value++
	game.GenLevel()
	game.Playground.RoomsToMap(&game.Level)
	game.Hero.SetPosition(game.Level.Rooms[START_ROOM].Y()+1, game.Level.Rooms[START_ROOM].X()+1)
	game.GenerateFog()

	if game.Level.value > 1 {
		if game.TimeCode == 0 {
			game.TimeCode = time.Now().Unix()
		}
		game.Store()
	}
}

func (game *GameSession) Contact() (*Enemy, bool) {
	x := game.Hero.position.x
	y := game.Hero.position.y
	for _, e := range game.Level.Enemies {
		if math.Abs(float64(e.position.x-x))+math.Abs(float64(e.position.y-y)) == 1 {
			e.invisibility = 0
			return e, true
		}
	}
	return nil, false
}

func (game *GameSession) GenLevel() {
	rooms := GenRooms()
	game.Level.Rooms = ConnectRooms(rooms)

	corridors, doors := GenCorridors(game.Level.Rooms)
	door := getExit(rooms)
	doors = append(doors, door)
	game.Level.Exit = &door
	game.Level.Corridors = corridors
	game.Level.Doors = doors
	game.Level.PopulateLevel(game.Level.Rooms)
	log.DebugLog("GenLevel:", game.Level.value)
}

func (level *Level) PopulateLevel(rooms []Room) {
	log.DebugLog("PopulateLevel:", level.value)
	var positions []Position
	positions = append(positions, Position{level.Rooms[START_ROOM].Y() + 1, level.Rooms[START_ROOM].X() + 1})
	positions = append(positions, Position{level.Rooms[START_ROOM].Y() + 1, level.Rooms[START_ROOM].X() + 1})
	level.Items = GenItems(rooms, level.value, positions)
	level.Enemies = GenEnemies(rooms, level.value, positions)
}

func (game *GameSession) UpdateGameState() {
	fsm := Fsm_instance()
	state := fsm.Get_cur_state()

	if fsm.states[state].exit != nil {
		fsm.states[state].exit(game)
	}

	state = fsm.states[state].action(game)

	if fsm.states[state].enter != nil {
		fsm.states[state].enter(game)
	}

	fsm.cur_state = state
	game.CurStats.state = state
}

func (game *GameSession) Atack(enemy *Enemy) (bool, bool) {
	enemy.atack(game)
	game.Hero.atack(game, enemy)
	return (enemy.health <= 0), (game.Hero.GetHealth() <= 0)
}

func (game *GameSession) RemoveEnemyByPtr(e *Enemy) {
	for i, enemy := range game.Level.Enemies {
		if enemy == e {
			game.Level.Enemies = append(game.Level.Enemies[:i], game.Level.Enemies[i+1:]...)
		}
	}
}

func (game *GameSession) MakeTreasure(e *Enemy) {
	for _, enemy := range game.Level.Enemies {
		if enemy == e {
			game.Level.Items = append(game.Level.Items, MakeItem(TreasureType, 0))
			game.Level.Items[len(game.Level.Items)-1].SetPosition(e.position)
			log.DebugLog("new tresure:", game.Level.Items[len(game.Level.Items)-1].GetPosition(), game.Level.Items[len(game.Level.Items)-1])
		}
	}
}

func (game *GameSession) GetItem(x, y int) (Item, bool) {
	for i, item := range game.Level.Items {
		if item.GetX() == x && item.GetY() == y {
			game.Level.Items = append(game.Level.Items[:i], game.Level.Items[i+1:]...)
			return item, true
		}
	}
	return nil, false
}

func (game *GameSession) PutIntoBag(item Item) bool {
	if game.Hero.backpack.items == nil {
		game.Hero.backpack.items = make(map[ItemType][]Item)
	}
	if _, ok := game.Hero.backpack.items[item.GetItemType()]; !ok {
		game.Hero.backpack.items[item.GetItemType()] = make([]Item, 0, 9)
	}
	added := false
	if len(game.Hero.backpack.items[item.GetItemType()]) < 9 {
		game.Hero.backpack.items[item.GetItemType()] = append(game.Hero.backpack.items[item.GetItemType()], item)
		added = true
	}
	return added
}

func (game *GameSession) TakeItem() bool {
	x := game.Hero.position.x
	y := game.Hero.position.y
	if item, ok := game.GetItem(x, y); ok {
		if item.GetItemType() == TreasureType {
			game.Hero.backpack.treasure = game.Hero.backpack.treasure + ValueType(item.GetValue())
			return true
		}
		return game.PutIntoBag(item)
	}
	return false
}

func (game *GameSession) GetEnemyAt(x, y int) *Enemy {
	for _, e := range game.Level.Enemies {
		if e.position.x == x && e.position.y == y {
			return e
		}
	}
	return nil
}

func (game *GameSession) CanMove(y, x int) bool {
	cells := []int{INNER_AREA_CHAR, CORRIDOR_CHAR, DOOR_CHAR}
	nc := game.Playground.Playground[y][x]
	enemies := game.Level.Enemies
	for _, enemy := range enemies {
		if enemy.position.x == x && enemy.position.y == y {
			return false
		}
	}
	return contains(cells, nc)
}

func (game *GameSession) ReachExit() bool {
	return game.Level.Exit.position == game.Hero.position
}

func (game *GameSession) MoveHero(dy, dx int) bool {
	success := false
	ny := game.Hero.position.y + dy
	nx := game.Hero.position.x + dx
	if game.CanMove(ny, nx) {
		success = true
		game.Hero.position.y = ny
		game.Hero.position.x = nx
		game.CurStats.AddTile()
	}
	return success
}

func (game *GameSession) MoveEnemies() {
	for _, enemy := range game.Level.Enemies {
		if enemy.enemy_type != OgreType || !enemy.freeze {
			enemy.move(game)
		} else {
			log.DebugLog("freeze enemy:", enemy.enemy_type)
		}
		enemy.freeze = false
	}
}

func (game *GameSession) DropWeapon(w Item) {
	permittedCells := []int{INNER_AREA_CHAR, CORRIDOR_CHAR, DOOR_CHAR}
	x := game.Hero.position.x
	y := game.Hero.position.y
	var position Position
	switch {
	case contains(permittedCells, game.Playground.Playground[y+1][x]):
		position = Position{x, y + 1}
	case contains(permittedCells, game.Playground.Playground[y-1][x]):
		position = Position{x, y - 1}
	case contains(permittedCells, game.Playground.Playground[y][x+1]):
		position = Position{x + 1, y}
	case contains(permittedCells, game.Playground.Playground[y][x-1]):
		position = Position{x - 1, y}
	default:
		return
	}
	w.SetPosition(position)
	game.Level.Items = append(game.Level.Items, w)
	w.SetPosition(position)
}

func (game *GameSession) UseItem(itemtype ItemType, index int) {
	items := game.Hero.GetItemByType(itemtype)
	if index == 0 && itemtype != WeaponType {
		return
	}
	if index == 0 && itemtype == WeaponType {
		game.Hero.EmptyHand()
		return
	}
	if len(items) >= index {
		item := items[index-1]
		item.UseItem(game)
		game.Hero.backpack.items[itemtype] = append(items[:index-1], items[index:]...)
		game.CurStats.AddItem(item)
	}
}

func (stats *Statistics) AddItem(item Item) {
	switch item.GetItemType() {
	case ElixirType:
		stats.elixirs += 1
	case TreasureType:
		stats.treasures += ValueType(item.GetValue())
	case FoodType:
		stats.food += 1
	case ScrollType:
		stats.scrolls += 1
	}

}

func (stats *Statistics) AddDefeatedEnemy() { //враг побежден
	stats.enemies += 1
}

func (stats *Statistics) AddTile() {
	stats.tiles += 1
}

func (stats *Statistics) AddEnemyHit() { //враг атакует
	stats.enemy_hits += 1
}

func (stats *Statistics) AddHit() { //игрок атакует
	stats.hits += 1
}

func (stats *Statistics) GetTreasures() ValueType {
	return stats.treasures
}

func (stats *Statistics) GetLevel() int {
	return stats.level
}

func (stats *Statistics) GetResult() string {
	switch stats.state {
	case WIN_S:
		return "WIN"
	case GAMEOVER_S:
		return "GAMEOVER"
	default:
		return "OTHER"
	}
}

func (stats *Statistics) GetElixirs() int {
	return stats.elixirs
}

func (stats *Statistics) GetEnemies() int {
	return stats.enemies
}

func (stats *Statistics) GetTiles() int {
	return stats.tiles
}

func (stats *Statistics) GetEnemyHits() int {
	return stats.enemy_hits
}

func (stats *Statistics) GetHits() int {
	return stats.hits
}

// If you have scrolls count field
func (stats *Statistics) GetScrolls() int {
	return stats.scrolls
}

// If you have food count field
func (stats *Statistics) GetFood() int {
	return stats.food
}
