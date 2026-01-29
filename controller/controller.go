package controller

import (
	"rogue/domain"
	"rogue/view"
	"time"
)

type GameMode int

const (
	ModeGame GameMode = iota
	ModeStart
	ModeStop
	ModeWeaponSelect
	ModeFoodSelect
	ModeElixirSelect
	ModeScrollSelect
	ModeStats
	ModeLoadSave
	ModeWin
	ModeGameOver
)

type Controller struct {
	currentMode GameMode
	model       *domain.GameSession
	view        *view.Screen
	running     bool
}

func MakeController(game *domain.GameSession, view *view.Screen) Controller {
	return Controller{
		currentMode: ModeStart,
		model:       game,
		view:        view,
		running:     false,
	}
}

func (c *Controller) SwitchMode(newMode GameMode) {
	time.Sleep(48 * time.Millisecond)
	c.currentMode = newMode
	c.Render()
}

func (c *Controller) GetCurrentInputHandler() InputHandler {
	switch c.currentMode {
	case ModeGame:
		return &GamePlayHandler{}
	case ModeStart:
		return &StartScreenHandler{}
	case ModeStop:
		return &StopScreenHandler{}
	case ModeWeaponSelect:
		return &WeaponSelectHandler{}
	case ModeFoodSelect:
		return &FoodSelectHandler{}
	case ModeElixirSelect:
		return &ElixirSelectHandler{}
	case ModeScrollSelect:
		return &ScrollSelectHandler{}
	case ModeStats:
		return &StatsScreenHandler{}
	case ModeLoadSave:
		return &LoadSaveHandler{}
	case ModeWin:
		return &WinScreenHandler{}
	case ModeGameOver:
		return &GameOverScreenHandler{}
	default:
		return nil
	}
}

func (c *Controller) Run() {
	c.running = true
	c.SwitchMode(ModeStart)

	for c.running {
		if ev := c.view.GetScreen().GetChar(); ev != -1 {
			handler := c.GetCurrentInputHandler()
			if handler != nil {
				handler.HandleInput(ev, c)
			}
		}
	}
}

func (c *Controller) Quit() {
	c.running = false
}
