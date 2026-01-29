package domain

import (
	"math"
	"math/rand/v2"
	"rogue/log"
)

type EnemyInter interface{}
type VampireInter interface{}
type GhostInter interface{}
type OgreInter interface{}
type SnakeMageInter interface{}
type ZombieInter interface{}
type MimicInter interface{}

const (
	ZombieType EnemyType = iota
	VampireType
	GhostType
	OgreType
	SnakeMageType
	MimicType
	HumanType
)

type EnemyBehavior interface {
	move(*Enemy, *GameSession)
	atack(*Enemy, *GameSession)
	String() string
}

const (
	ForwardDir DirectionType = iota
	BackDir
	LeftDir
	RightDir
	ForwardLeftDir
	ForwardRightDir
	BackLeftDir
	BackRightDir
	StopDir
)

type Enemy struct {
	birth_room *Room
	enemy_type EnemyType
	behavior   EnemyBehavior
	health     HealthType
	dexterity  DexterityType
	strength   StrengthType
	hostility  HostilityType
	// направление текущее, чтобы при выборе следующего направления
	// не выбрать случайным образом текущее
	direction DirectionType
	// положение текущее
	position Position
	// преследует или нет
	chasing bool
	// видимый или нет
	visible bool
	// ударил или нет
	hit bool
	// время невидимости
	invisibility int
	// количество атак на монстра
	atack_counter int
	// состояние обездвиживания
	freeze bool
}

func (b *Enemy) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"birth_room":    b.birth_room.ToMap(),
		"enemy_type":    b.enemy_type,
		"health":        b.health,
		"dexterity":     b.dexterity,
		"strength":      b.strength,
		"hostility":     b.hostility,
		"direction":     b.direction,
		"position":      b.position.ToMap(),
		"chasing":       b.chasing,
		"visible":       b.visible,
		"hit":           b.hit,
		"invisibility":  b.invisibility,
		"atack_counter": b.atack_counter,
		"freeze":        b.freeze,
	}
}

func (b *Enemy) FromMap(data map[string]interface{}) error {
	b.birth_room = &Room{}
	b.birth_room.FromMap(data["birth_room"].(map[string]interface{}))

	b.enemy_type = EnemyType(data["enemy_type"].(float64))
	b.health = HealthType(data["health"].(float64))
	b.dexterity = DexterityType(data["dexterity"].(float64))
	b.strength = StrengthType(data["strength"].(float64))
	b.hostility = HostilityType(data["hostility"].(float64))
	b.direction = DirectionType(data["direction"].(float64))

	b.position = Position{}
	b.position.FromMap(data["position"].(map[string]interface{}))

	b.chasing = data["chasing"].(bool)
	b.visible = data["visible"].(bool)
	b.hit = data["hit"].(bool)
	b.invisibility = int(data["invisibility"].(float64))
	b.atack_counter = int(data["atack_counter"].(float64))
	b.freeze = data["freeze"].(bool)

	b.FixBehavior()
	return nil
}

func (b *Enemy) EnemyType() EnemyType {
	return b.enemy_type
}

type Zombie struct{}

func (Zombie) move(enemy *Enemy, game *GameSession) {
	for stop := false; !stop; {
		direction := randInRange(0, 3)
		y_ := 0
		x_ := 0
		switch direction {
		case 0:
			y_ = 0
			x_ = ENEMY_STEP
		case 1:
			y_ = 0
			x_ = -ENEMY_STEP
		case 2:
			y_ = ENEMY_STEP
			x_ = 0
		case 3:
			y_ = -ENEMY_STEP
			x_ = 0
		}
		if enemy.CanMove(enemy.position.y+y_, enemy.position.x+x_, game) && enemy.InARoom(enemy.position.y+y_, enemy.position.x+x_) {
			enemy.position.y += y_
			enemy.position.x += x_
			stop = true
		}
	}
}
func (Zombie) atack(enemy *Enemy, game *GameSession) {
	if enemy.hit_check(game) {
		game.Hero.health -= HealthType(enemy.strength)
	}
}
func (Zombie) String() string {
	return "Zombie"
}

type Vampire struct{}

func (Vampire) move(enemy *Enemy, game *GameSession) {
	for stop := false; !stop; {
		direction := randInRange(0, 3)
		y_ := 0
		x_ := 0
		switch direction {
		case 0:
			y_ = 0
			x_ = -ENEMY_STEP
		case 1:
			y_ = 0
			x_ = ENEMY_STEP
		case 2:
			y_ = ENEMY_STEP
			x_ = 0
		case 3:
			y_ = -ENEMY_STEP
			x_ = 0
		}
		if enemy.CanMove(enemy.position.y+y_, enemy.position.x+x_, game) {
			enemy.position.y += y_
			enemy.position.x += x_
			stop = true
		}
	}
}
func (Vampire) atack(enemy *Enemy, game *GameSession) {
	if enemy.hit_check(game) {
		game.Hero.health -= HealthType(enemy.strength)
		game.Hero.maximum -= 10
	}
}
func (Vampire) String() string {
	return "Vampire"
}

