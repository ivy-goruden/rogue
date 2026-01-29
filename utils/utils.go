package utils

import (
	"fmt"
	"time"
)

type AssocArray map[string]interface{}

func ArrayMap[S any, T any](src []S, f func(S) T) []T {
	res := make([]T, 0, len(src))
	for _, v := range src {
		res = append(res, f(v))
	}
	return res
}

func GetSessName(t int64) string {
	timestamp := time.Unix(t, 0).Format("20060102_150405")
	return fmt.Sprintf("rogue_%s.session", timestamp)
}
