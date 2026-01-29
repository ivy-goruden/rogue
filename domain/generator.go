package domain

import (
	"math/rand/v2"
)

func checkGraph(rooms []Room) bool {
	visited := make([]*Room, 0)   //посещенные вершины графа
	toVisit := []*Room{&rooms[0]} //вершины, которые предстоит посетить
	for len(toVisit) > 0 {
		temp := make([]*Room, 0)
		for _, e := range toVisit {
			visited = append(visited, e)
			neighbours := e.connections
			for _, p := range neighbours {
				if !contains(visited, p) && !contains(temp, p) {
					temp = append(temp, p)
				}
			}

		}
		toVisit = temp
	}
	return len(visited) == MAX_ROOMS_NUMBER
}

func getExit(rooms []Room) Door {
	room := &rooms[START_ROOM]
	visited := []*Room{room}
	exitRoom, _ := bsf(room, visited, 0)
	for i := range rooms {
		if &rooms[i] == exitRoom {
			// Определяем строку в сетке
			if i/ROOMS_PER_SIDE == 0 {
				// Первая строка - дверь сверху
				x := exitRoom.dimentions.x + exitRoom.dimentions.width/2
				y := exitRoom.dimentions.y
				return Door{
					position: Position{x: x, y: y},
				}
			} else if i/ROOMS_PER_SIDE == ROOMS_PER_SIDE-1 {
				// Последняя строка - дверь снизу
				x := exitRoom.dimentions.x + exitRoom.dimentions.width/2
				y := exitRoom.dimentions.y + exitRoom.dimentions.height - 1
				return Door{
					position: Position{x: x, y: y},
				}
			} else {
				// Определяем столбец
				col := i % ROOMS_PER_SIDE
				if col == 0 {
					// Первый столбец - дверь слева
					x := exitRoom.dimentions.x
					y := exitRoom.dimentions.y + exitRoom.dimentions.height/2
					return Door{
						position: Position{x: x, y: y},
					}
				} else if col == ROOMS_PER_SIDE-1 {
					// Последний столбец - дверь справа
					x := exitRoom.dimentions.x + exitRoom.dimentions.width - 1
					y := exitRoom.dimentions.y + exitRoom.dimentions.height/2
					return Door{
						position: Position{x: x, y: y},
					}
				}
			}
			break
		}
	}

	return Door{
		position: Position{
			x: exitRoom.dimentions.x + exitRoom.dimentions.width/2,
			y: exitRoom.dimentions.y + exitRoom.dimentions.height/2,
		},
	}

}

func bsf(room *Room, visited []*Room, size int) (*Room, int) {
	neighs := room.connections
	if len(neighs) == 0 {
		return room, size
	}
	var max_room *Room
	var max int
	for _, e := range neighs {
		if !contains(visited, e) {
			temp := append(visited, e)
			max_room_temp, max_temp := bsf(e, temp, size+1)
			if max < max_temp {
				max = max_temp
				max_room = max_room_temp
			}
		}
	}
	if max_room == nil {
		return room, size
	}
	return max_room, max
}

func randInRange(min int, max int) int {
	return rand.IntN(max-min+1) + min
}

func contains[T comparable](arr []T, key T) bool {
	for _, e := range arr {
		if e == key {
			return true
		}
	}
	return false
}

func GenRooms() []Room {
	rooms := make([]Room, 0)
	for i := range MAX_ROOMS_NUMBER {
		room_w := randInRange(MIN_ROOM_W, MAX_ROOM_W)
		room_h := randInRange(MIN_ROOM_H, MAX_ROOM_H)
		room_range_x := ROOM_RANGE_W * (i % ROOMS_PER_SIDE)
		room_range_y := ROOM_RANGE_H * int(i/ROOMS_PER_SIDE)
		room_x := randInRange(room_range_x+1, room_range_x+ROOM_RANGE_W-1-room_w)
		room_y := randInRange(room_range_y+1, room_range_y+ROOM_RANGE_H-1-room_h)
		room := Room{
			dimentions: Dimention{
				width:  room_w,
				height: room_h,
				x:      room_x,
				y:      room_y,
			},
		}
		rooms = append(rooms, room)
	}
	return rooms
}

func addNeighbour(room1 *Room, room2 *Room) {
	if !contains(room1.connections, room2) {
		room1.connections = append(room1.connections, room2)
		// Добавляем обратное соединение
	}
	if !contains(room2.connections, room1) {
		room2.connections = append(room2.connections, room1)
	}
}

func connectRooms(rooms []Room) []Room {
	stop := false
	for !stop {
		for i := 0; i < MAX_ROOMS_NUMBER; i++ {
			rooms[i].connections = make([]*Room, 0, 4)
		}
		for i := 0; i < MAX_ROOMS_NUMBER; i++ {
			row := int(i / ROOMS_PER_SIDE)
			col := i % ROOMS_PER_SIDE

			// Нижний сосед
			if row < ROOMS_PER_SIDE-1 && rand.Float64() < CORIDOR_CHANCE {
				targetRoom := &rooms[i+ROOMS_PER_SIDE]
				addNeighbour(&rooms[i], targetRoom)
			}
			// Правый сосед
			if col < ROOMS_PER_SIDE-1 && rand.Float64() < CORIDOR_CHANCE {
				targetRoom := &rooms[i+1]
				addNeighbour(&rooms[i], targetRoom)
			}
		}
		stop = checkGraph(rooms)
	}

	return rooms

}

