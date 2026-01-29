package domain

const (
	FOG = 1
)

const (
	UNINITIALIZED = -1
	CONNECTED     = 0
	NOT_CONNECTED = 1
)

const (
	MAP_HEIGHT     = 30
	MAP_WIDTH      = 90
	ROOMS_PER_SIDE = 3
	ROOM_RANGE_H   = MAP_HEIGHT / ROOMS_PER_SIDE
	ROOM_RANGE_W   = MAP_WIDTH / ROOMS_PER_SIDE
	MAX_ROOM_W     = ROOM_RANGE_W - 2
	MAX_ROOM_H     = ROOM_RANGE_H - 2
	MIN_ROOM_W     = 6
	MIN_ROOM_H     = 6
)

const (
	MAX_ROOMS_NUMBER     = ROOMS_PER_SIDE * ROOMS_PER_SIDE
	MAX_CORRIDORS_NUMBER = 12
	MAX_SEGMENTS         = 3
	CORRIDOR_WIDTH       = 2
)

const (
	SECTOR_HEIGHT = MAP_HEIGHT / ROOMS_PER_SIDE
	SECTOR_WIDTH  = MAP_WIDTH / ROOMS_PER_SIDE
)

const (
	CORNER_VERT_RANGE = (SECTOR_HEIGHT - 6) / 2
	CORNER_HOR_RANGE  = (SECTOR_WIDTH - 6) / 2
)

const (
	ROOM_CHANCE         = 0.5
	SPAWN_SET_CHANCE    = 0.5
	CORIDOR_CHANCE      = 0.5
	ENEMY_SPAWN_CHANCE  = 1.0 / MAX_ROOMS_NUMBER
	SUMCHANCE_PER_LEVEL = 0.1  //прибавляется к шансу спавна врагов на каждом уровне
	SUBCHANCE_PER_LEVEL = 0.01 //отнимается от шанса спавна предметов на каждом уровне
	ITEM_SPAWN_CHANCE   = 1.0  /// ROOMS_PER_SIDE
	INVISIBILITY_CHANCE = 0.2
)

const (
	MAX_ENEMIES_PER_ROOM  = 5
	MAX_ITEMS_PER_ROOM    = 5
	MAX_ENTITIES_PER_ROOM = MAX_ENEMIES_PER_ROOM + MAX_ITEMS_PER_ROOM + 2
	MAX_ENEMIES_TOTAL     = MAX_ENEMIES_PER_ROOM * ROOMS_PER_SIDE * ROOMS_PER_SIDE
	MAX_ITEMS_TOTAL       = MAX_ITEMS_PER_ROOM * ROOMS_PER_SIDE * ROOMS_PER_SIDE
	MAX_ENTITIES_TOTAL    = MAX_ENEMIES_TOTAL + MAX_ITEMS_TOTAL + 2
	ENEMY_POOL_LEN        = 26
	ITEM_POOL_LEN         = 5
	FOG_RADIUS            = 5
)

const (
	TOP    = 0
	RIGHT  = 1
	BOTTOM = 2
	LEFT   = 3
)

const (
	LEFT_TO_RIGHT_CORRIDOR = 0
	LEFT_TURN_CORRIDOR     = 1
	RIGHT_TURN_CORRIDOR    = 2
	TOP_TO_BOTTOM_CORRIDOR = 3
)

const (
	UNOCCUPIED = 0
	OCCUPIED   = 1
)

const (
	PLAYER      = 0
	EXIT        = 1
	ENEMY       = 2
	ITEM        = 3
	PLAYER_CHAR = '@'
	EXIT_CHAR   = '|'
)

// const (
// 	WALL_CHAR_H          = '─'
// 	WALL_CHAR_V          = '│'
// 	WALL_CHAR_TOPLEFT    = '┌'
// 	WALL_CHAR_TOPRIGHT   = '┐'
// 	WALL_CHAR_BOTTOMLEFT = '└'
// 	WALLCHAR_BOTTOMRIGHT = '┘'
// 	CORRIDOR_CHAR        = '▒'
// 	OUTER_AREA_CHAR      = ' '
// 	INNER_AREA_CHAR      = '·'
// 	EMPTY_CHAR           = ' '
// 	MAP_V                = "║"
// 	MAP_H                = "═"
// 	DOOR_CHAR            = '╬'
// )

const (
	WALL_CHAR_H          = '+'
	WALL_CHAR_V          = '+'
	WALL_CHAR_TOPLEFT    = '"'
	WALL_CHAR_TOPRIGHT   = '"'
	WALL_CHAR_BOTTOMLEFT = '"'
	WALLCHAR_BOTTOMRIGHT = '"'
	CORRIDOR_CHAR        = '#'
	OUTER_AREA_CHAR      = ' '
	INNER_AREA_CHAR      = '.'
	EMPTY_CHAR           = ' '
	MAP_V                = "║"
	MAP_H                = "═"
	DOOR_CHAR            = ')'
	HEALTHBAR            = "#"
)

const (
	IS_OUTER = 0
	IS_INNER = 1
	IS_WALL  = 2
)

const (
	START_ROOM = 0
)

const (
	ENEMY_STEP      = 1
	HOSTILITY_VALUE = 1
)

const (
	MAX_LEVEL = 21
)
