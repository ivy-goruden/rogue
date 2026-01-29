package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const LOG_FILE = "rogue.log"

func DebugLog(args ...interface{}) error {
	var filename string = LOG_FILE
	// Создаём/открываем файл (дописывание)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %w", filename, err)
	}
	defer f.Close()

	// Буферизированный writer для производительности
	w := log.New(f, "", log.LstdFlags)

	// Получаем информацию о вызывающей функции
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	shortFile := filepath.Base(file)

	// Форматируем сообщение с timestamp и стек-трейсом
	caller := fmt.Sprintf("%s:%d", shortFile, line)

	// Соединяем все аргументы
	msg := fmt.Sprint(args...)
	logEntry := fmt.Sprintf("[%s] %s\n", caller, msg)

	// Записываем
	err = w.Output(2, logEntry)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}

	return nil
}