type Ghost struct{}

func (Ghost) move(enemy *Enemy, game *GameSession) {
	chance := rand.Float64()
	if chance < INVISIBILITY_CHANCE {
		if enemy.invisibility == 0 {
			enemy.invisibility = 1
		} else {
			enemy.invisibility = 0
		}
		return
	}
	for stop := false; !stop; {
		new_x := randInRange(enemy.birth_room.dimentions.x+1, enemy.birth_room.dimentions.x+enemy.birth_room.dimentions.width-2)
		new_y := randInRange(enemy.birth_room.dimentions.y+1, enemy.birth_room.dimentions.y+enemy.birth_room.dimentions.height-2)
		if enemy.CanMove(new_y, new_x, game) {
			enemy.position.x = new_x
			enemy.position.y = new_y
			stop = true
		}
	}
}
func (Ghost) atack(enemy *Enemy, game *GameSession) {
	if enemy.hit_check(game) {
		game.Hero.health -= HealthType(enemy.strength)
	}
}
func (Ghost) String() string {
	return "Ghost"
}

type Ogre struct{}

func (Ogre) move(enemy *Enemy, game *GameSession) {
	for stop := false; !stop; {
		direction := randInRange(0, 3)
		y_ := 0
		x_ := 0
		switch direction {
		case 0:
			y_ = 0
			x_ = -2
		case 1:
			y_ = 0
			x_ = 2
		case 2:
			y_ = 2
			x_ = 0
		case 3:
			y_ = -2
			x_ = 0
		}
		if enemy.CanMove(enemy.position.y+y_, enemy.position.x+x_, game) && enemy.InARoom(enemy.position.y+y_, enemy.position.x+x_) {
			enemy.position.y += y_
			enemy.position.x += x_
			stop = true
		}
	}
}
func (Ogre) atack(enemy *Enemy, game *GameSession) {
	if enemy.hit_check(game) {
		game.Hero.health -= HealthType(enemy.strength)
	}
}
func (Ogre) String() string {
	return "Ogre"
}

type SnakeMage struct{}

func (SnakeMage) move(enemy *Enemy, game *GameSession) {
	for stop := false; !stop; {
		direction := randInRange(0, 3)
		y_ := 0
		x_ := 0
		switch direction {
		case 0:
			y_ = ENEMY_STEP
			x_ = -ENEMY_STEP
		case 1:
			y_ = ENEMY_STEP
			x_ = ENEMY_STEP
		case 2:
			y_ = -ENEMY_STEP
			x_ = ENEMY_STEP
		case 3:
			y_ = -ENEMY_STEP
			x_ = -ENEMY_STEP
		}
		if enemy.CanMove(enemy.position.y+y_, enemy.position.x+x_, game) {
			enemy.position.y += y_
			enemy.position.x += x_
			stop = true
		}
	}
}
func (SnakeMage) atack(enemy *Enemy, game *GameSession) {
	if enemy.hit_check(game) {
		game.Hero.health -= HealthType(enemy.strength)
		if int(rand.Float64()*100) > 25 {
			game.Hero.put_to_sleep = true
		}
	}
}
func (SnakeMage) String() string {
	return "SnakeMage"
}

type Mimic struct{}

func (Mimic) move(enemy *Enemy, game *GameSession) {

}
func (Mimic) atack(enemy *Enemy, game *GameSession) {
	if enemy.hit_check(game) {
		game.Hero.health -= HealthType(enemy.strength)
	}
}
func (Mimic) String() string {
	return "Mimic"
}

func MakeEnemy(t EnemyType) *Enemy {
	e := &Enemy{enemy_type: t}
	e.FixBehavior()
	switch t {
	case ZombieType:
		e.dexterity = 25
		e.health = 60
		e.strength = 50
		e.hostility = 4
	case VampireType:
		e.dexterity = 60
		e.health = 75
		e.strength = 50
		e.hostility = 5
	case GhostType:
		e.dexterity = 60
		e.health = 25
		e.strength = 25
		e.hostility = 3
	case OgreType:
		e.dexterity = 25
		e.health = 100
		e.strength = 100
		e.hostility = 4
	case SnakeMageType:
		e.dexterity = 90
		e.health = 50
		e.strength = 50
		e.hostility = 5
	case MimicType:
		e.dexterity = 0
		e.health = 100
		e.strength = 0
		e.hostility = 2
	}
	return e
}

func (e *Enemy) FixBehavior() {
	switch e.enemy_type {
	case ZombieType:
		e.behavior = Zombie{}
	case VampireType:
		e.behavior = Vampire{}
	case GhostType:
		e.behavior = Ghost{}
	case OgreType:
		e.behavior = Ogre{}
	case SnakeMageType:
		e.behavior = SnakeMage{}
	case MimicType:
		e.behavior = Mimic{}
	}
}

func (enemy *Enemy) String() string {
	return enemy.behavior.String()
}

func (enemy *Enemy) hit_check(game *GameSession) bool {
	if rand.Float64()*100 > float64(game.Hero.GetDexterity()) {
		game.CurStats.AddEnemyHit()
		return true
	}
	return false
}

