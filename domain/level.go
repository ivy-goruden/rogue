package domain

import (
	"math"
	"rogue/utils"
)

type EntityType int
type GridSymbol int

const ChacterEntity EntityType = 0
const EnemyEntity EntityType = 1
const ItemEntity EntityType = 2

type Position struct {
	x, y int
}

func (p Position) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type": "Position",
		"x":    p.x,
		"y":    p.y,
	}
}

func (p *Position) FromMap(data map[string]interface{}) error {
	p.x = int(data["x"].(float64))
	p.y = int(data["y"].(float64))
	return nil
}

type Dimention struct {
	x, y, width, height int
}

func (d Dimention) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"x":      d.x,
		"y":      d.y,
		"width":  d.width,
		"height": d.height,
	}
}

func (d *Dimention) FromMap(data map[string]interface{}) error {
	d.x = int(data["x"].(float64))
	d.y = int(data["y"].(float64))
	d.width = int(data["width"].(float64))
	d.height = int(data["height"].(float64))
	return nil
}

type Level struct {
	value     int
	Exit      *Door
	Doors     []Door
	Corridors []Corridor
	Rooms     []Room
	Items     []Item
	Enemies   []*Enemy
	Fog       [MAP_HEIGHT][MAP_WIDTH]int //туман войны
}

func (l Level) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":      "Level",
		"value":     l.value,
		"exit":      l.Exit.ToMap(),
		"doors":     utils.ArrayMap(l.Doors, func(x Door) map[string]interface{} { return x.ToMap() }),
		"corridors": utils.ArrayMap(l.Corridors, func(x Corridor) map[string]interface{} { return x.ToMap() }),
		"rooms":     utils.ArrayMap(l.Rooms, func(x Room) map[string]interface{} { return x.ToMap() }),
		"items":     utils.ArrayMap(l.Items, func(x Item) map[string]interface{} { return ItemToMap(x) }),
		"enemies":   utils.ArrayMap(l.Enemies, func(x *Enemy) map[string]interface{} { return x.ToMap() }),
		"fog":       l.Fog,
	}
}

func (l *Level) FromMap(data map[string]interface{}) error {
	l.value = int(data["value"].(float64))

	if v, ok := data["fog"].([]interface{}); ok {
		if mat, ok := ToMatrix(v); ok {
			l.Fog = mat
		}
	}

	if v, ok := data["exit"].(map[string]interface{}); ok {
		l.Exit = &Door{}
		l.Exit.FromMap(v)
	}

	var arr []interface{}

	arr, _ = data["doors"].([]interface{})
	l.Doors = utils.ArrayMap(arr,
		func(x interface{}) Door {
			m, _ := x.(map[string]interface{})
			var d = Door{}
			d.FromMap(m)
			return d
		})

	arr, _ = data["corridors"].([]interface{})
	l.Corridors = utils.ArrayMap(arr,
		func(x interface{}) Corridor {
			m, _ := x.(map[string]interface{})
			var d = Corridor{}
			d.FromMap(m)
			return d
		})

	arr, _ = data["rooms"].([]interface{})
	l.Rooms = utils.ArrayMap(arr,
		func(x interface{}) Room {
			m, _ := x.(map[string]interface{})
			var d = Room{}
			d.FromMap(m)
			return d
		})

	arr, _ = data["items"].([]interface{})
	l.Items = utils.ArrayMap(arr,
		func(x interface{}) Item {
			m, _ := x.(map[string]interface{})
			i, _ := ItemFromMap(m)
			return i
		})

	arr, _ = data["enemies"].([]interface{})
	l.Enemies = utils.ArrayMap(arr,
		func(x interface{}) *Enemy {
			m, _ := x.(map[string]interface{})
			var d = &Enemy{}
			d.FromMap(m)
			return d
		})

	return nil
}

type Door struct {
	position Position
}

func (d Door) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":     "Door",
		"position": d.position.ToMap(),
	}
}

func (d *Door) FromMap(data map[string]interface{}) error {
	d.position = Position{}
	if v, ok := data["position"].(map[string]interface{}); ok {
		d.position.FromMap(v)
	}
	return nil
}

type Room struct {
	index       int
	dimentions  Dimention
	connections []*Room
	doors       [4]Door
}

func (r Room) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":       "Room",
		"index":      r.index,
		"dimentions": r.dimentions.ToMap(),
		"doors":      r.doors,
	}
}

func (r *Room) FromMap(data map[string]interface{}) error {
	r.dimentions = Dimention{}
	r.dimentions.FromMap(data["dimentions"].(map[string]interface{}))
	r.index = int(data["index"].(float64))

	arr, _ := data["doors"].([]interface{})
	doors := utils.ArrayMap(arr,
		func(x interface{}) Door {
			m, _ := x.(map[string]interface{})
			var d = Door{}
			d.FromMap(m)
			return d
		})
	r.doors = [4]Door{doors[0], doors[1], doors[2], doors[3]}

	return nil
}

type Corridor struct {
	startX int
	startY int
	endX   int
	endY   int
}

func (c Corridor) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"type":   "Corridor",
		"startX": c.startX,
		"startY": c.startY,
		"endX":   c.endX,
		"endY":   c.endY,
	}
}

func (c *Corridor) FromMap(data map[string]interface{}) error {
	c.startX = int(data["startX"].(float64))
	c.startY = int(data["startY"].(float64))
	c.endX = int(data["endX"].(float64))
	c.endY = int(data["endY"].(float64))
	return nil
}

type Entity struct {
	entity_type EntityType
	position    Position
	symbol      int
}

type Map struct {
	Playground [MAP_HEIGHT][MAP_WIDTH]int
}

func (m *Map) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Playground": m.Playground,
	}
}

