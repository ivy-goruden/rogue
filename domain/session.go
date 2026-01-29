package domain

import (
	"rogue/data"
	"rogue/domain/serializer"
	"rogue/log"
	"sort"
)

func LoadSession(file string) (GameSession, error) {
	loadedGame := &GameSession{}
	serial := serializer.MakeSerializer()
	fileHandler := serializer.MakeFileHandler("sessions")
	if err := fileHandler.LoadObject(file, loadedGame, serial); err != nil {
		log.DebugLog("Ошибка загрузки сессии: ", err)
		return GameSession{}, err
	}
	return *loadedGame, nil
}

func BetterStats(game1, game2 *GameSession) bool {
	return game1.Hero.Treasures() > game2.Hero.Treasures()
}

func GetGameSessions() []GameSession {
	result := []GameSession{}

	files, err := data.GetLatestSessionFile("sessions", 1000)

	if err == nil {
		games := []GameSession{}
		for _, v := range files {
			game, err := LoadSession(v)
			if err != nil {
				log.DebugLog("Некоторые сессии загрузить не удалось: ", err)
			} else {
				games = append(games, game)
			}
		}
		if len(games) > 0 {
			sort.Slice(games, func(i, j int) bool {
				return BetterStats(&games[i], &games[j])
			})
			nums := min(10, len(games))
			stats := make([]Statistics, nums)
			result = make([]GameSession, nums)
			for i := 0; i < nums; i++ {
				stats[i] = games[i].CurStats
				result[i] = games[i]
			}
			for i := range result {
				result[i].Stats = stats
			}
		}
	} else {
		log.DebugLog("Файлы получить не удалось: ", err)
	}
	return result
}
