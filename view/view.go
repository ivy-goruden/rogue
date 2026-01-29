package view

import (
	"fmt"
	"math"
	"os"
	"rogue/domain"
	"rogue/log"
	"rogue/utils"
	"sort"
	"strings"
	"time"

	"github.com/gbin/goncurses"
)

type Logo struct {
	Logo      rune
	ColorPair int16
}

type Screen struct {
	stdscr *goncurses.Window
}

func MakeScreen() Screen {
	return Screen{}
}

func (scr *Screen) ShowMessage(s string) {
	panic("unimplemented")
}

func (scr *Screen) RenderLoadSaveScreen() {
	scr.Erase()
	scr.DrawLoadScreen()
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderStatsScreen(session *domain.GameSession) {
	scr.Erase()
	scr.DrawLeaderboard(session)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderFoodSelection(session *domain.GameSession) {
	scr.Erase()
	scr.DisplayInventory(session, domain.FoodType)
	scr.DrawStats(session)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderElixirSelection(session *domain.GameSession) {
	scr.Erase()
	scr.DisplayInventory(session, domain.ElixirType)
	scr.DrawStats(session)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderScrollSelection(session *domain.GameSession) {
	scr.Erase()
	scr.DisplayInventory(session, domain.ScrollType)
	scr.DrawStats(session)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderWeaponSelection(session *domain.GameSession) {
	scr.Erase()
	scr.DisplayInventory(session, domain.WeaponType)
	scr.DrawStats(session)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderGameScreen(session *domain.GameSession) {
	scr.DisplayMap(session)
	scr.DisplayHero(session)
	scr.DisplayObjects(session)
	scr.DrawStats(session)
	scr.DrawStatistics(&session.CurStats)
	scr.DrawBuffs(&session.Hero)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderStartScreen(session *domain.GameSession) {
	scr.Erase()
	scr.DrawStartScreen()
	scr.DrawBuffs(&session.Hero)
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderWinScreen() {
	scr.Erase()
	scr.DrawWinScreen()
	scr.stdscr.Refresh()
}

func (scr *Screen) RenderGameOverScreen() {
	scr.Erase()
	scr.DrawGameOverScreen()
	scr.stdscr.Refresh()
}

func (scr *Screen) Close() {
	scr.GetScreen().Timeout(-1)
	goncurses.End()
}

func (scr *Screen) DrawStats(session *domain.GameSession) {
	level := fmt.Sprintf("Level: %d", session.Level.GetValue())
	health := fmt.Sprintf("Health: %s %d/%d", strings.Repeat(domain.HEALTHBAR, max(0, session.Hero.GetHealth()/10)), session.Hero.GetHealth(), session.Hero.GetMaxHealth())
	strength := fmt.Sprintf("Strength: %d", session.Hero.GetStrength())
	dexterity := fmt.Sprintf("Dexterity: %d", session.Hero.GetDexterity())
	gold := fmt.Sprintf("Gold: %d", session.Hero.GetGold())
	level_len := len(level)
	health_len := len(health)
	strength_len := len(strength)
	dexterity_len := len(dexterity)
	gold_len := len(gold)

	total_len := level_len + health_len + strength_len + dexterity_len + gold_len
	space := max((domain.MAP_WIDTH-total_len)/4, 1)

	scr.stdscr.Move(domain.MAP_HEIGHT, 0)
	scr.stdscr.ClearToEOL()
	scr.stdscr.MovePrint(domain.MAP_HEIGHT, 0, level)
	scr.stdscr.MovePrint(domain.MAP_HEIGHT, level_len+space, health)
	scr.stdscr.MovePrint(domain.MAP_HEIGHT, level_len+health_len+space*2, strength)
	scr.stdscr.MovePrint(domain.MAP_HEIGHT, level_len+health_len+strength_len+space*3, dexterity)
	scr.stdscr.MovePrint(domain.MAP_HEIGHT, level_len+health_len+strength_len+dexterity_len+space*4, gold)
}

func (scr *Screen) Init() {
	if scr.stdscr == nil {
		stdscr, err := goncurses.Init()
		if err != nil {
			log.DebugLog(err)
			os.Exit(1)
		}
		scr.stdscr = stdscr
		goncurses.StartColor()
		goncurses.Raw(true)
		goncurses.Echo(false)
		goncurses.Cursor(0)
		goncurses.StdScr().Timeout(16)
		scr.stdscr.Keypad(true)
		scr.InitColors()
	}
}

func (scr *Screen) GetScreen() *goncurses.Window {
	return scr.stdscr
}

func (scr *Screen) GetChar() goncurses.Key {
	return scr.stdscr.GetChar()
}

func (scr *Screen) Erase() {
	scr.stdscr.Clear()
	scr.stdscr.Refresh()
}

func (scr *Screen) InitColors() {
	goncurses.InitPair(1, goncurses.C_RED, goncurses.C_BLACK)
	goncurses.InitPair(2, goncurses.C_BLACK, goncurses.C_WHITE)
	goncurses.InitPair(3, goncurses.C_GREEN, goncurses.C_WHITE)
	goncurses.InitPair(4, goncurses.C_BLUE, goncurses.C_BLACK)
	goncurses.InitPair(5, goncurses.C_YELLOW, goncurses.C_BLACK)
	goncurses.InitPair(6, goncurses.C_CYAN, goncurses.C_BLACK)
	goncurses.InitPair(7, goncurses.C_MAGENTA, goncurses.C_BLACK)
	goncurses.InitPair(8, goncurses.C_WHITE, goncurses.C_BLACK) //Map
}

func (scr *Screen) GetEnemyLogo(t domain.EnemyType) Logo {
	var logo Logo
	switch t {
	case domain.ZombieType:
		logo = Logo{Logo: 'Z', ColorPair: 3}
	case domain.VampireType:
		logo = Logo{Logo: 'V', ColorPair: 1}
	case domain.GhostType:
		logo = Logo{Logo: 'G', ColorPair: 2}
	case domain.OgreType:
		logo = Logo{Logo: 'O', ColorPair: 4}
	case domain.SnakeMageType:
		logo = Logo{Logo: 'S', ColorPair: 4}
	case domain.MimicType:
		logo = Logo{Logo: 'M', ColorPair: 2}
	}
	return logo
}

func (scr *Screen) GetItemLogo() Logo {
	return Logo{Logo: '*', ColorPair: 5}
}

func (scr *Screen) DisplayHero(game_session *domain.GameSession) {
	scr.DrawHero(game_session)
}

func (scr *Screen) DisplayObjects(game_session *domain.GameSession) {
	scr.DrawItems(&game_session.Level)
	scr.DrawEnemies(&game_session.Level)
}

func (scr *Screen) DrawItems(level *domain.Level) {
	for _, item := range level.Items {
		if level.Fog[item.GetY()][item.GetX()] == 0 {
			scr.stdscr.ColorOn(scr.GetItemLogo().ColorPair)
			scr.stdscr.MovePrintf(item.GetY(), item.GetX(), "%c", scr.GetItemLogo().Logo)
			scr.stdscr.ColorOff(scr.GetItemLogo().ColorPair)
		}
	}
}

func (scr *Screen) DrawEnemies(level *domain.Level) {
	for _, enemy := range level.Enemies {
		if level.Fog[enemy.GetY()][enemy.GetX()] == 0 && enemy.IsVisible() {
			if enemy.EnemyType() == domain.MimicType && !enemy.IsChasing() {
				scr.stdscr.ColorOn(scr.GetItemLogo().ColorPair)
				scr.stdscr.MovePrintf(enemy.GetY(), enemy.GetX(), "%c", scr.GetItemLogo().Logo)
				scr.stdscr.ColorOff(scr.GetItemLogo().ColorPair)
			} else {
				scr.stdscr.ColorOn(scr.GetEnemyLogo(enemy.EnemyType()).ColorPair)
				scr.stdscr.MovePrintf(enemy.GetY(), enemy.GetX(), "%c", scr.GetEnemyLogo(enemy.EnemyType()).Logo)
				scr.stdscr.ColorOff(scr.GetEnemyLogo(enemy.EnemyType()).ColorPair)
			}
		}
	}
}

func (scr *Screen) DisplayMap(game_session *domain.GameSession) {
	scr.DrawMap(game_session)
}

func (scr *Screen) DrawHero(game_session *domain.GameSession) {
	x := game_session.Hero.Hx()
	y := game_session.Hero.Hy()
	disperseFog(&game_session.Level, x, y)
	scr.stdscr.MovePrintf(y, x, "%c", domain.PLAYER_CHAR)
}

func (scr *Screen) DrawMap(game_session *domain.GameSession) {
	scr.stdscr.ColorOn(8)
	for y, row := range game_session.Playground.Playground {
		for x, cell := range row {
			if game_session.Level.Fog[y][x] == 0 {
				scr.stdscr.MovePrintf(y, x, "%c", cell)
			} else {
				scr.stdscr.MovePrintf(y, x, "%c", ' ')
			}
		}
	}
	scr.stdscr.ColorOff(8)
}

func disperseFog(level *domain.Level, x, y int) {
	for dx := -domain.FOG_RADIUS; dx <= domain.FOG_RADIUS; dx++ {
		dy := int(math.Sqrt(float64(math.Pow(float64(domain.FOG_RADIUS), 2) - math.Pow(float64(dx), 2))))
		for i := -dy; i <= dy; i++ {
			if x+dx >= 0 && x+dx < domain.MAP_WIDTH && y+i >= 0 && y+i < domain.MAP_HEIGHT {
				level.Fog[y+i][x+dx] = 0
			}
		}
	}
}

func (scr *Screen) DisplayInventory(game_session *domain.GameSession, itemType domain.ItemType) {
	// Draw frame
	for i := range domain.MAP_WIDTH {
		scr.stdscr.MovePrintf(0, i, "=")
		scr.stdscr.MovePrintf(domain.MAP_HEIGHT-1, i, "=")
	}
	for i := range domain.MAP_HEIGHT {
		scr.stdscr.MovePrintf(i, 0, "=")
		scr.stdscr.MovePrintf(i, domain.MAP_WIDTH-1, "=")
	}

	scr.stdscr.MovePrintf(1, (domain.MAP_WIDTH-len("Inventory:"))/2, "Inventory:")
	items := game_session.Hero.GetItemByType(itemType)

	var drawer ItemDrawer
	switch itemType {
	case domain.WeaponType:
		drawer = WeaponList(items)
	case domain.FoodType:
		drawer = FoodList(items)
	case domain.ElixirType:
		drawer = ElixirList(items)
	case domain.ScrollType:
		drawer = ScrollList(items)
	default:
		return
	}
	drawer.Draw(scr, game_session.Hero)
}

func (scr *Screen) DrawDeadScreen() {
	maxY := domain.MAP_HEIGHT
	maxX := domain.MAP_WIDTH
	gameOver := []string{
		"=========================",
		"                         ",
		"    /\\_/\\              ",
		"    ( o.o )  GAME OVER   ",
		"     > ^ <               ",
		"                         ",
		"=========================",
		"       Q  - выход        ",
	}
	startY := maxY/2 - len(gameOver)/2
	for i, line := range gameOver {
		startX := maxX/2 - len(line)/2
		scr.stdscr.MovePrintf(startY+i, startX, "%s", line)
	}
}

func (scr *Screen) DrawGameOverScreen() {
	maxY, maxX := scr.stdscr.MaxYX()
	lines := []string{
		"=========================",
		"                         ",
		"    /\\_/\\              ",
		"    ( o.o )  GAME OVER   ",
		"     > ^ <               ",
		"                         ",
		"=========================",
		"     ESC  -  to menu     ",
	}
	startY := max(0, maxY/2-len(lines)/2)
	for i, line := range lines {
		if len(line) > maxX {
			line = line[:maxX]
		}
		startX := max(0, maxX/2-len(line)/2)
		scr.stdscr.MovePrintf(startY+i, startX, "%s", line)
	}
}

func (scr *Screen) DrawWinScreen() {
	maxY, maxX := scr.stdscr.MaxYX()
	lines := []string{
		"=========================",
		"                         ",
		"    /\\_/\\              ",
		"    ( ^.^ )   YOU WIN    ",
		"     > ^ <               ",
		"                         ",
		"=========================",
		"     ESC  -  to menu     ",
	}
	startY := max(0, maxY/2-len(lines)/2)
	for i, line := range lines {
		if len(line) > maxX {
			line = line[:maxX]
		}
		startX := max(0, maxX/2-len(line)/2)
		scr.stdscr.MovePrintf(startY+i, startX, "%s", line)
	}
}

func (scr *Screen) DrawStartScreen() {
	maxY := domain.MAP_HEIGHT
	maxX := domain.MAP_WIDTH
	start := []string{
		"=========================",
		"                         ",
		"    /\\_/\\              ",
		"    ( o.o )  ROGUE       ",
		"     > ^ <               ",
		"                         ",
		"=========================",
		"       Press a Key       ",
	}
	startY := (maxY - len(start)) / 2
	for i, line := range start {
		startX := maxX/2 - len(line)/2
		scr.stdscr.MovePrintf(startY+i, startX, "%s", line)
	}
	scr.stdscr.MovePrintf(maxY/2+5, maxX/2-len("1. Start")/2, "1. Start")
	scr.stdscr.MovePrintf(maxY/2+6, maxX/2-len("2.  Load")/2, "2. Load")
	scr.stdscr.MovePrintf(maxY/2+7, maxX/2-len("3.  Leaderboard")/2, "3. Leaderboard")
	scr.stdscr.MovePrintf(maxY/2+8, maxX/2-len("4.  Quit")/2, "4. Quit")
}

func (scr *Screen) DrawStopScreen() {
	maxY := domain.MAP_HEIGHT
	maxX := domain.MAP_WIDTH
	start := []string{
		"=========================",
		"                         ",
		"   /\\_/\\               ",
		"   ( o.o ) Хочешь Выйти? ",
		"    > ^ <                ",
		"                         ",
		"=========================",
		"    Да(Q)     Нет(Esc)   ",
	}
	startY := (maxY - len(start)) / 2
	for i, line := range start {
		startX := maxX/2 - len(line)/2
		scr.stdscr.MovePrintf(startY+i, startX, "%s", line)
	}
}

type ItemDrawer interface {
	Draw(scr *Screen, player domain.Player)
}

type WeaponList []domain.Item

func (items WeaponList) Draw(scr *Screen, player domain.Player) {
	if len(items) == 0 && player.GetHandWeapon().GetStrength() == 0 {
		scr.stdscr.MovePrintf(4, (domain.MAP_WIDTH-len("== No Weapons =="))/2, "== No Weapons ==")
		return
	}
	scr.stdscr.MovePrintf(3, 3, "No.")
	scr.stdscr.MovePrintf(3, 13, "Strength")
	if player.GetHandWeapon().GetStrength() != 0 {
		scr.stdscr.MovePrintf(4, 3, "%d", 0)
		scr.stdscr.MovePrintf(4, 13, "%d", player.GetHandWeapon().GetStrength())
		scr.stdscr.MovePrintf(4, 18, "%s", "Away the weapon")
	}
	for i, item := range items {
		scr.stdscr.MovePrintf(5+i, 3, "%d", i+1)
		scr.stdscr.MovePrintf(5+i, 13, "%d", item.GetStrength())
	}
}

type FoodList []domain.Item

func (items FoodList) Draw(scr *Screen, player domain.Player) {
	if len(items) == 0 {
		scr.stdscr.MovePrintf(4, (domain.MAP_WIDTH-len("== No Food =="))/2, "== No Food ==")
		return
	}
	scr.stdscr.MovePrintf(3, 3, "No.")
	scr.stdscr.MovePrintf(3, 13, "Health")
	for i, item := range items {
		scr.stdscr.MovePrintf(4+i, 3, "%d", i+1)
		scr.stdscr.MovePrintf(4+i, 13, "%d", item.GetHealth())
	}
}

type ElixirList []domain.Item

func (items ElixirList) Draw(scr *Screen, player domain.Player) {
	if len(items) == 0 {
		scr.stdscr.MovePrintf(4, (domain.MAP_WIDTH-len("== No Elixir =="))/2, "== No Elixir ==")
		return
	}
	scr.stdscr.MovePrintf(3, 3, "No.")
	scr.stdscr.MovePrintf(3, 9, "Name")
	scr.stdscr.MovePrintf(3, 36, "Value")
	scr.stdscr.MovePrintf(3, 46, "Duration")
	for i, item := range items {
		name := ""
		switch item.GetFeature() {
		case domain.FeatureType(domain.DexterityFeature):
			name = "Elixir of Dexterity"
		case domain.FeatureType(domain.StrengthFeature):
			name = "Elixir of Strength"
		case domain.FeatureType(domain.HealthFeature):
			name = "Elixir of Max Health"
		}
		scr.stdscr.MovePrintf(4+i, 3, "%d", i+1)
		scr.stdscr.MovePrintf(4+i, 6, "%s", name)
		scr.stdscr.MovePrintf(4+i, 6, "%s", name)
		scr.stdscr.MovePrintf(4+i, 36, "%d", item.GetValue())
		scr.stdscr.MovePrintf(4+i, 46, "%d", item.GetDuration())
	}
}

type ScrollList []domain.Item

func (items ScrollList) Draw(scr *Screen, player domain.Player) {
	if len(items) == 0 {
		scr.stdscr.MovePrintf(4, (domain.MAP_WIDTH-len("== No Scroll =="))/2, "== No Scroll ==")
		return
	}
	scr.stdscr.MovePrintf(3, 3, "No.")
	scr.stdscr.MovePrintf(3, 9, "Name")
	scr.stdscr.MovePrintf(3, 36, "Value")
	for i, item := range items {
		name := ""
		switch item.GetFeature() {
		case domain.FeatureType(domain.DexterityFeature):
			name = "Scroll of Dexterity"
		case domain.FeatureType(domain.StrengthFeature):
			name = "Scroll of Strength"
		case domain.FeatureType(domain.HealthFeature):
			name = "Scroll of Max Health"
		}
		scr.stdscr.MovePrintf(4+i, 3, "%d", i+1)
		scr.stdscr.MovePrintf(4+i, 6, "%s", name)
		scr.stdscr.MovePrintf(4+i, 6, "%s", name)
		scr.stdscr.MovePrintf(4+i, 36, "%d", item.GetValue())
	}
}

func (scr *Screen) DrawLoadScreen() {
	// Draw frame
	for i := range domain.MAP_WIDTH {
		scr.stdscr.MovePrintf(0, i, "=")
		scr.stdscr.MovePrintf(domain.MAP_HEIGHT-1, i, "=")
	}
	for i := range domain.MAP_HEIGHT {
		scr.stdscr.MovePrintf(i, 0, "=")
		scr.stdscr.MovePrintf(i, domain.MAP_WIDTH-1, "=")
	}
	sessions := domain.GetGameSessions()

	if len(sessions) == 0 {
		scr.stdscr.MovePrintf(4, (domain.MAP_WIDTH-len("== No Sessions =="))/2, "== No Sessions ==")
		return
	}
	scr.DrawSessions(sessions)
}

func (scr *Screen) DrawSessions(sessions []domain.GameSession) {
	start_y := 3
	start_x := 3
	scr.stdscr.MovePrintf(start_y, start_x, "%5s%s %12s %30s %10s", "N", "", "Level", "Session", "Tres")
	for i, v := range sessions {
		scr.stdscr.MovePrintf(start_y+i, start_x, "%5d%s %10s%2d %30s %10d", i, ".", "Level ", v.Level.GetValue(), utils.GetSessName(v.TimeCode), v.Hero.GetGold())
	}
}

func (scr *Screen) DrawStatistics(stats *domain.Statistics) {
	// Позиция для статистики - справа от игрового поля
	statX := domain.MAP_WIDTH + 2 // Отступ от правой границы поля
	startY := 1                   // Начинаем с первой строки под верхней границей

	// Очищаем область статистики
	for i := 0; i < 15; i++ {
		scr.stdscr.MovePrintf(startY+i, statX, "%-30s", " ") // 30 символов для статистики
	}

	scr.stdscr.MovePrintf(startY, statX, "Enemies killed: %d", stats.GetEnemies())
	startY++

	scr.stdscr.MovePrintf(startY, statX, "Food collected: %d", stats.GetFood())
	startY++

	scr.stdscr.MovePrintf(startY, statX, "Elixirs used: %d", stats.GetElixirs())
	startY++

	scr.stdscr.MovePrintf(startY, statX, "Scrolls used: %d", stats.GetScrolls())
	startY++

	scr.stdscr.MovePrintf(startY, statX, "Player hits: %d", stats.GetHits())
	startY++

	scr.stdscr.MovePrintf(startY, statX, "Enemy hits: %d", stats.GetEnemyHits())
	startY++

	scr.stdscr.MovePrintf(startY, statX, "Tiles explored: %d", stats.GetTiles())
}

func (scr *Screen) DrawBuffs(player *domain.Player) {
	buffs := player.GetBuffs()
	totalBuffs := len(buffs.GetHealthBuffs()) + len(buffs.GetDexterityBuffs()) + len(buffs.GetStrengthBuffs())
	if totalBuffs == 0 {
		return
	}

	// Позиция для баффов - под статистикой
	statX := domain.MAP_WIDTH + 2 // Та же колонка, что и статистика
	startY := 9                   // Начинаем после статистики (предположительно 8 строк статистики)

	// Очищаем область баффов (например, 10 строк)
	for i := 0; i < 10; i++ {
		scr.stdscr.MovePrintf(startY+i, statX, "%-30s", " ")
	}

	// Заголовок баффов
	scr.stdscr.MovePrintf(startY, statX, "+=[ BUFFS ]================+")
	startY++

	currentTime := int(time.Now().Unix())
	buffCount := 1

	// Здоровье
	for _, b := range buffs.GetHealthBuffs() {
		remaining := b.GetStart() + b.GetEffect() - currentTime
		if remaining > 0 {
			seconds := remaining % 60
			scr.stdscr.MovePrintf(startY, statX, "| %d+%d HP (%dс)",
				buffCount, b.GetValue(), seconds)
			startY++
			buffCount++
		}
	}

	// Ловкость
	for _, b := range buffs.GetDexterityBuffs() {
		remaining := b.GetStart() + b.GetEffect() - currentTime
		if remaining > 0 {
			seconds := remaining % 60
			scr.stdscr.MovePrintf(startY, statX, "| %d+%d DEX (%dс)",
				buffCount, b.GetValue(), seconds)
			startY++
			buffCount++
		}
	}

	// Сила
	for _, b := range buffs.GetStrengthBuffs() {
		remaining := b.GetStart() + b.GetEffect() - currentTime
		if remaining > 0 {
			seconds := remaining % 60
			scr.stdscr.MovePrintf(startY, statX, "| %d+%d STR (%dс)",
				buffCount, b.GetValue(), seconds)
			startY++
			buffCount++
		}
	}

	// Нижняя граница
	scr.stdscr.MovePrintf(startY, statX, "+==========================+")
}

func (scr *Screen) DrawLeaderboard(session *domain.GameSession) {
	sessions := domain.GetGameSessions()
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CurStats.GetTreasures() > sessions[j].CurStats.GetTreasures()
	})
	startY := 1
	startX := 3
	scr.stdscr.MovePrintf(startY, startX, "%3s  %9s  %5s  %-40s  %-10s", "N", "Treasures", "Level", "Session", "Result")
	scr.stdscr.MovePrintf(startY+1, startX, "%s", strings.Repeat("-", 3+2+9+2+5+2+40+2+10))
	i := 0
	for _, v := range sessions {
		if v.CurStats.GetResult() == "OTHER" {
			continue
		}
		row1Y := startY + 2 + i*2
		row2Y := row1Y + 1
		scr.stdscr.MovePrintf(row1Y, startX, "%3d.  %9d  %5d  %-40s  %-10s", i+1, v.CurStats.GetTreasures(), v.Level.GetValue(), utils.GetSessName(v.TimeCode), v.CurStats.GetResult())
		scr.stdscr.MovePrintf(
			row2Y,
			startX,
			"%3s   E:%3d  F:%3d  Elx:%3d  Scr:%3d  Hits:%3d  EH:%3d  Tiles:%4d",
			"",
			v.CurStats.GetEnemies(),
			v.CurStats.GetFood(),
			v.CurStats.GetElixirs(),
			v.CurStats.GetScrolls(),
			v.CurStats.GetHits(),
			v.CurStats.GetEnemyHits(),
			v.CurStats.GetTiles(),
		)
		i++
	}
}