func (enemy *Enemy) atack(game *GameSession) {
	enemy.behavior.atack(enemy, game)
}

func (enemy *Enemy) move(game *GameSession) {
	can, next_pos := enemy.IsChase(game)
	if !can {
		enemy.chasing = false
		enemy.behavior.move(enemy, game)
	} else {
		enemy.position = next_pos
		enemy.chasing = true
	}
}

func (enemy *Enemy) GetY() int {
	return enemy.position.y
}

func (enemy *Enemy) GetX() int {
	return enemy.position.x
}

func (enemy *Enemy) CanMove(y, x int, game *GameSession) bool {
	if x < 0 || y < 0 || x >= MAP_WIDTH || y >= MAP_HEIGHT {
		return false
	}
	if enemy.freeze {
		log.DebugLog("freeze enemy:", enemy)
		enemy.freeze = false
		return false
	}
	cells := []int{INNER_AREA_CHAR, CORRIDOR_CHAR, DOOR_CHAR}
	nc := game.Playground.Playground[y][x]
	game.Hero.GetX()
	game.Hero.GetY()
	if y == game.Hero.GetY() && x == game.Hero.GetX() {
		return false
	}
	return contains(cells, nc)
}

func (enemy *Enemy) IsChase(game *GameSession) (bool, Position) {
	radius := enemy.hostility * HOSTILITY_VALUE
	distanceX := math.Abs(float64(game.Hero.GetX() - enemy.GetX()))
	distanceY := math.Abs(float64(game.Hero.GetY() - enemy.GetY()))
	euclideanDistance := math.Sqrt(float64(distanceX*distanceX + distanceY*distanceY))
	if euclideanDistance > 0 && euclideanDistance <= float64(radius) {
		log.DebugLog("folowing you")
		return enemy.CanCatch(game, radius)
	}
	return false, Position{}
}

func (enemy *Enemy) CanCatch(game *GameSession, radius HostilityType) (bool, Position) {
	startX, startY := enemy.GetX(), enemy.GetY()
	endX, endY := game.Hero.GetX(), game.Hero.GetY()

	// Проверяем начальную и конечную позиции
	startPos := Position{startX, startY}
	endPos := Position{endX, endY}

	// Если уже на одной позиции
	if startPos == endPos {
		return false, Position{} // Уже на месте
	}

	// Проверяем допустимость ячеек
	cells := []int{INNER_AREA_CHAR, CORRIDOR_CHAR, DOOR_CHAR}

	// Используем BFS для поиска пути
	visited := make(map[Position]bool)
	parent := make(map[Position]Position)
	queue := []Position{startPos}
	visited[startPos] = true

	found := false
	var current Position

	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]

		// Если достигли цели
		if current == endPos {
			found = true
			break
		}

		// Проверяем расстояние по прямой от текущей позиции до врага
		distanceX := math.Abs(float64(current.x - startX))
		distanceY := math.Abs(float64(current.y - startY))
		euclideanDistance := math.Sqrt(distanceX*distanceX + distanceY*distanceY)

		// Если вышли за радиус враждебности
		if euclideanDistance > float64(radius) {
			continue
		}

		// Проверяем все направления
		directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

		for _, dir := range directions {
			newX, newY := current.x+dir[0], current.y+dir[1]
			newPos := Position{newX, newY}

			// Проверяем, не посещали ли уже эту позицию
			if visited[newPos] {
				continue
			}

			// Проверяем валидность ячейки
			if !isValidCell(newX, newY) {
				continue
			}

			// Проверяем тип ячейки
			cellType := game.Playground.Playground[newY][newX]
			if !contains(cells, cellType) {
				continue
			}

			// Добавляем в очередь
			visited[newPos] = true
			parent[newPos] = current
			queue = append(queue, newPos)
		}
	}

	// Если путь найден, восстанавливаем первый шаг
	if found {
		// Восстанавливаем путь от цели к началу
		path := []Position{}
		for pos := endPos; pos != startPos; pos = parent[pos] {
			path = append([]Position{pos}, path...)
		}

		// Возвращаем true и следующую позицию для движения
		if len(path) > 0 {
			return true, path[0] // Первая позиция после старта
		}
	}

	return false, Position{}
}

func (enemy *Enemy) InARoom(y, x int) bool {
	if enemy.birth_room.dimentions.x <= x && x <= enemy.birth_room.dimentions.x+enemy.birth_room.dimentions.width-1 &&
		enemy.birth_room.dimentions.y <= y && y <= enemy.birth_room.dimentions.y+enemy.birth_room.dimentions.height-1 {
		return true
	}
	return false
}

func isValidCell(x, y int) bool {
	return x >= 0 && x < MAP_WIDTH &&
		y >= 0 && y < MAP_HEIGHT
}

func (enemy *Enemy) IsVisible() bool {
	return enemy.invisibility == 0
}

func (enemy *Enemy) SetVisible() {
	enemy.invisibility = 0
}

func (enemy *Enemy) IsChasing() bool {
	return enemy.chasing
}
