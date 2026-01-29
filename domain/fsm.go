package domain

import (
	"rogue/log"
	"time"
)

var fsm *Fsm

type ActionFunc func(*GameSession) int
type ExitFunc func(*GameSession)
type EnterFunc func(*GameSession)

const (
	IDLE_S = iota
	MOVES_S
	ENEMY_S
	FIGHT_S
	NEXTLEVEL_S
	GAMEOVER_S
	WIN_S
	STATES_NUM_S
)

const (
	UpAction = iota
	DownAction
	LeftAction
	RightAction
	StartAction
	NoAction
)

type FsmState struct {
	action ActionFunc
	exit   ExitFunc
	enter  EnterFunc
}

type Fsm struct {
	states    map[int]FsmState
	cur_state int
}

func Fsm_instance() *Fsm {
	if fsm == nil {
		fsm = &Fsm{}
		fsm.init()
	}
	return fsm
}

func (fsm *Fsm) Get_cur_state() int {
	return fsm.cur_state
}

func (fms *Fsm) init() {
	fsm.states = map[int]FsmState{
		IDLE_S:      {idle_action, nil, nil},
		MOVES_S:     {moves_action, nil, nil},
		ENEMY_S:     {enemy_action, nil, nil},
		FIGHT_S:     {fight_action, nil, fight_enter},
		NEXTLEVEL_S: {nextlevel_action, nil, nil},
		GAMEOVER_S:  {gameover_action, nil, nil},
		WIN_S:       {win_action, nil, nil},
	}
	fsm.cur_state = MOVES_S
}

func (fsm *Fsm) ResetState() {
	fsm.cur_state = MOVES_S
}

func idle_action(game *GameSession) int {
	switch game.Action {
	case LeftAction, RightAction, UpAction, DownAction:
		return MOVES_S
	}
	return MOVES_S
}

func moves_action(game *GameSession) int {
	game.Hero.UpdateBuffs()
	if game.ReachExit() {
		if game.Level.value >= MAX_LEVEL {
			if game.TimeCode == 0 {
				game.TimeCode = time.Now().Unix()
			}
			game.CurStats.state = WIN_S
			game.Store()
			return WIN_S
		}
		return NEXTLEVEL_S
	}

	if game.Hero.put_to_sleep {
		game.Hero.put_to_sleep = false
		return ENEMY_S
	}

	dx, dy := 0, 0
	switch game.Action {
	case LeftAction:
		dx = -1
	case RightAction:
		dx = 1
	case UpAction:
		dy = -1
	case DownAction:
		dy = 1
	}

	if dx != 0 || dy != 0 {
		targetX := game.Hero.position.x + dx
		targetY := game.Hero.position.y + dy

		if enemy := game.GetEnemyAt(targetX, targetY); enemy != nil {
			game.FightEnemy = enemy
			game.FightEnemy.freeze = false
			return fight_action(game)
		} else {
			game.MoveHero(dy, dx)
			game.TakeItem()
			game.Armistice = false
			return ENEMY_S
		}
	}

	game.TakeItem()
	game.Armistice = false
	return MOVES_S
}

func enemy_action(game *GameSession) int {
	state := IDLE_S
	if e, ok := game.Contact(); ok && !game.Armistice {
		log.DebugLog("contact:", game.Hero.position, e.position)
		game.FightEnemy = e
		state = FIGHT_S
	} else {
		game.MoveEnemies()
		game.Armistice = false
	}
	return state
}

func nextlevel_action(game *GameSession) int {
	game.NextLevel(false)
	log.DebugLog("nextlevel_action:", game.Level.value)
	return IDLE_S
}

func fight_action(game *GameSession) int {
	state := IDLE_S
	pwin, ewin := game.Atack(game.FightEnemy)
	if ewin {
		log.DebugLog("gameover, ewin:", game.Hero.health, game.FightEnemy.health)
		state = GAMEOVER_S
	} else if pwin {
		log.DebugLog("pwin:", game.Hero.health, game.FightEnemy.health)
		game.MakeTreasure(game.FightEnemy)
		game.RemoveEnemyByPtr(game.FightEnemy)
		game.CurStats.AddDefeatedEnemy()
	} else {
		if game.FightEnemy.enemy_type == OgreType {
			game.FightEnemy.freeze = true
		}
	}
	return state
}

func fight_enter(game *GameSession) {
	game.FightEnemy.freeze = false
	game.Armistice = true
	log.DebugLog("Armistice:", game.Armistice)
}

func gameover_action(game *GameSession) int {
	return GAMEOVER_S
}

func win_action(game *GameSession) int {
	return WIN_S
}
