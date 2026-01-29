package data

import (
	"fmt"
	"path/filepath"
	"rogue/log"
	"rogue/utils"
	"strings"
	"time"
)

func GetLatestSessionFile(dir string, num int) ([]string, error) {
	files, err := filepath.Glob(dir + "/rogue_*.session")
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("файлы сессий не найдены")
	}
	result := make([]string, min(num, len(files)))
	for i := range files[:min(num, len(files))] {
		result[i] = strings.Split(files[i], "/")[1]
	}
	return result, nil
}

func newerFile(f1, f2 string) bool {
	t1 := parseTimestamp(f1)
	t2 := parseTimestamp(f2)
	return t1.After(t2)
}

func parseTimestamp(filename string) time.Time {
	// rogue_20260106_125955.session -> 20260106_125955
	parts := strings.Split(strings.TrimSuffix(filepath.Base(filename), ".session"), "_")
	if len(parts) < 2 {
		return time.Time{}
	}
	timestampStr := parts[1] + "_" + parts[2]
	t, _ := time.Parse("20060102_150405", timestampStr)
	return t
}

func FindSessionFilename(dir string) string {
	file, err := GetLatestSessionFile(dir, 1)
	if err != nil {
		log.DebugLog("loading:filename error:", err)
		return ""
	}
	return file[0]
}

func GenerateSessionFilename() string {
	return utils.GetSessName(time.Now().Unix())
}