func ConnectRooms(rooms []Room) []Room {
	return connectRooms(rooms)
}

func genCorridors(rooms []Room) ([]Corridor, []Door) {
	var corridors []Corridor
	var doors []Door
	for i := range rooms {
		for i2 := range rooms {
			if i2 <= i {
				continue
			}

			if !contains(rooms[i].connections, &rooms[i2]) {
				continue
			}

			//left
			if i+1 == i2 {
				// Коридор от правой стены первой комнаты до левой стены второй
				corr := Corridor{
					startX: rooms[i].dimentions.x + rooms[i].dimentions.width,
					startY: rooms[i].dimentions.y + rooms[i].dimentions.height/2, // середина по высоте
					endX:   rooms[i2].dimentions.x,                               // левая стена комнаты i2
					endY:   rooms[i2].dimentions.y + rooms[i2].dimentions.height/2,
				}
				corridors = append(corridors, corr)

				// Дверь в правой стене комнаты i
				door := Door{
					position: Position{
						x: rooms[i].dimentions.x + rooms[i].dimentions.width - 1,
						y: rooms[i].dimentions.y + rooms[i].dimentions.height/2,
					},
				}
				doors = append(doors, door)

				// Дверь в левой стене комнаты i2
				door = Door{
					position: Position{
						x: rooms[i2].dimentions.x,
						y: rooms[i2].dimentions.y + rooms[i2].dimentions.height/2,
					},
				}
				doors = append(doors, door)
			}
			//up
			if i+ROOMS_PER_SIDE == i2 {
				corr := Corridor{
					startX: rooms[i].dimentions.x + rooms[i].dimentions.width/2,
					startY: rooms[i].dimentions.y + rooms[i].dimentions.height,
					endX:   rooms[i2].dimentions.x + rooms[i2].dimentions.width/2,
					endY:   rooms[i2].dimentions.y - 1,
				}
				corridors = append(corridors, corr)
				door := Door{
					position: Position{
						x: rooms[i].dimentions.x + rooms[i].dimentions.width/2,
						y: rooms[i].dimentions.y + rooms[i].dimentions.height - 1,
					},
				}
				doors = append(doors, door)
				door = Door{
					position: Position{
						x: rooms[i2].dimentions.x + rooms[i2].dimentions.width/2,
						y: rooms[i2].dimentions.y - 1 + 1,
					},
				}
				doors = append(doors, door)
			}
		}
	}

	return corridors, doors
}

func GenCorridors(rooms []Room) ([]Corridor, []Door) {
	return genCorridors(rooms)
}

func GenPosition(r *Room) Position {
	_x := randInRange(r.dimentions.x+1, r.dimentions.x+r.dimentions.width-2)
	_y := randInRange(r.dimentions.y+1, r.dimentions.y+r.dimentions.height-2)
	return Position{x: _x, y: _y}
}

func GenItems(rooms []Room, level int, positions []Position) []Item {
	var items []Item
	count := 0
	for index, room := range rooms {
		if index == START_ROOM {
			// continue
		}
		if count >= MAX_ITEMS_TOTAL {
			break
		}
		chance := rand.Float64()
		if chance < ITEM_SPAWN_CHANCE-float64(level)*SUBCHANCE_PER_LEVEL {
			count++
			Itemtype := randInRange(0, int(NoType)-1)
			Feature := randInRange(0, NoFeature-1)
			items = append(items, MakeItem(ItemType(Itemtype), FeatureType(Feature)))
			var position Position
			stop := false
			for !stop {
				position = GenPosition(&room)
				stop = true
				for _, pos := range positions {
					if position.x == pos.x && position.y == pos.y {
						stop = false
						break
					}
				}
			}
			items[len(items)-1].SetPosition(position)
			positions = append(positions, position)
		}
	}
	return items
}

func GenEnemies(rooms []Room, level int, positions []Position) []*Enemy {
	var enemies []*Enemy
	count := 0
	for index, room := range rooms {
		if index == START_ROOM {
			continue
		}
		if count >= MAX_ENEMIES_TOTAL {
			break
		}
		chance := rand.Float64()
		if chance < ENEMY_SPAWN_CHANCE+float64(level)*SUMCHANCE_PER_LEVEL {
			count++
			enemy_type := randInRange(0, 5)
			enemies = append(enemies, MakeEnemy(EnemyType(enemy_type)))
			var position Position
			stop := false
			for !stop {
				position = GenPosition(&room)
				stop = true
				for _, pos := range positions {
					if position.x == pos.x && position.y == pos.y {
						stop = false
						break
					}
				}
			}
			enemies[len(enemies)-1].position = position
			enemies[len(enemies)-1].birth_room = &room
			positions = append(positions, position)
		}
	}
	return enemies
}