func ToMatrix(data []interface{}) ([30][90]int, bool) {
	if len(data) != 30 {
		return [30][90]int{}, false
	}
	var mat [30][90]int
	for i := 0; i < 30; i++ {
		row, ok := data[i].([]interface{})
		if !ok || len(row) != 90 {
			return [30][90]int{}, false
		}
		for j := 0; j < 90; j++ {
			val, ok := row[j].(float64)
			if !ok {
				return [30][90]int{}, false
			}
			mat[i][j] = int(val)
		}
	}
	return mat, true
}

func (m *Map) FromMap(data map[string]interface{}) error {
	if v, ok := data["Playground"].([]interface{}); ok {
		if mat, ok := ToMatrix(v); ok {
			m.Playground = mat
		}
	}
	return nil
}

func (r Room) X() int      { return r.dimentions.x }
func (r Room) Y() int      { return r.dimentions.y }
func (r Room) Width() int  { return r.dimentions.width }
func (r Room) Height() int { return r.dimentions.height }

func (d Door) X() int { return d.position.x }
func (d Door) Y() int { return d.position.y }

func (c Corridor) StartX() int { return c.startX }
func (c Corridor) StartY() int { return c.startY }
func (c Corridor) EndX() int   { return c.endX }
func (c Corridor) EndY() int   { return c.endY }

func (m *Map) RoomsToMap(level *Level) {
	grid := make([][]rune, MAP_HEIGHT)
	for i := range grid {
		grid[i] = make([]rune, MAP_WIDTH)
		for j := range grid[i] {
			grid[i][j] = EMPTY_CHAR
		}
	}

	DrawRooms(level.Rooms, &grid)
	DrawCorridors(level.Corridors, &grid)
	DrawDoors(level.Doors, &grid)

	for i, row := range grid {
		for j, cell := range row {
			m.Playground[i][j] = int(cell)
		}
	}
}

func DrawRooms(rooms []Room, grid_p *[][]rune) {
	grid := *grid_p
	// Draw each room
	for i := range rooms {
		room := &rooms[i]
		x := room.X()
		y := room.Y()
		w := room.Width()
		h := room.Height()

		// Draw room walls
		for row := y; row < y+h; row++ {
			for col := x; col < x+w; col++ {
				if row < 0 || row >= MAP_HEIGHT || col < 0 || col >= MAP_WIDTH {
					continue
				}

				// Draw walls
				if row == y || row == y+h-1 {
					grid[row][col] = WALL_CHAR_H
				} else if col == x || col == x+w-1 {
					grid[row][col] = WALL_CHAR_V
				} else {
					grid[row][col] = INNER_AREA_CHAR
				}
			}
		}

		// Draw corners
		if y >= 0 && y < MAP_HEIGHT && x >= 0 && x < MAP_WIDTH {
			grid[y][x] = WALL_CHAR_TOPLEFT
		}
		if y >= 0 && y < MAP_HEIGHT && x+w-1 >= 0 && x+w-1 < MAP_WIDTH {
			grid[y][x+w-1] = WALL_CHAR_TOPRIGHT
		}
		if y+h-1 >= 0 && y+h-1 < MAP_HEIGHT && x >= 0 && x < MAP_WIDTH {
			grid[y+h-1][x] = WALL_CHAR_BOTTOMLEFT
		}
		if y+h-1 >= 0 && y+h-1 < MAP_HEIGHT && x+w-1 >= 0 && x+w-1 < MAP_WIDTH {
			grid[y+h-1][x+w-1] = WALLCHAR_BOTTOMRIGHT
		}
		//DrawRoomNumber(room, i, grid_p)
	}
}

func DrawRoomNumber(room *Room, roomIndex int, grid_p *[][]rune) {
	grid := *grid_p
	x := room.X()
	y := room.Y()
	w := room.Width()
	h := room.Height()
	centerX := x + w/2
	centerY := y + h/2
	if centerY < 0 || centerY >= MAP_HEIGHT || centerX < 0 || centerX >= MAP_WIDTH {
		return
	}
	if centerY >= 0 && centerY < MAP_HEIGHT && centerX >= 0 && centerX < MAP_WIDTH {
		grid[centerY][centerX] = rune('0' + roomIndex)
	}
}

func DrawDoors(doors []Door, grid_p *[][]rune) {
	grid := *grid_p
	for _, door := range doors {
		grid[door.Y()][door.X()] = DOOR_CHAR
	}

}

func DrawCorridors(corridors []Corridor, grid_p *[][]rune) {
	for _, coor := range corridors {
		DrawCorridor(coor.StartX(), coor.StartY(), coor.EndX(), coor.EndY(), grid_p)
	}
}

func DrawCorridor(startX, startY, endX, endY int, grid_p *[][]rune) {
	// Безопасное получение grid
	if grid_p == nil || len(*grid_p) == 0 {
		return
	}
	grid := *grid_p
	dx := math.Abs(float64(endX - startX))
	dy := math.Abs(float64(endY - startY))
	sx := 1
	if startX > endX {
		sx = -1
	}
	sy := 1
	if startY > endY {
		sy = -1
	}
	err := dx - dy
	x, y := startX, startY
	for {
		if y >= 0 && y < len(grid) && x >= 0 && x < len(grid[0]) {
			grid[y][x] = CORRIDOR_CHAR
		}
		if x == endX && y == endY {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			grid[y][x] = CORRIDOR_CHAR
			x += sx
		}
		if e2 < dx {
			err += dx
			grid[y][x] = CORRIDOR_CHAR
			y += sy
		}
	}
}

func (entity *Entity) GetX() int {
	return entity.position.x
}

func (entity *Entity) GetY() int {
	return entity.position.y
}

func (level Level) GetValue() int {
	return level.value
}
