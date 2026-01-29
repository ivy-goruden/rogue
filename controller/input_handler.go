package controller

import (
	"rogue/domain"

	"github.com/gbin/goncurses"
)

type InputHandler interface {
	HandleInput(ev goncurses.Key, c *Controller)
}

type GamePlayHandler struct{}
type StartScreenHandler struct{}
type StopScreenHandler struct{}
type WeaponSelectHandler struct{}
type FoodSelectHandler struct{}
type ElixirSelectHandler struct{}
type ScrollSelectHandler struct{}
type StatsScreenHandler struct{}
type LoadSaveHandler struct{}
type WinScreenHandler struct{}
type GameOverScreenHandler struct{}

func (h GamePlayHandler) HandleInput(ev goncurses.Key, c *Controller) {
	if c.GetFsmState() == domain.GAMEOVER_S {
		c.SwitchMode(ModeGameOver)
		return
	}
	if c.GetFsmState() == domain.WIN_S {
		c.SwitchMode(ModeWin)
		return
	}

	switch ev {
	case 'q', 'Q', 27:
		c.SwitchMode(ModeStart)
	case 'e', 'E':
		c.SwitchMode(ModeScrollSelect)
	case 'h', 'H':
		c.SwitchMode(ModeWeaponSelect)
	case 'j', 'J':
		c.SwitchMode(ModeFoodSelect)
	case 'k', 'K':
		c.SwitchMode(ModeElixirSelect)
	default:
		c.UpdateGameState(ev)
	}
}

func (h StartScreenHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case '1':
		c.StartNewGame()
	case '2':
		c.SwitchMode(ModeLoadSave)
	case '3':
		c.SwitchMode(ModeStats)
	case '4', 27:
		c.Quit()
	}
}

func (h StopScreenHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 'q', 'Q':
		c.Quit()
	}
}

func (h WeaponSelectHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeGame)
	case 'e', 'E':
		c.SwitchMode(ModeScrollSelect)
	case 'h', 'H':
		c.SwitchMode(ModeWeaponSelect)
	case 'j', 'J':
		c.SwitchMode(ModeFoodSelect)
	case 'k', 'K':
		c.SwitchMode(ModeElixirSelect)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		c.UseItem(domain.WeaponType, int(ev-'0'))
	}
}

func (h FoodSelectHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeGame)
	case 'e', 'E':
		c.SwitchMode(ModeScrollSelect)
	case 'h', 'H':
		c.SwitchMode(ModeWeaponSelect)
	case 'j', 'J':
		c.SwitchMode(ModeFoodSelect)
	case 'k', 'K':
		c.SwitchMode(ModeElixirSelect)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		c.UseItem(domain.FoodType, int(ev-'0'))
	}
}

func (h ElixirSelectHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeGame)
	case 'e', 'E':
		c.SwitchMode(ModeScrollSelect)
	case 'h', 'H':
		c.SwitchMode(ModeWeaponSelect)
	case 'j', 'J':
		c.SwitchMode(ModeFoodSelect)
	case 'k', 'K':
		c.SwitchMode(ModeElixirSelect)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		c.UseItem(domain.ElixirType, int(ev-'0'))
	}
}

func (h ScrollSelectHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeGame)
	case 'e', 'E':
		c.SwitchMode(ModeScrollSelect)
	case 'h', 'H':
		c.SwitchMode(ModeWeaponSelect)
	case 'j', 'J':
		c.SwitchMode(ModeFoodSelect)
	case 'k', 'K':
		c.SwitchMode(ModeElixirSelect)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		c.UseItem(domain.ScrollType, int(ev-'0'))
	}
}

func (h StatsScreenHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeStart)
	}
}

func (h LoadSaveHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeStart)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		c.LoadGame(int(ev) - '0')
	}
}

func (h WinScreenHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeStart)
	}
}

func (h GameOverScreenHandler) HandleInput(ev goncurses.Key, c *Controller) {
	switch ev {
	case 27:
		c.SwitchMode(ModeStart)
	}
}
