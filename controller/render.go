package controller

func (c *Controller) Render() {
	switch c.currentMode {
	case ModeStart:
		session := c.model.GetSession()
		c.view.RenderStartScreen(session)

	case ModeGame:
		session := c.model.GetSession()
		c.view.RenderGameScreen(session)

	case ModeWeaponSelect:
		session := c.model.GetSession()
		c.view.RenderWeaponSelection(session)

	case ModeFoodSelect:
		session := c.model.GetSession()
		c.view.RenderFoodSelection(session)

	case ModeElixirSelect:
		session := c.model.GetSession()
		c.view.RenderElixirSelection(session)

	case ModeScrollSelect:
		session := c.model.GetSession()
		c.view.RenderScrollSelection(session)

	case ModeStats:
		session := c.model.GetSession()
		c.view.RenderStatsScreen(session)

	case ModeLoadSave:
		c.view.RenderLoadSaveScreen()

	case ModeWin:
		c.view.RenderWinScreen()

	case ModeGameOver:
		c.view.RenderGameOverScreen()
	}
}
