package controller

import (
	"rogue/domain"

	"github.com/gbin/goncurses"
)

func (c *Controller) ResetState() {
	fsm := domain.Fsm_instance()
	fsm.ResetState()
}

func (c *Controller) GetFsmState() int {
	fsm := domain.Fsm_instance()
	return fsm.Get_cur_state()
}

func (c *Controller) InputToGameAction(input goncurses.Key) int {
	var action int
	switch input {
	case 'a', 'A':
		action = domain.LeftAction
	case 'd', 'D':
		action = domain.RightAction
	case 'w', 'W':
		action = domain.UpAction
	case 's', 'S':
		action = domain.DownAction
	default:
		action = domain.NoAction
	}
	return action
}

func (c *Controller) UpdateGameState(key goncurses.Key) {
	c.model.HandleUserInput(c.InputToGameAction(key))
	c.model.UpdateGameState()
	if c.GetFsmState() == domain.IDLE_S {
		c.Render()
	}
}

func (c *Controller) UseItem(itemtype domain.ItemType, index int) {
	c.model.UseItem(itemtype, index)
	c.Render()
}

func (c *Controller) StartNewGame() {
	c.model.NewGame()
	c.ResetState()
	c.SwitchMode(ModeGame)
}

func (c *Controller) LoadGame(index int) {
	err := c.model.LoadFromSession(index)
	if err == nil {
		c.ResetState()
		c.SwitchMode(ModeGame)
	}
}
