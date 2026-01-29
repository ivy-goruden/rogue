package main

import (
	"rogue/controller"
	"rogue/domain"
	"rogue/view"
)
func main() {
	game := domain.MakeGameSession()
	screen := view.MakeScreen()
	screen.Init()
	defer screen.Close()
	dispatcher := controller.MakeController(&game, &screen)
	dispatcher.Run()
}
